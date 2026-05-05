package data

import (
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/problemset"
	"cascade-oj/ent/user"
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
)

type ContestRepo struct {
	data *Data
	log  *log.Helper
}

func NewContestRepo(data *Data, logger log.Logger) biz.ContestRepo {
	return &ContestRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (contestRepo *ContestRepo) GetContests(ctx context.Context) ([]*biz.Contest, error) {
	entContests, err := contestRepo.data.db.ProblemSet.
		Query().
		Select(problemset.FieldID,
			problemset.FieldName,
			problemset.FieldStartTime,
			problemset.FieldEndTime,
			problemset.FieldStatus,
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	contests := make([]*biz.Contest, 0, len(entContests))
	for i := 0; i < len(entContests); i++ {
		contests = append(contests,
			&biz.Contest{
				ID:        entContests[i].ID,
				Title:     entContests[i].Name,
				StartTime: entContests[i].StartTime,
				EndTime:   entContests[i].EndTime,
				Status:    string(entContests[i].Status),
			})
	}
	return contests, nil
}
func (contestRepo *ContestRepo) GetSingleContest(ctx context.Context, contestID int64) (*biz.DetailedContest, error) {
	entContest, err := contestRepo.data.db.ProblemSet.
		Query().Select(
		problemset.FieldID, problemset.FieldName,
		problemset.FieldStartTime,
		problemset.FieldEndTime,
		problemset.FieldStatus,
		problemset.FieldDescription,
	).
		Where(
			problemset.IDEQ(contestID),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return &biz.DetailedContest{
		Contest: biz.Contest{
			ID:        entContest.ID,
			Title:     entContest.Name,
			StartTime: entContest.StartTime,
			EndTime:   entContest.EndTime,
			Status:    string(entContest.Status),
		},
		Description: entContest.Description,
	}, nil
}
func (contestRepo *ContestRepo) PostContest(ctx context.Context, contestCreateInfo biz.ContestCreateInfo) (int64, error) {
	tx, err := contestRepo.data.db.Tx(ctx)
	if err != nil {
		return -1, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	contest, err := tx.ProblemSet.
		Create().
		SetName(contestCreateInfo.Title).
		SetStartTime(contestCreateInfo.StartTime).
		SetEndTime(contestCreateInfo.EndTime).
		Save(ctx)
	if err != nil {
		return -1, err
	}
	problems := make([]*ent.ProblemSetIncludesCreate, 0, len(contestCreateInfo.ProblemIDList))
	for problemOrder, problemID := range contestCreateInfo.ProblemIDList {
		problems = append(problems,
			tx.ProblemSet_Includes.
				Create().
				SetProblemID(problemID).
				SetProblemSetID(contest.ID).
				SetProblemOrder(problemOrder+1),
		)
	}

	if err = tx.ProblemSet_Includes.CreateBulk(problems...).Exec(ctx); err != nil {
		return -1, nil
	}
	if err = tx.Commit(); err != nil {
		return -1, err
	} else {
		return contest.ID, nil
	}
}
func (contestRepo *ContestRepo) PutContest(ctx context.Context, contestEditInfo biz.ContestEditInfo) (bool, error) {
	err := contestRepo.data.db.ProblemSet.
		UpdateOneID(contestEditInfo.ID).
		SetName(contestEditInfo.Title).
		SetDescription(contestEditInfo.Description).
		SetStartTime(contestEditInfo.StartTime).
		SetEndTime(contestEditInfo.EndTime).
		Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func (contestRepo *ContestRepo) DeleteContest(ctx context.Context, contestID int64) (bool, error) {
	err := contestRepo.data.db.ProblemSet.DeleteOneID(contestID).Exec(ctx)
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func (contestRepo *ContestRepo) GetRanks(ctx context.Context, contestID int64) ([]*biz.ContestRank, error) {
	entCompetitorScores, err := contestRepo.data.db.Competitor_List.
		Query().
		Select(
			competitor_list.FieldUserID,
			competitor_list.FieldTotalScore,
		).
		WithUser(func(uq *ent.UserQuery) {
			uq.Select(
				user.FieldUsername,
			)
		}).
		Where(
			competitor_list.ProblemSetIDEQ(contestID),
		).
		Order(ent.Desc(competitor_list.FieldTotalScore)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	ranks := make([]*biz.ContestRank, 0, len(entCompetitorScores))
	for i := 0; i < len(entCompetitorScores); i++ {
		ranks = append(ranks,
			&biz.ContestRank{
				UserID:   entCompetitorScores[i].UserID,
				Username: entCompetitorScores[i].Edges.User.Username,
				Rank:     (int32(i + 1)),
				Score:    int32(entCompetitorScores[i].TotalScore),
			})
	}
	return ranks, nil
}

func (contestRepo *ContestRepo) GetContestUsers(ctx context.Context, contestID int64) ([]*biz.User, error) {
	entCompetitors, err := contestRepo.data.db.Competitor_List.
		Query().
		Where(competitor_list.ProblemSetIDEQ(contestID)).
		WithUser(func(uq *ent.UserQuery) {
			uq.Select(user.FieldID, user.FieldUsername, user.FieldEmail)
		}).
		All(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*biz.User, 0, len(entCompetitors))
	for _, competitor := range entCompetitors {
		if competitor.Edges.User == nil {
			continue
		}
		users = append(users, &biz.User{
			UserID:   competitor.Edges.User.ID,
			Username: competitor.Edges.User.Username,
			Email:    competitor.Edges.User.Email,
		})
	}
	return users, nil
}

func (contestRepo *ContestRepo) AddContestUser(ctx context.Context, contestID int64, userID int64) (bool, error) {
	contestExists, err := contestRepo.data.db.ProblemSet.
		Query().
		Where(problemset.IDEQ(contestID)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	if !contestExists {
		return false, errors.New("contest not found")
	}

	userExists, err := contestRepo.data.db.User.
		Query().
		Where(user.IDEQ(userID)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	if !userExists {
		return false, errors.New("user not found")
	}

	isJoined, err := contestRepo.data.db.Competitor_List.
		Query().
		Where(competitor_list.And(
			competitor_list.ProblemSetIDEQ(contestID),
			competitor_list.UserIDEQ(userID),
		)).
		Exist(ctx)
	if err != nil {
		return false, err
	}
	if isJoined {
		return true, nil
	}

	_, err = contestRepo.data.db.Competitor_List.
		Create().
		SetProblemSetID(contestID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (contestRepo *ContestRepo) RemoveContestUser(ctx context.Context, contestID int64, userID int64) (bool, error) {
	_, err := contestRepo.data.db.Competitor_List.
		Delete().
		Where(competitor_list.And(
			competitor_list.ProblemSetIDEQ(contestID),
			competitor_list.UserIDEQ(userID),
		)).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	return false, nil
}
