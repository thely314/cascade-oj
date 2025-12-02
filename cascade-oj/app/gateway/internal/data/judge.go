package data

import (
	"context"
	"encoding/json"
	"time"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

type selfTestDTO struct {
	ID       string `json:"id,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
	Code     string `json:"code,omitempty"`
	Language string `json:"language,omitempty"`
	Input    string `json:"input,omitempty"`
	Token    string `json:"token,omitempty"`
}

type submissionDTO struct {
	ID         string    `json:"id,omitempty"`
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
	// store into mq
	q, err := repo.data.mq_channel.QueueDeclare(
		"self_test_queue", // name
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
		ID:       selfTest.ID,
		UserID:   selfTest.UserID,
		Code:     selfTest.Code,
		Language: selfTest.Language,
		Input:    selfTest.Input,
		Token:    repo.data.cache.token,
	}
	// 结构体 slice 转为 JSON
	jsonBody, err := json.Marshal(sDTO)
	if err != nil {
		return "", err
	}
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
	// update cache if needed
	return selfTest.ID, nil
}

// create submission record
func (repo *judgeRepo) CreateSubmission(ctx context.Context, submission *biz.Submission) (string, error) {
	// store into mq
	q, err := repo.data.mq_channel.QueueDeclare(
		"submission_queue", // name
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
		ID:         submission.ID,
		UserID:     submission.UserID,
		ProblemID:  submission.ProblemID,
		Code:       submission.Code,
		Language:   submission.Language,
		Status:     submission.Status,
		CreateTime: submission.CreateTime,
		Score:      submission.Score,
		TimeCost:   submission.TimeCost,
		MemoryCost: submission.MemoryCost,
		Token:      repo.data.cache.token,
	}
	// 结构体 slice 转为 JSON
	jsonBody, err := json.Marshal(sDTO)
	if err != nil {
		return "", err
	}
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
	// update cache if needed
	return submission.ID, nil
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
			SubmissionID: v.SubmissionId,
			ProblemID:    v.ProblemId,
			UserID:       v.UserId,
			Status:       v.Status,
			SubmitTime:   v.SubmitTime.AsTime(),
			Score:        int32(v.Score),
		})
	}
	return res, nil
}

func (repo *judgeRepo) GetSingleSubmission(ctx context.Context, submissionID string) (*biz.Submission, error) {
	submission, err := repo.data.grpcUserClient.GetSingleSubmission(ctx, &pb.GetSingleSubmissionRequest{
		SubmissionId: submissionID,
	})
	if err != nil {
		return nil, err
	}
	var res = &biz.Submission{
		ID:         submissionID,
		UserID:     submission.Metadata.UserId,
		ProblemID:  submission.Metadata.ProblemId,
		Code:       submission.Code,
		Language:   submission.Language,
		Status:     submission.Metadata.Status,
		CreateTime: submission.Metadata.SubmitTime.AsTime(),
		Score:      int64(submission.Metadata.Score),
		// TODO refactor the proto and add these fields
		// TimeCost:   submission.TimeCost,
		// MemoryCost: submission.MemoryCost,
	}
	return res, nil
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
