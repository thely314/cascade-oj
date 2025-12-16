package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/judgerecord"
	"cascade-oj/ent/submissionrecord"
	"cascade-oj/pkg/middleware/auth"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

type selfTestDTO struct {
	UUID       string `json:"id,omitempty"`
	UserID     int64  `json:"user_id,omitempty"`
	Code       string `json:"code,omitempty"`
	Language   string `json:"language,omitempty"`
	Input      string `json:"input,omitempty"`
	IsCompiled bool   `json:"is_compiled,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	TimeCost   int64  `json:"time_cost,omitempty"`
	MemoryCost int64  `json:"memory_cost,omitempty"`
	Token      string `json:"token,omitempty"`
}

type submissionDTO struct {
	UUID        string    `json:"id,omitempty"`
	UserID      int64     `json:"user_id,omitempty"`
	ProblemID   int64     `json:"problem_id,omitempty"`
	Code        string    `json:"code,omitempty"`
	Language    string    `json:"language,omitempty"`
	Status      string    `json:"status,omitempty"`
	Score       int       `json:"score,omitempty"`
	CreateTime  time.Time `json:"create_time"`
	TimeCost    int64     `json:"time_cost,omitempty"`
	MemoryCost  int64     `json:"memory_cost,omitempty"`
	Stderr      string    `json:"stderr,omitempty"`
	CaseVersion int16     `json:"case_version,omitempty"`
	Token       string    `json:"token,omitempty"`
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (repo *judgeRepo) CreateSelfTest(ctx context.Context, selfTest *biz.SelfTest) (string, error) {
	_, err := repo.data.redis.Get(ctx, "user:contest:problem:"+strconv.FormatInt(selfTest.ProblemID, 10)).Result()
	if err != nil {
		return "", err
	}
	q, err := repo.data.mq_channel.QueueDeclare(
		"self_test_queue", // TODO: name may be wrong
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return "", err
	}
	sDTO := &selfTestDTO{
		UUID:       selfTest.UUID,
		UserID:     selfTest.UserID,
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
	jsonBody, err := json.Marshal(sDTO)
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
	return selfTest.UUID, nil
}

func (repo *judgeRepo) CreateSubmission(ctx context.Context, submission *biz.Submission) (string, error) {
	_, err := repo.data.redis.Get(ctx, "user:contest:problem:"+strconv.FormatInt(submission.ProblemID, 10)).Result()
	if err != nil {
		return "", err
	}
	q, err := repo.data.mq_channel.QueueDeclare(
		"submission_queue", // TODO: name may be wrong
		false,              // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return "", err
	}
	sDTO := &submissionDTO{
		UUID:       submission.UUID,
		UserID:     submission.UserID,
		ProblemID:  submission.ProblemID,
		Code:       submission.Code,
		Language:   submission.Language,
		Status:     "Waiting",
		Score:      0,
		CreateTime: submission.CreateTime,
		TimeCost:   0,
		MemoryCost: 0,
		Token:      "",
	}
	// 结构体 slice 转为 JSON
	jsonBody, err := json.Marshal(sDTO)
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
	set := repo.data.redis.Set(ctx, fmt.Sprintf("submission:%d:%s", submission.UserID, submission.UUID), jsonBody, 2*time.Hour)
	if set.Err() != nil {
		return "", set.Err()
	}
	return submission.UUID, nil
}

func (repo *judgeRepo) GetSelfTest(ctx context.Context, selfTestID string) (*biz.SelfTest, error) {
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	var selfTest *biz.SelfTest
	var sDTO *selfTestDTO
	result, err := repo.data.redis.Get(ctx, fmt.Sprintf("self-test:%d:%s", userID, selfTestID)).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(result), &sDTO)
	if err != nil {
		return nil, err
	}
	selfTest = &biz.SelfTest{
		UUID:       sDTO.UUID,
		UserID:     sDTO.UserID,
		Code:       sDTO.Code,
		Input:      sDTO.Input,
		IsCompiled: sDTO.IsCompiled,
		Stdout:     sDTO.Stdout,
		Stderr:     sDTO.Stderr,
		TimeCost:   sDTO.TimeCost,
		MemoryCost: sDTO.MemoryCost,
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
		Order(ent.Desc(submissionrecord.FieldSubmissionTime)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	var submissions []*biz.Submission
	for _, v := range po {
		var j *ent.JudgeRecord
		if v.Edges.Judge != nil {
			j = v.Edges.Judge
		}
		var language string
		var code string
		var uid int64
		var status string
		if j != nil {
			language = string(j.Language)
			code = j.Code
			uid = j.UserID
			status = j.Status.String()
		}
		submissions = append(submissions, &biz.Submission{
			UUID:       strconv.FormatInt(v.ID, 10),
			UserID:     uid,
			ProblemID:  v.ProblemID,
			Code:       code,
			Language:   language,
			Status:     status,
			Score:      v.Score,
			CreateTime: v.SubmissionTime,
		})
	}
	return submissions, nil
}

func (repo *judgeRepo) GetSingleSubmission(ctx context.Context, submissionID string) (*biz.Submission, error) {
	id, err := strconv.ParseInt(submissionID, 10, 64)
	userID := ctx.Value("userInfo").(*auth.Claims).UserID
	var res *biz.Submission
	if err != nil {
		// try to get from redis by UUID
		var sDTO submissionDTO
		result, err := repo.data.redis.Get(ctx, fmt.Sprintf("submission:%d:%s", userID, submissionID)).Result()
		if err != nil {
			return nil, fmt.Errorf("no submission found: %s", submissionID)
		}
		err = json.Unmarshal([]byte(result), &sDTO)
		if err != nil {
			return nil, err
		}
		res = &biz.Submission{
			UUID:       sDTO.UUID,
			UserID:     sDTO.UserID,
			ProblemID:  sDTO.ProblemID,
			Code:       sDTO.Code,
			Language:   sDTO.Language,
			Status:     sDTO.Status,
			CreateTime: sDTO.CreateTime,
			Score:      sDTO.Score,
		}
		return res, nil
	}
	// get from db
	po, err := repo.data.db.SubmissionRecord.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("no submission found: %s", submissionID)
	}
	// load judge to get user/code/language/status
	j, err := po.QueryJudge().Only(ctx)
	if err != nil {
		return nil, err
	}
	if j.UserID != userID {
		return nil, errors.New("permission denied")
	}
	res = &biz.Submission{
		UUID:       strconv.FormatInt(po.ID, 10),
		UserID:     j.UserID,
		ProblemID:  po.ProblemID,
		Code:       j.Code,
		Language:   string(j.Language),
		Status:     j.Status.String(),
		CreateTime: po.SubmissionTime,
		Score:      po.Score,
	}
	return res, nil
}

func (repo *judgeRepo) GetCases(ctx context.Context, submissionID string) ([]*biz.Case, error) {
	// id, err := strconv.ParseInt(submissionID, 10, 64)
	// userID := ctx.Value("userInfo").(*auth.Claims).UserID
	var cases []*biz.Case
	// get from db
	// TODO: apply ent CaseGroupResult to support dbQuery
	// po, err := repo.data.db.CaseGroupResult.
	// 	Query().
	// 	Where(CaseGroupResult.SubmissionIDEQ(id)).
	// 	WithSubmissions(func(q *ent.CaseGroupResultQuery) {
	// 		q.Where(submission.UserIDEQ(userID))
	// 	}).
	// 	All(ctx)
	// if err != nil {
	// 	repo.log.Errorf("get cases failed: %v", err)
	// 	return nil, errors.New("permission denied")
	// }
	cases = make([]*biz.Case, 0)
	// for i, v := range po {
	// 	cases = append(cases, &biz.Case{
	// 		Index:  int32(i),
	// 		Score:  int32(v.score),
	// 		Status:  int32(v.Status),
	// 		TimeCost:   v.TimeCost,
	// 		MemoryCost: v.MemoryCost,
	// 	})
	// }
	return cases, nil
}
