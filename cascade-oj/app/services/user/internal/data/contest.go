package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/problemset"
	"cascade-oj/pkg/cache"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestRepo struct {
	data *Data
	log  *log.Helper
}

func (contestRepo *ContestRepo) GetContests(ctx context.Context) ([]*biz.Contest, error) {
	cacheBytes, err := contestRepo.data.GetCache(ctx, cache.ContestListCacheKey)
	// cache hit
	if err == nil && cacheBytes != nil {
		contestRepo.log.Infof("[cache] hit from user calling GetContests")
		var contests []*biz.Contest
		if err := json.Unmarshal(cacheBytes, &contests); err == nil {
			return contests, nil
		}
		contestRepo.log.Errorf("[cache] failed to unmarshal contest list from cache: %v", err)
	}

	// singleflight
	// 返回值 shared 表示结果是否被多个调用者共享，可用于监控
	res, err, _ := sfGroup.Do(cache.ContestListCacheKey, func() (interface{}, error) {
		// 必须超时控制
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// cache miss
		contestRepo.log.Infof("[cache] miss from user calling GetContests")

		queryContests, err := contestRepo.data.db.ProblemSet.Query().All(ctx)
		if err != nil {
			return nil, err
		}

		contests := make([]*biz.Contest, 0, len(queryContests))
		for i := 0; i < len(queryContests); i++ {
			contests = append(contests, &biz.Contest{
				ID:        queryContests[i].ID,
				Title:     queryContests[i].Name,
				StartTime: queryContests[i].StartTime,
				EndTime:   queryContests[i].EndTime,
				Status:    string(queryContests[i].Status),
			})
		}

		// set cache
		cacheBytes, err = json.Marshal(contests)
		if err != nil {
			contestRepo.log.Errorf("[cache] failed to marshal contest list for caching: %v", err)
		} else {
			err = contestRepo.data.SetCache(ctx, cache.ContestListCacheKey, cacheBytes, 5*time.Minute)
			if err != nil {
				contestRepo.log.Errorf("[cache] failed to set cache for contest list: %v", err)
			}
		}

		return contests, nil
	})

	if err != nil {
		return nil, err
	}
	return res.([]*biz.Contest), nil
}

func (contestRepo *ContestRepo) GetSingleContest(ctx context.Context, contestID int64) (*biz.DetailedContest, error) {
	cacheBytes, err := contestRepo.data.GetCache(ctx, fmt.Sprintf(cache.ContestDetailCacheKeyFmt, contestID))
	// cache hit
	if err == nil && cacheBytes != nil {
		contestRepo.log.Infof("[cache] hit from user calling GetSingleContest")
		var contest *biz.DetailedContest
		if err := json.Unmarshal(cacheBytes, &contest); err == nil {
			return contest, nil
		}
		contestRepo.log.Errorf("[cache] failed to unmarshal contest from cache: %v", err)
	}

	// singleflight
	// 返回值 shared 表示结果是否被多个调用者共享，可用于监控
	res, err, _ := sfGroup.Do(fmt.Sprintf(cache.ContestDetailCacheKeyFmt, contestID), func() (interface{}, error) {
		// 必须超时控制
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		// cache miss
		contestRepo.log.Infof("[cache] miss from user calling GetSingleContest")

		queryContest, err := contestRepo.data.db.ProblemSet.Query().
			Where(problemset.IDEQ(contestID)).
			Only(ctx)
		if err != nil {
			return nil, err
		}

		now := time.Now()
		calculatedStatus := CalculateContestStatus(now, queryContest.StartTime, queryContest.EndTime)
		if queryContest.Status.String() != calculatedStatus {
			queryContest.Status = problemset.Status(calculatedStatus)
			EnqueueStatusUpdate(queryContest.ID, calculatedStatus)
		}

		detailedContest := &biz.DetailedContest{
			Contest: biz.Contest{
				ID:        queryContest.ID,
				Title:     queryContest.Name,
				StartTime: queryContest.StartTime,
				EndTime:   queryContest.EndTime,
				Status:    string(queryContest.Status),
			},
			Description: queryContest.Description,
		}

		// set cache
		cacheBytes, err = json.Marshal(detailedContest)
		if err != nil {
			contestRepo.log.Errorf("[cache] failed to marshal contest for caching: %v", err)
		} else {
			err = contestRepo.data.SetCache(ctx, fmt.Sprintf(cache.ContestDetailCacheKeyFmt, contestID), cacheBytes, 5*time.Minute)
			if err != nil {
				contestRepo.log.Errorf("[cache] failed to set cache for contest %d: %v", contestID, err)
			}
		}

		return detailedContest, nil
	})

	if err != nil {
		return nil, err
	}
	return res.(*biz.DetailedContest), nil
}

