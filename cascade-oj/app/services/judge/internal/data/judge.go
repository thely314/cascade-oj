package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"cascade-oj/app/services/judge/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/judgerecord"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/submissionrecord"
	"cascade-oj/pkg/cache"
	filemanage "cascade-oj/pkg/file_manage"
	"cascade-oj/pkg/mq"
	"cascade-oj/pkg/util"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
)

type judgeRepo struct {
	data *Data
	log  *log.Helper
}

// arg msg_submission comes with status = pending(0 in i16)
func (repo *judgeRepo) JudgeSubmission(ctx context.Context, msg_submission *mq.SubmissionMessage) error {
	// set submission to judging status
	msg_submission.Status = util.StatusToInt16(util.Judging)
	// set submission meta info
	msg_submission.MemoryCost = 0
	msg_submission.TimeCost = 0
	msg_marshal, err := json.Marshal(msg_submission)
	if err != nil {
		return err
	}
	// update redis cache
	repo.data.redis.Set(ctx, fmt.Sprintf(cache.SubmissionDetailCacheKeyFmt, msg_submission.UserID, msg_submission.UUID), msg_marshal, 1*time.Hour)

	// default error status, recover when judge completed
	msg_submission.Status = util.StatusToInt16(util.SystemError)

	// update cache after judge completed
	defer func() {
		msg_marshal, err := json.Marshal(msg_submission)
		if err != nil {
			repo.log.Errorf("failed to marshal submission message: %v", err)
		}
		err = repo.data.redis.Set(ctx, fmt.Sprintf(cache.SubmissionDetailCacheKeyFmt, msg_submission.UserID, msg_submission.UUID), msg_marshal, 1*time.Hour).Err()
		if err != nil {
			repo.log.Errorf("failed to update submission message in redis: %v", err)
		}
	}()

	// get problem data from cache or database
	var problemResult *ent.Problem
	problemString, err := repo.data.redis.Get(ctx, fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, msg_submission.ProblemID)).Result()
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
		err = repo.data.redis.Set(ctx, fmt.Sprintf(cache.ProblemDetailCacheKeyFmt, msg_submission.ProblemID), problemBytes, 1*time.Hour).Err()
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
			SetJudgeType(mq.SubmissionType).
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

		// update ranks
		// get all submissions from this user
		// then calculate total rank score
		// each submission group by the same problem should only count the highest score
		var existingJudges []struct {
			ProblemID int64 `json:"problem_id"`
			Max       int   `json:"max"`
		}
		err = repo.data.db.SubmissionRecord.Query().
			Where(submissionrecord.And(
				submissionrecord.HasJudgeWith(judgerecord.UserIDEQ(msg_submission.UserID)),
				submissionrecord.ProblemSetIDEQ(msg_submission.ProblemSetID),
			)).
			GroupBy(submissionrecord.FieldProblemID).
			Aggregate(ent.Max(submissionrecord.FieldScore)).
			Scan(ctx, &existingJudges)
		if err != nil {
			repo.log.Errorf("failed to query existing submission records: %v", err)
			return
		}
		var newRankScore int = 0
		for _, record := range existingJudges {
			newRankScore += record.Max
		}

		// calculate new total rank score
		oldRank, err := repo.data.db.Competitor_List.Query().
			Where(competitor_list.UserIDEQ(msg_submission.UserID)).
			First(ctx)
		if err != nil {
			repo.log.Errorf("failed to query competitor list: %v", err)
			return
		}
		if newRankScore > oldRank.TotalScore {
			// update ranking info in competitor list
			_, err = repo.data.db.Competitor_List.Update().
				Where(competitor_list.UserIDEQ(msg_submission.UserID)).
				SetTotalScore(newRankScore).
				Save(ctx)
			if err != nil {
				repo.log.Errorf("failed to update competitor list: %v", err)
				return
			}
		}

		// skip save case if system error or compile error
		if msg_submission.Status == util.StatusToInt16(util.SystemError) ||
			msg_submission.Status == util.StatusToInt16(util.CompileError) {
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
			msg_submission.Status = util.StatusToInt16(util.CompileError)
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
	msg_submission.Status = util.StatusToInt16(util.Accepted)

	// test each case
	// fileManager get config, read caseGroup from config
	for idx, caseGroup := range judgeConfig.CaseGroups {
		// 需要原子操作更新
		var caseGroupStatus int16 = util.StatusToInt16(util.Accepted)
		var atomicCaseGroupStatus int32 = int32(util.StatusToInt16(util.Accepted))
		var caseGroupScore int32 = 0
		var totalTime uint64 = 0
		var maxMemory uint64 = 0
		caseBuilders := make([]*ent.CaseResultCreate, len(caseGroup.Cases))

		caseGroupBuilder := repo.data.db.CaseGroupResult.Create().
			SetStatus(util.StatusToInt16(util.Pending)).
			SetTotalTimeCostMs(0).
			SetMaxMemoryCostKB(0).
			SetScore(0)

		wg := sync.WaitGroup{}
		wg.Add(len(caseGroup.Cases))
		for idx2, testCase := range caseGroup.Cases {
			go func(idx2 int, testCase *filemanage.TestCase) {
				caseBuilder := repo.data.db.CaseResult.Create().
					SetStatus(util.StatusToInt16(util.Pending)).
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
					// goroutine 不能直接返回 error
					caseBuilder.SetStatus(util.StatusToInt16(util.SystemError))
					atomic.StoreInt32(&atomicCaseGroupStatus, int32(util.StatusToInt16(util.SystemError)))
					return
				}
				// handle CRLF for both files
				inputBytes = util.RemoveCR(inputBytes)
				ansBytes = util.RemoveCR(ansBytes)

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
					caseBuilder.SetStatus(util.StatusToInt16(util.RuntimeError))
					if result != nil {
						caseBuilder.SetStderr(util.ParseSignalToString(result.ExitStatus))
					} else {
						caseBuilder.SetStderr("Runtime error without stderr message")
					}
					// the whole group is runtime error
					atomic.StoreInt32(&atomicCaseGroupStatus, int32(util.StatusToInt16(util.RuntimeError)))
					return
				}
				repo.log.Infof("submission: %s case: %s result: %+v", msg_submission.UUID, testCase.InputFileLocation, result)

				// check result status
				resStatus := util.ParseGojudgeStatus(result.Status)
				out := result.Files["stdout"]
				caseBuilder.SetStdout(string(out))
				if resStatus == util.StatusToInt16(util.Accepted) && out != nil {
					// compare output
					// TODO may need to support more comparison methods
					if !util.BytesCompareIgnoreSpacesAndNewlines(out, ansBytes) {
						resStatus = util.StatusToInt16(util.WrongAnswer)
					}
				}

				// set caseBuilder metadata
				resTime := result.Time / 1000 / 1000 // ns -> ms
				resMemory := result.Memory / 1024    // Bytes -> KB
				caseBuilder.SetStatus(resStatus)
				caseBuilder.SetTimeCostMs(resTime)
				caseBuilder.SetMemoryCostKB(resMemory)

				// update case group info
				atomic.AddUint64(&totalTime, resTime) // ms
				// maxMemory 需要使用循环原子操作
				// KB
				for {
					oldMax := atomic.LoadUint64(&maxMemory)
					if resMemory <= oldMax {
						break
					}
					if atomic.CompareAndSwapUint64(&maxMemory, oldMax, resMemory) {
						break
					}
				}

				// update submission info
				atomic.AddUint64(&msg_submission.TimeCost, resTime) // ms
				// msg_submission.MemoryCost 需要使用循环原子操作
				for {
					oldMax := atomic.LoadUint64(&msg_submission.MemoryCost)
					if resMemory <= oldMax {
						break
					}
					if atomic.CompareAndSwapUint64(&msg_submission.MemoryCost, oldMax, resMemory) {
						break
					}
				}

				// set case score
				if resStatus != util.StatusToInt16(util.Accepted) {
					caseBuilder.SetScore(0)
					// 原子地更新 caseGroupStatus
					// 只有在当前状态是 AC 时才更新为第一个遇到的非 AC 状态
					currentStatus := atomic.LoadInt32(&atomicCaseGroupStatus)
					if currentStatus == int32(util.StatusToInt16(util.Accepted)) {
						atomic.CompareAndSwapInt32(&atomicCaseGroupStatus, currentStatus, int32(resStatus))
					}
				} else {
					caseBuilder.SetScore(testCase.SubScore)
					atomic.AddInt32(&caseGroupScore, int32(testCase.SubScore))
				}
			}(idx2, &testCase)
		}
		wg.Wait()

		// 从原子变量读取最终状态
		caseGroupStatus = int16(atomic.LoadInt32(&atomicCaseGroupStatus))

		// calculate case group score
		caseGroupBuilder.SetScore(int(caseGroupScore))
		caseGroupBuilder.SetStatus(caseGroupStatus)
		caseGroupBuilder.SetTotalTimeCostMs(totalTime) // ms
		caseGroupBuilder.SetMaxMemoryCostKB(maxMemory) // KB
		caseGroupBuilders[idx] = caseGroupBuilder

		// fill cases array
		casesArrays[idx] = caseBuilders

		// update submission info
		msg_submission.Score += int(caseGroupScore)
		if msg_submission.Status == util.StatusToInt16(util.Accepted) &&
			caseGroupStatus != util.StatusToInt16(util.Accepted) {
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
		err = repo.data.redis.Set(ctx, fmt.Sprintf(cache.SelfTestCacheKeyFmt, msg_self_test.UserID, msg_self_test.UUID), msg_marshal, 1*time.Hour).Err()
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
		repo.log.Errorf("failed to compiled self-test code: %v", err)
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
	msg_self_test.TimeCost = result.Time / 1000 / 1000 // ns -> ms
	msg_self_test.MemoryCost = result.Memory / 1024    // Bytes -> KB
	return nil
}

func NewJudgeRepo(data *Data, logger log.Logger) biz.JudgeRepo {
	return &judgeRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
