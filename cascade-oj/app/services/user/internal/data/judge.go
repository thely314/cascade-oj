package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/judgerecord"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/problemset"
	"cascade-oj/ent/submissionrecord"
	"cascade-oj/pkg/cache"
	"cascade-oj/pkg/middleware/auth"
	"cascade-oj/pkg/mq"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *judgeRepo) CreateSelfTest(ctx context.Context, selfTest *biz.SelfTest) (string, error) {
	conn, err := repo.data.mq_channel.GetConnection()
	if err != nil {
		log.Errorf("failed to get mq connection: %v", err)
		return "", err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("failed opening a channel")
		return "", err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.GojudgeSelfTestQueueName, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
		log.Errorf("error from QueueDeclare: %v", err)
		return "", err
	}
	selfTestMsg := &mq.SelfTestMessage{
		UUID:       selfTest.UUID,
		UserID:     selfTest.UserID,
		ProblemID:  selfTest.ProblemID,
		Code:       selfTest.Code,
		Language:   selfTest.Language,
		Input:      selfTest.Input,
		IsCompiled: false,
		Stdout:     "",
		Stderr:     "",
		TimeCost:   0,
		MemoryCost: 0,
		Token:      "",
	}
	// 结构体 slice 转为 JSON
	jsonBody, err := json.Marshal(selfTestMsg)
	if err != nil {
		log.Errorf("error when marshaling selfTestMsg: %v", err)
		return "", err
	}
	log.Infof("publish message: %s", jsonBody)

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBody,
		},
	)
	if err != nil {
		log.Errorf("error from PublishWithContext: %v", err)
		return "", err
	}

	// set cache
	err = repo.data.SetCache(ctx, fmt.Sprintf(cache.SelfTestCacheKeyFmt, selfTest.UserID, selfTest.UUID), jsonBody, 2*time.Hour)
	if err != nil {
		repo.log.Errorf("[cache] failed to set cache for self test: %v", err)
		return "", err
	}
	return selfTest.UUID, nil
}

func (repo *judgeRepo) CreateSubmission(ctx context.Context, submission *biz.Submission) (string, error) {
	// check if user has joined the contest
	_, err := repo.data.db.Competitor_List.Query().
		Where(competitor_list.And(
			competitor_list.ProblemSetIDEQ(submission.ProblemSetID),
			competitor_list.UserIDEQ(submission.UserID),
		)).
		Only(ctx)
	if err != nil {
		return "", errors.New("user has not joined the contest")
	}

	// check if contest is ongoing
	queryContest, err := repo.data.db.ProblemSet.Query().
		Where(problemset.IDEQ(submission.ProblemSetID)).
		Only(ctx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	calculatedStatus := CalculateContestStatus(now, queryContest.StartTime, queryContest.EndTime)
	if queryContest.Status.String() != calculatedStatus {
		queryContest.Status = problemset.Status(calculatedStatus)
		EnqueueStatusUpdate(queryContest.ID, calculatedStatus)
	}
	if calculatedStatus != "ongoing" {
		return "", errors.New("contest is not ongoing")
	}

	conn, err := repo.data.mq_channel.GetConnection()
	if err != nil {
		log.Errorf("failed to get mq connection: %v", err)
		return "", err
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Errorf("failed opening a channel")
		return "", err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.GojudgeSubmissionQueueName, // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		log.Errorf("error from QueueDeclare: %v", err)
		return "", err
	}
	// get case version from problemTarget
	problemTarget, err := repo.data.db.Problem.Query().
		Select(problem.FieldCaseVersion).
		Where(problem.IDEQ(submission.ProblemID)).
		Only(ctx)
	if err != nil {
		return "", err
	}
	submissionMsg := &mq.SubmissionMessage{
		UUID:         submission.UUID,
		UserID:       submission.UserID,
		ProblemID:    submission.ProblemID,
		ProblemSetID: submission.ProblemSetID,
		Code:         submission.Code,
		Status:       0, // Pending
		Score:        0,
		CreateTime:   submission.CreateTime,
		TimeCost:     0,
		MemoryCost:   0,
		Language:     submission.Language,
		Stderr:       "",
		CaseVersion:  problemTarget.CaseVersion,
		Token:        "",
	}
	// 结构体 slice 转为 JSON
	jsonBody, err := json.Marshal(submissionMsg)
	if err != nil {
		log.Errorf("error when marshaling submissionMsg: %v", err)
		return "", err
	}
	log.Infof("publish message: %s", jsonBody)

	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonBody,
		},
	)
	if err != nil {
		log.Errorf("error from PublishWithContext: %v", err)
		return "", err
	}

	// delete cache for submission list
	repo.submissionListCacheUpdate(ctx, submission.UserID, submission.ProblemSetID, submission.ProblemID)
	// set cache
	err = repo.data.SetCache(ctx, fmt.Sprintf(cache.SubmissionDetailCacheKeyFmt, submission.UserID, submission.UUID), jsonBody, 2*time.Hour)
	if err != nil {
		repo.log.Errorf("[cache] failed to set cache for submission: %v", err)
		return "", err
	}
	return submission.UUID, nil
}

