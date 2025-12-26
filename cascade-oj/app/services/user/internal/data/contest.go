package data

import (
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/problemset"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestRepo struct {
	data *Data
	log  *log.Helper
}

func (contestRepo *ContestRepo) GetContests(ctx context.Context) ([]*ent.ProblemSet, error) {
	queryContests, err := contestRepo.data.db.ProblemSet.Query().All(ctx)

	if err != nil {
		return nil, err
	}

	return queryContests, nil
}

func (contestRepo *ContestRepo) GetSingleContest(ctx context.Context, contestID int64) (*ent.ProblemSet, error) {
	queryContest, err := contestRepo.data.db.ProblemSet.Query().Where(problemset.IDEQ(contestID)).Only(ctx)

	if err != nil {
		return nil, err
	}

	return queryContest, nil
}

func (contestRepo *ContestRepo) JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
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

func NewContestRepo(data *Data, logger log.Logger) biz.ContestRepo {
	return &ContestRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
