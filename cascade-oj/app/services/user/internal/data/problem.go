package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/problemset_includes"
	"cascade-oj/pkg/cache"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
)

type ProblemRepo struct {
	data *Data
	log  *log.Helper
}

func (problemRepo *ProblemRepo) GetProblems(ctx context.Context, contestID int64) ([]*biz.Problem, error) {
	cacheBytes, err := problemRepo.data.GetCache(ctx, fmt.Sprintf(cache.ProblemListCacheKeyFmt, contestID))
	// cache hit
	if err == nil && cacheBytes != nil {
		problemRepo.log.Infof("[cache] hit from user calling GetProblems")
		var problems []*biz.Problem
		if err := json.Unmarshal(cacheBytes, &problems); err == nil {
			return problems, nil
		}
		problemRepo.log.Errorf("[cache] failed to unmarshal problem list from cache: %v", err)
	}

	// singleflight
	// 返回值 shared 表示结果是否被多个调用者共享，可用于监控
	res, err, _ := sfGroup.Do(fmt.Sprintf(cache.ProblemListCacheKeyFmt, contestID), func() (interface{}, error) {
		// 必须超时控制
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// cache miss
		problemRepo.log.Infof("[cache] miss from user calling GetProblems")

		queryProblemsID, err := problemRepo.data.db.ProblemSet_Includes.Query().
			Select(problemset_includes.FieldProblemID).
			Where(problemset_includes.ProblemSetIDEQ(contestID)).
			All(ctx)
		if err != nil {
			return nil, err
		}

		problems := make([]*biz.Problem, 0, len(queryProblemsID))
		for i := 0; i < len(queryProblemsID); i++ {
			queryProblem, err := problemRepo.data.db.Problem.Query().
				Where(problem.IDEQ(queryProblemsID[i].ProblemID)).
				Only(ctx)
			if err != nil {
				return nil, err
			}
			problems = append(problems, &biz.Problem{
				ID:            queryProblem.ID,
				Title:         queryProblem.Title,
				TimeLimitMs:   int32(queryProblem.TimeLimitMs),
				MemoryLimitMb: int32(queryProblem.MemoryLimitKB / 1024),
			})
		}

		// set cache
		cacheBytes, err = json.Marshal(problems)
		if err != nil {
			problemRepo.log.Errorf("[cache] failed to marshal problem list for caching: %v", err)
		} else {
			err = problemRepo.data.SetCache(ctx, fmt.Sprintf(cache.ProblemListCacheKeyFmt, contestID), cacheBytes, 5*time.Minute)
			if err != nil {
				problemRepo.log.Errorf("[cache] failed to set cache for problem list: %v", err)
			}
		}

		return problems, nil
	})

	if err != nil {
		return nil, err
	}
	return res.([]*biz.Problem), nil
}

func (problemRepo *ProblemRepo) GetSingleProblem(ctx context.Context, problemID int64) (*biz.DetailedProblem, error) {
	cacheBytes, err := problemRepo.data.GetCache(ctx, fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, problemID))
	// cache hit
	if err == nil && cacheBytes != nil {
		problemRepo.log.Infof("[cache] hit from user calling GetSingleProblem")
		var problem *biz.DetailedProblem
		if err := json.Unmarshal(cacheBytes, &problem); err == nil {
			return problem, nil
		}
		problemRepo.log.Errorf("[cache] failed to unmarshal problem from cache: %v", err)
	}

	// singleflight
	// 返回值 shared 表示结果是否被多个调用者共享，可用于监控
	res, err, _ := sfGroup.Do(fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, problemID), func() (interface{}, error) {
		// 必须超时控制
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// cache miss
		problemRepo.log.Infof("[cache] miss from user calling GetSingleProblem")

		queryProblem, err := problemRepo.data.db.Problem.Query().
			WithCreator().
			Where(problem.IDEQ(problemID)).
			Only(ctx)
		if err != nil {
			return nil, err
		}

		detailedProblem := &biz.DetailedProblem{
			Problem: biz.Problem{
				ID:            queryProblem.ID,
				Title:         queryProblem.Title,
				TimeLimitMs:   int32(queryProblem.TimeLimitMs),
				MemoryLimitMb: int32(queryProblem.MemoryLimitKB / 1024),
			},
			CreatorUsername: queryProblem.Edges.Creator.Username,
			Description:     queryProblem.Description,
		}

		// set cache
		cacheBytes, err = json.Marshal(detailedProblem)
		if err != nil {
			problemRepo.log.Errorf("[cache] failed to marshal problem for caching: %v", err)
		} else {
			err = problemRepo.data.SetCache(ctx, fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, problemID), cacheBytes, 5*time.Minute)
			if err != nil {
				problemRepo.log.Errorf("[cache] failed to set cache for problem: %v", err)
			}
		}

		return detailedProblem, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*biz.DetailedProblem), nil
}

func NewProblemRepo(data *Data, logger log.Logger) biz.ProblemRepo {
	return &ProblemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 延迟双删 Problem 列表或单个 Problem 缓存
func (problemRepo *ProblemRepo) ProblemCacheConsumer(ctx context.Context, msg *mq.ProblemCacheMsg) error {
	var key string
	if msg.Scale == "single" {
		key = fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, msg.ProblemID)
	} else {
		// handle as list cache by default
		// problemID will be used as contestID in this case
		key = fmt.Sprintf(cache.ProblemListCacheKeyFmt, msg.ProblemID)
	}

	err := problemRepo.data.DeleteCache(ctx, key)
	if err != nil {
		problemRepo.log.Errorf("[cache] failed to delete cache for problem %d: %v", *msg, err)
		return err
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		err := problemRepo.data.DeleteCache(ctx, key)
		if err != nil {
			problemRepo.log.Errorf("[cache] failed to delayed delete cache for problem %d: %v", *msg, err)
		}
	}()
	return nil
}
