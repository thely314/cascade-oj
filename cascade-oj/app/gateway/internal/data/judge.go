package data

import (
	"context"
	"encoding/json"
	"time"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"
	"cascade-oj/pkg/mq"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

type selfTestDTO struct {
	UUID     string `json:"id,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
	Code     string `json:"code,omitempty"`
	Language string `json:"language,omitempty"`
	Input    string `json:"input,omitempty"`
	Token    string `json:"token,omitempty"`
}

type submissionDTO struct {
	UUID       string    `json:"id,omitempty"`
	UserID     int64     `json:"user_id,omitempty"`
	ProblemID  int64     `json:"problem_id,omitempty"`
	Code       string    `json:"code,omitempty"`
	Language   string    `json:"language,omitempty"`
	Status     string    `json:"status,omitempty"`
	CreateTime time.Time `json:"create_time"`
	Score      int64     `json:"score,omitempty"`
	TimeCost   int64     `json:"time_cost,omitempty"`
	MemoryCost int64     `json:"memory_cost,omitempty"`
	Token      string    `json:"token,omitempty"`
}

// create self test record
func (repo *judgeRepo) CreateSelfTest(ctx context.Context, selfTest *biz.SelfTest) (string, error) {
	q, err := repo.data.mq_channel.QueueDeclare(
		mq.GojudgeSelfTestQueueName, // name
		false,                       // durable
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
	// judge microservice is responsible for updating cache
	return selfTest.UUID, nil
}

// create submission record
func (repo *judgeRepo) CreateSubmission(ctx context.Context, submission *biz.Submission) (string, error) {
	q, err := repo.data.mq_channel.QueueDeclare(
		mq.GojudgeSubmissionQueueName, // name
		false,                         // durable
		false,                         // delete when unused
		false,                         // exclusive
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		return "", err
	}
	// TODO
	// cant get case version from problemTarget
	// problemTarget, err := repo.data.db.Problem.Query().
	// 	Select(problem.FieldCaseVersion).
	// 	Where(problem.IDEQ(submission.ProblemID)).
	// 	Only(ctx)
	// if err != nil {
	// 	return "", err
	// }
	submissionMsg := &mq.SubmissionMessage{
		UUID:         submission.UUID,
		UserID:       submission.UserID,
		ProblemID:    submission.ProblemID,
		ProblemSetID: 0, // TODO: need problem set id
		Code:         submission.Code,
		Status:       0, // Pending
		Score:        0,
		CreateTime:   submission.CreateTime,
		TimeCost:     0,
		MemoryCost:   0,
		Language:     submission.Language,
		Stderr:       "",
		CaseVersion:  0,
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
	// judge microservice is responsible for updating cache
	return submission.UUID, nil
}

func (repo *judgeRepo) GetSubmissions(ctx context.Context, userID int64, problemID int64) ([]*biz.SubmissionMetadata, error) {
	submissions, err := repo.data.grpcUserClient.GetSubmissions(ctx, &pb.GetSubmissionsRequest{
		UserId:    userID,
		ProblemId: problemID,
	})
	if err != nil {
		return nil, err
	}
	var res []*biz.SubmissionMetadata
	for _, v := range submissions.Submissions {
		res = append(res, &biz.SubmissionMetadata{
			SubmissionID: v.SubmissionUuid,
			ProblemID:    v.ProblemId,
			UserID:       v.UserId,
			Status:       v.Status,
			SubmitTime:   v.SubmitTime.AsTime(),
			Score:        int32(v.Score),
		})
	}
	return res, nil
}

func (repo *judgeRepo) GetSingleSubmission(ctx context.Context, submissionID string) (*biz.SingleSubmissionDTO, error) {
	submission, err := repo.data.grpcUserClient.GetSingleSubmission(ctx, &pb.GetSingleSubmissionRequest{
		SubmissionUuid: submissionID,
	})
	if err != nil {
		return nil, err
	}
	var casesResults []*biz.Case
	for i, v := range submission.CaseResults.Cases {
		casesResults = append(casesResults, &biz.Case{
			Index:      int32(i),
			Score:      v.Score,
			Status:     v.Status,
			TimeCost:   uint64(v.TimeCost),
			MemoryCost: uint64(v.MemoryCost),
		})
	}
	var res = &biz.SingleSubmissionDTO{
		Submission: &biz.Submission{
			UUID:       submissionID,
			UserID:     submission.Metadata.UserId,
			ProblemID:  submission.Metadata.ProblemId,
			Code:       submission.Code,
			Language:   submission.Language,
			Status:     submission.Metadata.Status,
			CreateTime: submission.Metadata.SubmitTime.AsTime(),
			Score:      int(submission.Metadata.Score),
			TimeCost:   int64(submission.TimeCost),
			MemoryCost: int64(submission.MemoryCost),
		},
		CasesResults: casesResults,
	}
	return res, nil
}

func (repo *judgeRepo) GetSelfTest(ctx context.Context, selfTestUUID string) (*biz.SelfTest, error) {
	selfTest, err := repo.data.grpcUserClient.GetSelfTestResult(ctx, &pb.GetSelfTestResultRequest{
		SelftestUuid: selfTestUUID,
	})
	if err != nil {
		return nil, err
	}
	var res = &biz.SelfTest{
		IsCompiled: selfTest.IsCompiled,
		Stdout:     selfTest.Stdout,
		Stderr:     selfTest.Stderr,
		TimeCost:   int64(selfTest.TimeCost),
		MemoryCost: int64(selfTest.MemoryCost),
	}
	return res, nil
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
