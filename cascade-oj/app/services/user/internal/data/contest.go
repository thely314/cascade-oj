package data

import (
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/problemset"
	"cascade-oj/ent/problemset_user"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestRepo struct {
	data *Data
	log  *log.Helper
}

func (contestRepo *ContestRepo) GetContests(ctx context.Context) ([]*ent.ProblemSet, error) {
	queryContests, err := contestRepo.data.db.ProblemSet.Query().Select(problemset.FieldID, problemset.FieldName, problemset.FieldStartTime, problemset.FieldEndTime, problemset.FieldStatus).All(ctx)

	if err != nil {
		return nil, err
	}

	return queryContests, nil
}

func (contestRepo *ContestRepo) GetSingleContest(ctx context.Context, contestID int64) (*ent.ProblemSet, error) {
	queryContest, err := contestRepo.data.db.ProblemSet.Query().Select(problemset.FieldID, problemset.FieldName, problemset.FieldStartTime, problemset.FieldEndTime, problemset.FieldStatus).Where(problemset.IDEQ(contestID)).Only(ctx)

	if err != nil {
		return nil, err
	}

	return queryContest, nil
}

func (contestRepo *ContestRepo) JoinContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	queryContestRecord, err := contestRepo.data.db.ProblemSet_User.Query().Where(problemset_user.And(problemset_user.ProblemSetIDEQ(contestID), problemset_user.UserIDEQ(userID))).Exist(ctx)
	if err != nil {
		return false, nil
	}
	if queryContestRecord {
		contestRepo.log.Errorf("failed to join contest because user has already joined the same contest")
		return false, nil
	}
	_, err = contestRepo.data.db.ProblemSet_User.Create().SetUserID(userID).SetProblemSetID(contestID).Save(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (contestRepo *ContestRepo) QuitContest(ctx context.Context, contestID int64, userID int64) (bool, error) {
	queryContestRecord, err := contestRepo.data.db.ProblemSet_User.Query().Where(problemset_user.And(problemset_user.ProblemSetIDEQ(contestID), problemset_user.UserIDEQ(userID))).Exist(ctx)
	if err != nil {
		return false, nil
	}
	if !queryContestRecord {
		contestRepo.log.Errorf("failed to quit contest because user has not joined the contest yet")
		return false, nil
	}
	_, err = contestRepo.data.db.ProblemSet_User.Delete().Where(problemset_user.And(problemset_user.ProblemSetIDEQ(contestID), problemset_user.UserIDEQ(userID))).Exec(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func NewContestRepo(data *Data, logger log.Logger) biz.ContestRepo {
	return &ContestRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