func (repo *judgeRepo) GetSelfTest(ctx context.Context, selfTestUUID string) (*biz.SelfTest, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID

	cacheBytes, err := repo.data.GetCache(ctx, fmt.Sprintf(cache.SelfTestCacheKeyFmt, userID, selfTestUUID))
	// cache hit
	if err == nil && cacheBytes != nil {
		repo.log.Infof("[cache] hit from user calling GetSelfTest")
		var selftest *biz.SelfTest
		if err := json.Unmarshal(cacheBytes, &selftest); err == nil {
			return selftest, nil
		}
		repo.log.Errorf("[cache] failed to unmarshal self test from cache: %v", err)
	}

	// cache miss, self test is not durable, no need to continue
	return nil, err
}

func (repo *judgeRepo) GetSubmissions(ctx context.Context, contestID int64, problemID int64) ([]*biz.Submission, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	cacheBytes, err := repo.data.GetCache(ctx, fmt.Sprintf(cache.SubmissionListCacheKeyFmt, userID, contestID, problemID))
	// cache hit
	if err == nil && cacheBytes != nil {
		repo.log.Infof("[cache] hit from user calling GetSubmissions")
		var submissions []*biz.Submission
		if err := json.Unmarshal(cacheBytes, &submissions); err == nil {
			return submissions, nil
		}
		repo.log.Errorf("[cache] failed to unmarshal submission list from cache: %v", err)
	}

	// cache miss
	repo.log.Infof("[cache] miss from user calling GetSubmissions")

	// verify problem set exists
	_, err = repo.data.db.ProblemSet.Get(ctx, contestID)
	if err != nil {
		return nil, err
	}

	// submission records don't carry user id directly, join via judge edge
	po, err := repo.data.db.SubmissionRecord.Query().
		Where(submissionrecord.ProblemSetIDEQ(contestID), submissionrecord.ProblemIDEQ(problemID), submissionrecord.HasJudgeWith(judgerecord.UserIDEQ(userID))).
		WithJudge().
		WithProblem().
		Order(ent.Desc(submissionrecord.FieldSubmissionTime)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var submissions []*biz.Submission
	for _, v := range po {
		submissions = append(submissions, &biz.Submission{
			UUID:        v.Edges.Judge.UUID,
			UserID:      v.Edges.Judge.UserID,
			ProblemID:   v.ProblemID,
			Code:        v.Edges.Judge.Code,
			Language:    v.Edges.Judge.Language,
			Status:      util.StatusToString(v.Edges.Judge.Status),
			CreateTime:  v.SubmissionTime,
			Score:       v.Score,
			TimeCost:    int64(v.Edges.Judge.TimeCostMs),
			MemoryCost:  int64(v.Edges.Judge.MemoryCostKB),
			CaseVersion: int8(v.Edges.Problem.CaseVersion),
		})
	}

	// set cache
	cacheBytes, err = json.Marshal(submissions)
	if err != nil {
		repo.log.Errorf("[cache] failed to marshal submission list for caching: %v", err)
	} else {
		err = repo.data.SetCache(ctx, fmt.Sprintf(cache.SubmissionListCacheKeyFmt, userID, contestID, problemID), cacheBytes, 5*time.Minute)
		if err != nil {
			repo.log.Errorf("[cache] failed to set cache for submission list: %v", err)
		}
	}

	return submissions, nil
}

func (repo *judgeRepo) GetSingleSubmission(ctx context.Context, submissionUUID string) (*biz.Submission, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	cacheBytes, err := repo.data.GetCache(ctx, fmt.Sprintf(cache.SubmissionDetailCacheKeyFmt, userID, submissionUUID))
	// cache hit
	if err == nil && cacheBytes != nil {
		repo.log.Infof("[cache] hit from user calling GetSingleSubmission")
		var submission *biz.Submission
		if err := json.Unmarshal(cacheBytes, &submission); err == nil {
			return submission, nil
		}
		repo.log.Errorf("[cache] failed to unmarshal submission from cache: %v", err)
	}

	// cache miss
	repo.log.Infof("[cache] miss from user calling GetSingleSubmission")

	var res *biz.Submission

	// get from db
	po, err := repo.data.db.SubmissionRecord.Query().
		Where(submissionrecord.HasJudgeWith(judgerecord.UUIDEQ(submissionUUID))).
		WithJudge().
		WithProblem().
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("no submission found: %s", submissionUUID)
	}
	// load judge to get user/code/language/status
	j, err := po.QueryJudge().Only(ctx)
	if err != nil {
		return nil, err
	}
	if j.UserID != userID {
		return nil, fmt.Errorf("permission denied")
	}
	res = &biz.Submission{
		UUID:        submissionUUID,
		UserID:      j.UserID,
		ProblemID:   po.ProblemID,
		Code:        j.Code,
		Language:    j.Language,
		Status:      util.StatusToString(j.Status),
		CreateTime:  po.SubmissionTime,
		Score:       po.Score,
		TimeCost:    int64(j.TimeCostMs),
		MemoryCost:  int64(j.MemoryCostKB),
		CaseVersion: int8(po.Edges.Problem.CaseVersion),
	}

	// set cache
	cacheBytes, err = json.Marshal(res)
	if err != nil {
		repo.log.Errorf("[cache] failed to marshal submission for caching: %v", err)
	} else {
		err = repo.data.SetCache(ctx, fmt.Sprintf(cache.SubmissionDetailCacheKeyFmt, userID, submissionUUID), cacheBytes, 5*time.Minute)
		if err != nil {
			repo.log.Errorf("[cache] failed to set cache for submission: %v", err)
		}
	}

	return res, nil
}

