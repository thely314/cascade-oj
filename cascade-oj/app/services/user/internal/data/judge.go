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
	q, err := repo.data.mq_channel.QueueDeclare(
		mq.GojudgeSelfTestQueueName, // name
		true,                        // durable
		false,                       // delete when unused
		false,                       // exclusive
		false,                       // no-wait
		nil,                         // arguments
	)
	if err != nil {
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
		return "", err
	}
	log.Infof("publish message: %s", jsonBody)
	err = repo.data.mq_channel.PublishWithContext(
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
		return "", err
	}
	set := repo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", mq.SelfTestType, selfTest.UserID, selfTest.UUID), jsonBody, 2*time.Hour)
	if set.Err() != nil {
		return "", set.Err()
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

	q, err := repo.data.mq_channel.QueueDeclare(
		mq.GojudgeSubmissionQueueName, // name
		true,                          // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
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
		return "", err
	}
	log.Infof("publish message: %s", jsonBody)
	err = repo.data.mq_channel.PublishWithContext(
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
		return "", err
	}
	set := repo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", mq.SubmissionType, submission.UserID, submission.UUID), jsonBody, 2*time.Hour)
	if set.Err() != nil {
		return "", set.Err()
	}
	return submission.UUID, nil
}

func (repo *judgeRepo) GetSelfTest(ctx context.Context, selfTestUUID string) (*biz.SelfTest, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	var selfTestMsg *mq.SelfTestMessage
	result, err := repo.data.redis.Get(ctx, fmt.Sprintf("%s:%d:%s", mq.SelfTestType, userID, selfTestUUID)).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(result), &selfTestMsg)
	if err != nil {
		return nil, err
	}
	selfTest := &biz.SelfTest{
		UUID:       selfTestMsg.UUID,
		UserID:     selfTestMsg.UserID,
		ProblemID:  selfTestMsg.ProblemID,
		Code:       selfTestMsg.Code,
		Language:   selfTestMsg.Language,
		Input:      selfTestMsg.Input,
		IsCompiled: selfTestMsg.IsCompiled,
		Stdout:     selfTestMsg.Stdout,
		Stderr:     selfTestMsg.Stderr,
		TimeCost:   int64(selfTestMsg.TimeCost),
		MemoryCost: int64(selfTestMsg.MemoryCost),
	}
	return selfTest, nil
}

func (repo *judgeRepo) GetSubmissions(ctx context.Context, contestID int64, problemID int64) ([]*biz.Submission, error) {
	// verify problem set exists
	_, err := repo.data.db.ProblemSet.Get(ctx, contestID)
	if err != nil {
		return nil, err
	}
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
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
	return submissions, nil
}

func (repo *judgeRepo) GetSingleSubmission(ctx context.Context, submissionUUID string) (*biz.Submission, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	var res *biz.Submission
	// try to get from redis by UUID
	var submissionMsg *mq.SubmissionMessage
	result, err := repo.data.redis.Get(ctx, fmt.Sprintf("%s:%d:%s", mq.SubmissionType, userID, submissionUUID)).Result()
	if err != nil {
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
		return res, nil
	}
	err = json.Unmarshal([]byte(result), &submissionMsg)
	if err != nil {
		return nil, err
	}
	res = &biz.Submission{
		UUID:        submissionMsg.UUID,
		UserID:      submissionMsg.UserID,
		ProblemID:   submissionMsg.ProblemID,
		Code:        submissionMsg.Code,
		Language:    submissionMsg.Language,
		Status:      util.StatusToString(submissionMsg.Status),
		CreateTime:  submissionMsg.CreateTime,
		Score:       submissionMsg.Score,
		TimeCost:    int64(submissionMsg.TimeCost),
		MemoryCost:  int64(submissionMsg.MemoryCost),
		CaseVersion: int8(submissionMsg.CaseVersion),
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
