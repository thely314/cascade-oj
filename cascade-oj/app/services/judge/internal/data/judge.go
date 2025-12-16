package data

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"cascade-oj/app/services/judge/internal/biz"
	"cascade-oj/app/services/judge/pkg/gojudge"
	"cascade-oj/ent"
	"cascade-oj/ent/problem"
	"cascade-oj/pkg/mq"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

// judge types
const (
	SubmissionType = "test_case"
	SelfTestType   = "custom_test_case"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

// arg msg_submission comes with status = pending(0 in i16)
func (repo *judgeRepo) JudgeSubmission(ctx context.Context, msg_submission *mq.SubmissionMessage) error {
	// set submission to judging status
	msg_submission.Status = gojudge.StatusToInt16(gojudge.Judging)
	// set submission meta info
	msg_submission.MemoryCost = 0
	msg_submission.TimeCost = 0
	msg_marshal, err := json.Marshal(msg_submission)
	if err != nil {
		return err
	}
	// update redis cache
	repo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", SubmissionType, msg_submission.UserID, msg_submission.UUID), msg_marshal, 1*time.Hour)

	// default error status, recover when judge completed
	msg_submission.Status = gojudge.StatusToInt16(gojudge.SystemError)

	// update cache after judge completed
	defer func() {
		msg_marshal, err := json.Marshal(msg_submission)
		if err != nil {
			repo.log.Errorf("failed to marshal submission message: %v", err)
		}
		err = repo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", SubmissionType, msg_submission.UserID, msg_submission.UUID), msg_marshal, 1*time.Hour).Err()
		if err != nil {
			repo.log.Errorf("failed to update submission message in redis: %v", err)
		}
	}()

	// get problem data from cache or database
	var problemResult *ent.Problem
	problemString, err := repo.data.redis.Get(ctx, fmt.Sprintf("problem:%d", msg_submission.ProblemID)).Result()
	if err != nil {
		// get from database
		problemResult, err = repo.data.db.Problem.Query().
			Where(problem.IDEQ(msg_submission.ProblemID)).
			First(ctx)
		if err != nil {
			return err
		}
		// update cache
		problemBytes, err := json.Marshal(problemResult)
		if err != nil {
			return err
		}
		err = repo.data.redis.Set(ctx, fmt.Sprintf("problem:%d", msg_submission.ProblemID), problemBytes, 1*time.Hour).Err()
		if err != nil {
			return err
		}
	} else {
		_ = json.Unmarshal([]byte(problemString), &problemResult)
	}
	// set case version
	msg_submission.CaseVersion = problemResult.CaseVersion

	// fileManager get judge config
	judgeConfig, err := repo.data.fileManager.GetJudgeConfig(msg_submission.ProblemID)
	if err != nil {
		return err
	}

	// setup builders arrays to construct case groups in memory
	// builders will be updated during judge process
	// later commit each case result to db in one transaction
	// reduce interaction times with db
	caseGroupBuilders := make([]*ent.CaseGroupResultCreate, len(judgeConfig.CaseGroups))
	casesArrays := make([][]*ent.CaseResultCreate, len(judgeConfig.CaseGroups))

	// save submission to db after judge
	defer func() {
		judgeRecord, err := repo.data.db.JudgeRecord.Create().
			SetProblemID(msg_submission.ProblemID).
			SetUserID(msg_submission.UserID).
			SetUUID(msg_submission.UUID).
			SetStatus(msg_submission.Status).
			SetJudgeStartTime(time.Now()).
			SetTimeCostMs(msg_submission.TimeCost).
			SetMemoryCostKB(msg_submission.MemoryCost).
			SetCode(msg_submission.Code).
			SetLanguage(msg_submission.Language).
			SetJudgeType(SubmissionType).
			Save(ctx)
		if err != nil {
			repo.log.Errorf("failed to save judge record to db: %v", err)
			return
		}

		submissionRecord, err := repo.data.db.SubmissionRecord.Create().
			SetJudgeID(judgeRecord.ID).
			SetProblemID(msg_submission.ProblemID).
			SetProblemSetID(msg_submission.ProblemSetID).
			SetSubmissionTime(msg_submission.CreateTime).
			SetScore(int(msg_submission.Score)).
			Save(ctx)
		if err != nil {
			repo.log.Errorf("failed to save submission record to db: %v", err)
			return
		}

		// skip save case if system error or compile error
		if msg_submission.Status == gojudge.StatusToInt16(gojudge.SystemError) ||
			msg_submission.Status == gojudge.StatusToInt16(gojudge.CompileError) {
			return
		}

		// save case groups to db in one transaction
		for _, builder := range caseGroupBuilders {
			builder.SetSubmissionID(submissionRecord.ID)
		}
		caseGroups, err := repo.data.db.CaseGroupResult.CreateBulk(caseGroupBuilders...).Save(ctx)
		if err != nil {
			repo.log.Errorf("failed to save case group results to db: %v", err)
			return
		}

		// save cases to db in one transaction
		caseBuilders := make([]*ent.CaseResultCreate, 0)
		for idx, casesArray := range casesArrays {
			for _, builder := range casesArray {
				builder.SetCaseGroupResultID(caseGroups[idx].ID)
				caseBuilders = append(caseBuilders, builder)
			}
		}
		_, err = repo.data.db.CaseResult.CreateBulk(caseBuilders...).Save(ctx)
		if err != nil {
			repo.log.Errorf("failed to save case results to db: %v", err)
		}
	}()

	// compile
	fileID, result, err := repo.data.gojudge.Compile([]byte(msg_submission.Code), msg_submission.Language, msg_submission.UUID)
	if err != nil {
		if result != nil {
			msg_submission.Status = gojudge.StatusToInt16(gojudge.CompileError)
			// default stderr message
			msg_submission.Stderr = "Compile error without stderr message"
			stderr, is_ok := result.Files["stderr"]
			if is_ok {
				msg_submission.Stderr = string(stderr)
			}
		}
		return err
	}

	// delete compiled files after judge
	defer func() {
		err := repo.data.gojudge.DeleteFile(fileID)
		if err != nil {
			repo.log.Errorf("failed to delete compiled file: %v", err)
		}
	}()

	// begin judge
	// pre set status to accepted, change when any case failed
	msg_submission.Status = gojudge.StatusToInt16(gojudge.Accepted)

	// test each case
	// fileManager get config, read caseGroup from config
	for idx, caseGroup := range judgeConfig.CaseGroups {
		var caseGroupStatus int16 = gojudge.StatusToInt16(gojudge.Accepted)
		var caseGroupScore int = 0
		var totalTime uint64 = 0
		var maxMemory uint64 = 0
		caseBuilders := make([]*ent.CaseResultCreate, len(caseGroup.Cases))

		caseGroupBuilder := repo.data.db.CaseGroupResult.Create().
			SetStatus(gojudge.StatusToInt16(gojudge.Pending)).
			SetTotalTimeCostMs(0).
			SetMaxMemoryCostKB(0).
			SetScore(0)

		wg := sync.WaitGroup{}
		wg.Add(len(caseGroup.Cases))
		for idx2, testCase := range caseGroup.Cases {
			err := func() error {
				caseBuilder := repo.data.db.CaseResult.Create().
					SetStatus(gojudge.StatusToInt16(gojudge.Pending)).
					SetTimeCostMs(0).
					SetMemoryCostKB(0).
					SetScore(0)

				// fill builder into array after judge
				defer func() {
					caseBuilders[idx2] = caseBuilder
					wg.Done()
				}()

				// get case input, answer bytes
				inputBytes, ansBytes, err := repo.data.fileManager.GetCase(
					msg_submission.ProblemID,
					testCase.InputFileLocation,
					testCase.AnswerFileLocation,
				)
				if err != nil {
					return err
				}
				// TODO
				// handle CRLF for both files

				// judge case
				result, err := repo.data.gojudge.CaseJudge(msg_submission.ProblemID,
					testCase.InputFileLocation,
					inputBytes,
					msg_submission.Language,
					fileID,
					uuid.New().String(), // new uuid for each case
					uint64(judgeConfig.TimeResourceLimit),
					uint64(judgeConfig.TimeResourceLimit*2),
					uint64(judgeConfig.MemoryResourceLimit),
				)
				if err != nil {
					repo.log.Errorf("submission: %s runtime error: %v", msg_submission.UUID, err)
					caseBuilder.SetStatus(gojudge.StatusToInt16(gojudge.RuntimeError))
					caseBuilder.SetStderr(util.ParseSignalToString(result.ExitStatus))
					// the whole group is runtime error
					caseGroupStatus = gojudge.StatusToInt16(gojudge.RuntimeError)
					return nil
				}
				repo.log.Infof("submission: %s case: %s result: %+v", msg_submission.UUID, testCase.InputFileLocation, result)

				// check result status
				resStatus := gojudge.ParseGojudgeStatus(result.Status)
				out := result.Files["stdout"]
				if resStatus == gojudge.StatusToInt16(gojudge.Accepted) && out != nil {
					// compare output
					// TODO need util impl
					// if util...
					if !bytes.Equal(out, ansBytes) {
						resStatus = gojudge.StatusToInt16(gojudge.WrongAnswer)
					}
				}

				// set caseBuilder metadata
				caseBuilder.SetStatus(resStatus)
				caseBuilder.SetTimeCostMs(result.Time)
				caseBuilder.SetMemoryCostKB(result.Memory)

				// update case group info
				totalTime += result.Time
				maxMemory = max(maxMemory, result.Memory)

				// update submission info
				msg_submission.TimeCost += result.Time
				msg_submission.MemoryCost = max(msg_submission.MemoryCost, result.Memory)

				// set case score
				if resStatus != gojudge.StatusToInt16(gojudge.Accepted) {
					caseBuilder.SetScore(0)
					caseGroupStatus = resStatus
				} else {
					caseBuilder.SetScore(testCase.SubScore)
					caseGroupScore += testCase.SubScore
				}

				return nil
			}()
			if err != nil {
				msg_submission.Status = gojudge.StatusToInt16(gojudge.SystemError)
				return err
			}
		}
		wg.Wait()

		// calculate case group score
		caseGroupBuilder.SetScore(caseGroupScore)
		caseGroupBuilder.SetStatus(caseGroupStatus)
		caseGroupBuilder.SetTotalTimeCostMs(totalTime)
		caseGroupBuilder.SetMaxMemoryCostKB(maxMemory)
		caseGroupBuilders[idx] = caseGroupBuilder

		// fill cases array
		casesArrays[idx] = caseBuilders

		// update submission info
		msg_submission.Score += caseGroupScore
		if msg_submission.Status == gojudge.StatusToInt16(gojudge.Accepted) &&
			caseGroupStatus != gojudge.StatusToInt16(gojudge.Accepted) {
			msg_submission.Status = caseGroupStatus
		}
	}
	return nil
}

func (repo *judgeRepo) JudgeSelfTest(ctx context.Context, msg_self_test *mq.SelfTestMessage) error {
	// set self-test meta info
	msg_self_test.IsCompiled = false
	msg_self_test.Stdout = ""
	msg_self_test.Stderr = ""
	msg_self_test.TimeCost = 0
	msg_self_test.MemoryCost = 0

	defer func() {
		msg_marshal, err := json.Marshal(msg_self_test)
		if err != nil {
			repo.log.Errorf("failed to marshal self-test message: %v", err)
			return
		}
		err = repo.data.redis.Set(ctx, fmt.Sprintf("%s:%d:%s", SelfTestType, msg_self_test.UserID, msg_self_test.UUID), msg_marshal, 1*time.Hour).Err()
		if err != nil {
			repo.log.Errorf("failed to update redis cache: %v", err)
		}
	}()

	// get self-test config
	commands := *repo.data.gojudge.Commands
	selfTestConfig, is_ok := commands[msg_self_test.Language]
	if !is_ok {
		return errors.New("unsupported language: " + msg_self_test.Language)
	}

	// compile
	fileID, result, err := repo.data.gojudge.Compile([]byte(msg_self_test.Code), msg_self_test.Language, msg_self_test.UUID)
	if err != nil {
		if result != nil {
			// default stderr message
			msg_self_test.Stderr = "Compile error without stderr message"
			stderr, is_ok := result.Files["stderr"]
			if is_ok {
				msg_self_test.Stderr = string(stderr)
			}
		}
		return err
	}
	msg_self_test.IsCompiled = true

	// delete compiled files after judge
	defer func() {
		err := repo.data.gojudge.DeleteFile(fileID)
		if err != nil {
			repo.log.Errorf("failed to delete compiled file: %v", err)
		}
	}()

	// begin self-test judge
	result, err = repo.data.gojudge.CaseJudge(
		-1,
		"",
		[]byte(msg_self_test.Input),
		msg_self_test.Language,
		fileID,
		msg_self_test.UUID,
		selfTestConfig.RunConfig.CpuTimeLimit,
		selfTestConfig.CompileConfig.ClockTimeLimit,
		selfTestConfig.RunConfig.MemoryLimit,
	)
	if err != nil {
		return err
	}

	// update self-test info
	msg_self_test.Stdout = string(result.Files["stdout"])
	msg_self_test.Stderr = string(result.Files["stderr"])
	msg_self_test.TimeCost = result.Time
	msg_self_test.MemoryCost = result.Memory
	return nil
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
