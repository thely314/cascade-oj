package data

import (
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/problemset"
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestRepo struct {
	data *Data
	log  *log.Helper
}

func (contestRepo *ContestRepo) GetContests(ctx context.Context) ([]*biz.Contest, error) {
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

	return contests, nil
}

func (contestRepo *ContestRepo) GetSingleContest(ctx context.Context, contestID int64) (*biz.DetailedContest, error) {
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

	return &biz.DetailedContest{
		Contest: biz.Contest{
			ID:        queryContest.ID,
			Title:     queryContest.Name,
			StartTime: queryContest.StartTime,
			EndTime:   queryContest.EndTime,
			Status:    string(queryContest.Status),
		},
		Description: queryContest.Description,
	}, nil
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