func (repo *judgeRepo) GetCases(ctx context.Context, submissionUUID string) ([]*biz.Case, error) {
	// get from db
	po, err := repo.data.db.SubmissionRecord.Query().
		Where(submissionrecord.HasJudgeWith(judgerecord.UUIDEQ(submissionUUID))).
		WithCaseGroupResults().
		Only(ctx)
	if err != nil {
		repo.log.Errorf("get cases failed: %v", err)
		return nil, errors.New("permission denied")
	}
	caseGroupResultList, err := po.QueryCaseGroupResults().All(ctx)
	if err != nil {
		return nil, err
	}
	caseResultList := make([]*ent.CaseResult, 0)
	for _, caseGroupResults := range caseGroupResultList {
		caseResults, err := caseGroupResults.QueryCaseResults().All(ctx)
		if err != nil {
			return nil, err
		}
		caseResultList = append(caseResultList, caseResults...)
	}

	cases := make([]*biz.Case, 0)
	for i, v := range caseResultList {
		cases = append(cases, &biz.Case{
			Index:      int32(i),
			Score:      int32(v.Score),
			Status:     util.StatusToString(v.Status),
			TimeCost:   v.TimeCostMs,
			MemoryCost: v.MemoryCostKB,
		})
	}
	return cases, nil
}

// 延迟双删 Submission 列表
func (repo *judgeRepo) submissionListCacheUpdate(ctx context.Context, userID, contestID, problemID int64) error {
	key := fmt.Sprintf(cache.SubmissionListCacheKeyFmt, userID, contestID, problemID)
	err := repo.data.DeleteCache(ctx, key)
	if err != nil {
		repo.log.Errorf("[cache] failed to delete cache for submission list: %v", err)
		return err
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		err := repo.data.DeleteCache(ctx, key)
		if err != nil {
			repo.log.Errorf("[cache] failed to delayed delete cache for submission list: %v", err)
		}
	}()
	return nil
}