func (contestRepo *ContestRepo) JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	// check if contest ended
	queryContest, err := contestRepo.data.db.ProblemSet.Query().
		Where(problemset.IDEQ(contestID)).
		Only(ctx)
	if err != nil {
		return false, err
	}
	now := time.Now()
	calculatedStatus := CalculateContestStatus(now, queryContest.StartTime, queryContest.EndTime)
	if queryContest.Status.String() != calculatedStatus {
		queryContest.Status = problemset.Status(calculatedStatus)
		EnqueueStatusUpdate(queryContest.ID, calculatedStatus)
	}
	if queryContest.Status == problemset.StatusEnded {
		contestRepo.log.Errorf("failed to join contest because contest has already ended")
		return false, nil
	}

	queryContestRecord, err := contestRepo.data.db.Competitor_List.Query().
		Where(competitor_list.And(competitor_list.ProblemSetIDEQ(contestID), competitor_list.UserIDEQ(userID))).
		Exist(ctx)
	if err != nil {
		return false, nil
	}
	if queryContestRecord {
		contestRepo.log.Errorf("failed to join contest because user has already joined the same contest")
		return true, nil
	}
	_, err = contestRepo.data.db.Competitor_List.Create().SetUserID(userID).SetProblemSetID(contestID).Save(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (contestRepo *ContestRepo) QuitContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	// check if contest started or ended
	queryContest, err := contestRepo.data.db.ProblemSet.Query().
		Where(problemset.IDEQ(contestID)).
		Only(ctx)
	if err != nil {
		contestRepo.log.Errorf("failed to quit contest because cannot find contest %d: %v", contestID, err)
		return false, err
	}
	now := time.Now()
	calculatedStatus := CalculateContestStatus(now, queryContest.StartTime, queryContest.EndTime)
	if queryContest.Status.String() != calculatedStatus {
		queryContest.Status = problemset.Status(calculatedStatus)
		EnqueueStatusUpdate(queryContest.ID, calculatedStatus)
	}
	if queryContest.Status == problemset.StatusOngoing || queryContest.Status == problemset.StatusEnded {
		contestRepo.log.Errorf("failed to quit contest because contest has already started or ended")
		return true, nil
	}

	queryContestRecord, err := contestRepo.data.db.Competitor_List.Query().
		Where(competitor_list.And(competitor_list.ProblemSetIDEQ(contestID), competitor_list.UserIDEQ(userID))).
		Exist(ctx)
	if err != nil {
		return true, nil
	}
	if !queryContestRecord {
		contestRepo.log.Errorf("failed to quit contest because user has not joined the contest yet")
		return false, nil
	}
	_, err = contestRepo.data.db.Competitor_List.Delete().
		Where(competitor_list.And(competitor_list.ProblemSetIDEQ(contestID), competitor_list.UserIDEQ(userID))).
		Exec(ctx)
	if err != nil {
		return true, err
	}
	return false, nil
}

func (contestRepo *ContestRepo) GetJoinStatus(ctx context.Context, contestID int64, userID int64) (bool, error) {
	isJoinedExist, err := contestRepo.data.db.Competitor_List.Query().
		Where(competitor_list.And(competitor_list.ProblemSetIDEQ(contestID), competitor_list.UserIDEQ(userID))).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	return isJoinedExist, nil
}

func NewContestRepo(data *Data, logger log.Logger) biz.ContestRepo {
	return &ContestRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// 延迟双删 Contest 列表或单个 Contest 缓存
func (contestRepo *ContestRepo) ContestCacheConsumer(ctx context.Context, msg *mq.ContestCacheMsg) error {
	var key string
	if msg.Scale == "single" {
		key = fmt.Sprintf(cache.ContestDetailCacheKeyFmt, msg.ContestID)
	} else {
		// handle as list cache by default
		key = cache.ContestListCacheKey
	}

	err := contestRepo.data.DeleteCache(ctx, key)
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to delete cache for contest %d: %v", *msg, err)
		return err
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		err := contestRepo.data.DeleteCache(ctx, key)
		if err != nil {
			contestRepo.log.Errorf("[cache] failed to delayed delete cache for contest %d: %v", *msg, err)
		}
	}()
	return nil
}
