package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/casegroupresult"
	"cascade-oj/ent/judgerecord"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/submissionrecord"
	"cascade-oj/pkg/mq"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

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

type SubmissionRepo struct {
	data *Data
	log  *log.Helper
}

func NewSubmissionRepo(data *Data, logger log.Logger) biz.SubmissionRepo {
	return &SubmissionRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (submissionRepo *SubmissionRepo) GetSubmissions(ctx context.Context, request biz.SubmissionsRequestInfo) ([]*biz.Submission, error) {
	entSubmissions, err := submissionRepo.data.db.SubmissionRecord.
		Query().
		Select(
			submissionrecord.FieldProblemID,
			submissionrecord.FieldSubmissionTime,
			submissionrecord.FieldScore,
		).
		Where(
			submissionrecord.ProblemIDEQ(request.ProblemID),
			submissionrecord.ProblemSetIDEQ(request.ContestID),
			submissionrecord.HasJudgeWith(
				judgerecord.UserIDEQ(request.UserID),
			),
		).
		WithJudge(func(jrq *ent.JudgeRecordQuery) {
			jrq.Select(
				judgerecord.FieldUUID,
				judgerecord.FieldUserID,
				judgerecord.FieldStatus,
			)
		}).
		Offset(int((max(request.Page-1, 0) * request.PageSize))).
		Limit(int(request.PageSize)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	submissions := make([]*biz.Submission, 0, len(entSubmissions))
	for i := 0; i < len(entSubmissions); i++ {
		submissions = append(submissions,
			&biz.Submission{
				SubmissionUUID: entSubmissions[i].Edges.Judge.UUID,
				ProblemID:      entSubmissions[i].ProblemID,
				UserID:         entSubmissions[i].Edges.Judge.UserID,
				Status:         util.StatusToString(entSubmissions[i].Edges.Judge.Status),
				SubmitTime:     entSubmissions[i].SubmissionTime,
				Score:          int32(entSubmissions[i].Score),
			})
	}
	return submissions, nil
}

func (submissionRepo *SubmissionRepo) GetSingleSubmission(ctx context.Context, submissionUUID string) (*biz.DetailedSubmission, error) {
	entSubmission, err := submissionRepo.data.db.SubmissionRecord.
		Query().
		Select(
			submissionrecord.FieldProblemID,
			submissionrecord.FieldSubmissionTime,
			submissionrecord.FieldScore,
		).
		Where(
			submissionrecord.HasJudgeWith(
				judgerecord.UUIDEQ(submissionUUID),
			),
		).
		WithJudge(
			func(jrq *ent.JudgeRecordQuery) {
				jrq.Select(
					judgerecord.FieldUUID,
					judgerecord.FieldUserID,
					judgerecord.FieldStatus,
					judgerecord.FieldCode,
					judgerecord.FieldLanguage,
					judgerecord.FieldTimeCostMs,
					judgerecord.FieldMemoryCostKB,
				)
			}).
		WithCaseGroupResults(
			func(cgrq *ent.CaseGroupResultQuery) {
				cgrq.Select(
					casegroupresult.FieldScore,
					casegroupresult.FieldStatus,
					casegroupresult.FieldTotalTimeCostMs,
					casegroupresult.FieldMaxMemoryCostKB,
				)
			},
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	testCases := make([]*biz.TestCase, 0, len(entSubmission.Edges.CaseGroupResults))
	for i := 0; i < len(entSubmission.Edges.CaseGroupResults); i++ {
		testCases = append(testCases,
			&biz.TestCase{
				Score:      int32(entSubmission.Edges.CaseGroupResults[i].Score),
				Status:     util.StatusToString(entSubmission.Edges.CaseGroupResults[i].Status),
				TimeCost:   int32(entSubmission.Edges.CaseGroupResults[i].TotalTimeCostMs),
				MemoryCost: int32(entSubmission.Edges.CaseGroupResults[i].MaxMemoryCostKB),
			},
		)
	}

	return &biz.DetailedSubmission{
		Submission: biz.Submission{
			SubmissionUUID: entSubmission.Edges.Judge.UUID,
			ProblemID:      entSubmission.ProblemID,
			UserID:         entSubmission.Edges.Judge.UserID,
			Status:         util.StatusToString(entSubmission.Edges.Judge.Status),
			SubmitTime:     entSubmission.SubmissionTime,
			Score:          int32(entSubmission.Score),
		},
		Code:       entSubmission.Edges.Judge.Code,
		Language:   entSubmission.Edges.Judge.Language,
		TimeCost:   int32(entSubmission.Edges.Judge.TimeCostMs),
		MemoryCost: int32(entSubmission.Edges.Judge.MemoryCostKB),
		TestCases:  testCases,
	}, nil
}

func (submissionRepo *SubmissionRepo) RejudgeSubmission(ctx context.Context, submissionUUID string) (string, error) {
	// get from redis first
	// TODO get userID
	// submissionMsg, err := submissionRepo.data.redis.Get(ctx, fmt.Sprintf("%s:%d:%s", mq.SubmissionType, userID, submissionUUID)).Result()
	// if err != nil {
	// 	// get from db
	// }
	entSubmission, err := submissionRepo.data.db.SubmissionRecord.
		Query().
		Select(
			submissionrecord.FieldProblemID,
			submissionrecord.FieldProblemSetID,
			submissionrecord.FieldSubmissionTime,
			submissionrecord.FieldScore,
		).
		Where(
			submissionrecord.HasJudgeWith(
				judgerecord.UUIDEQ(submissionUUID),
			),
		).
		WithJudge(
			func(jrq *ent.JudgeRecordQuery) {
				jrq.Select(
					judgerecord.FieldUUID,
					judgerecord.FieldUserID,
					judgerecord.FieldStatus,
					judgerecord.FieldCode,
					judgerecord.FieldLanguage,
					judgerecord.FieldTimeCostMs,
					judgerecord.FieldMemoryCostKB,
				)
			}).
		WithCaseGroupResults(
			func(cgrq *ent.CaseGroupResultQuery) {
				cgrq.Select(
					casegroupresult.FieldScore,
					casegroupresult.FieldStatus,
					casegroupresult.FieldTotalTimeCostMs,
					casegroupresult.FieldMaxMemoryCostKB,
				)
			},
		).
		Only(ctx)
	if err != nil {
		return "", err
	}
	q, err := submissionRepo.data.mq_channel.QueueDeclare(
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
	problemTarget, err := submissionRepo.data.db.Problem.Query().
		Select(problem.FieldCaseVersion).
		Where(problem.IDEQ(entSubmission.ProblemID)).
		Only(ctx)
	if err != nil {
		return "", err
	}
	newUUID := uuid.New().String()
	submissionMsg := &mq.SubmissionMessage{
		UUID:         newUUID,
		UserID:       entSubmission.Edges.Judge.UserID,
		ProblemID:    entSubmission.ProblemID,
		ProblemSetID: entSubmission.ProblemSetID,
		Code:         entSubmission.Edges.Judge.Code,
		Status:       0, // Pending
		Score:        0,
		CreateTime:   time.Now(),
		TimeCost:     0,
		MemoryCost:   0,
		Language:     entSubmission.Edges.Judge.Language,
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
	err = submissionRepo.data.mq_channel.PublishWithContext(
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
	// TODO
	set := submissionRepo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", mq.SubmissionType, submissionMsg.UserID, submissionMsg.UUID), jsonBody, 2*time.Hour)
	if set.Err() != nil {
		return "", set.Err()
	}
	return newUUID, nil
}
