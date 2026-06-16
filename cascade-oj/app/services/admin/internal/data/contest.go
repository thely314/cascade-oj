package data

import (
	"context"
	"time"

	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/competitor_list"
	"cascade-oj/ent/problemset"
	"cascade-oj/ent/problemset_includes"
	"cascade-oj/ent/problemsetmanager"
	"cascade-oj/ent/user"
	"cascade-oj/pkg/mq"

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

// 属于该管理员管理的比赛才能被查询到
func (contestRepo *ContestRepo) GetContests(ctx context.Context, adminID int64) ([]*biz.Contest, error) {
	now := time.Now()
	entContests, err := contestRepo.data.db.ProblemSet.
		Query().
		Select(problemset.FieldID,
			problemset.FieldName,
			problemset.FieldStartTime,
			problemset.FieldEndTime,
			problemset.FieldStatus,
		).
		Where(problemset.HasProblemSetManagerWith(problemsetmanager.AdminIDEQ(adminID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	contests := make([]*biz.Contest, 0, len(entContests))
	for i := 0; i < len(entContests); i++ {
		status := getContestStatus(entContests[i].StartTime, entContests[i].EndTime, now)
		contests = append(contests,
			&biz.Contest{
				ID:        entContests[i].ID,
				Title:     entContests[i].Name,
				StartTime: entContests[i].StartTime,
				EndTime:   entContests[i].EndTime,
				Status:    status,
			})
	}

	return sortContests(contests, now), nil
}

func getContestStatus(startTime, endTime, now time.Time) string {
	if now.Before(startTime) {
		return "upcoming"
	}
	if now.After(endTime) {
		return "ended"
	}
	return "ongoing"
}

func sortContests(contests []*biz.Contest, now time.Time) []*biz.Contest {
	if len(contests) <= 1 {
		return contests
	}

	var (
		endedContests   []*biz.Contest
		ongoingContests []*biz.Contest
	)

	for _, contest := range contests {
		if contest.EndTime.Before(now) {
			endedContests = append(endedContests, contest)
		} else {
			ongoingContests = append(ongoingContests, contest)
		}
	}

	for i := 0; i < len(endedContests)-1; i++ {
		for j := 0; j < len(endedContests)-1-i; j++ {
			if endedContests[j].EndTime.Before(endedContests[j+1].EndTime) {
				endedContests[j], endedContests[j+1] = endedContests[j+1], endedContests[j]
			} else if endedContests[j].EndTime.Equal(endedContests[j+1].EndTime) {
				if endedContests[j].StartTime.Before(endedContests[j+1].StartTime) {
					endedContests[j], endedContests[j+1] = endedContests[j+1], endedContests[j]
				}
			}
		}
	}

	for i := 0; i < len(ongoingContests)-1; i++ {
		for j := 0; j < len(ongoingContests)-1-i; j++ {
			if ongoingContests[j].EndTime.After(ongoingContests[j+1].EndTime) {
				ongoingContests[j], ongoingContests[j+1] = ongoingContests[j+1], ongoingContests[j]
			} else if ongoingContests[j].EndTime.Equal(ongoingContests[j+1].EndTime) {
				if ongoingContests[j].StartTime.After(ongoingContests[j+1].StartTime) {
					ongoingContests[j], ongoingContests[j+1] = ongoingContests[j+1], ongoingContests[j]
				}
			}
		}
	}

	return append(ongoingContests, endedContests...)
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

	entProblemIncludes, err := contestRepo.data.db.ProblemSet_Includes.
		Query().
		Select(problemset_includes.FieldProblemID).
		Where(
			problemset_includes.ProblemSetIDEQ(contestID),
		).
		Order(ent.Asc(problemset_includes.FieldProblemOrder)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	problemIDs := make([]int64, 0, len(entProblemIncludes))
	for _, pi := range entProblemIncludes {
		problemIDs = append(problemIDs, pi.ProblemID)
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
		ProblemIDs:  problemIDs,
	}, nil
}

func (contestRepo *ContestRepo) PostContest(ctx context.Context, contestCreateInfo biz.ContestCreateInfo, adminID int64) (int64, error) {
	if len(contestCreateInfo.ProblemIDList) == 0 {
		return -1, nil
	}

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
		SetDescription(contestCreateInfo.Description).
		SetStartTime(contestCreateInfo.StartTime).
		SetEndTime(contestCreateInfo.EndTime).
		SetStatus(problemset.StatusUpcoming).
		Save(ctx)
	if err != nil {
		return -1, err
	}

	if _, err = tx.ProblemSetManager.
		Create().
		SetAdminID(adminID).
		SetProblemSetID(contest.ID).
		Save(ctx); err != nil {
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
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	// delete user cache for contest list
	err = contestRepo.data.publishContestInvalidation(ctx, &mq.ContestCacheMsg{
		ContestID: contest.ID,
		Scale:     "list",
	})
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to publish contest cache invalidation message for contest list: %v", err)
	}

	return contest.ID, nil
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
	}

	currentProblems, err := contestRepo.data.db.ProblemSet_Includes.
		Query().
		Select(problemset_includes.FieldProblemID).
		Where(problemset_includes.ProblemSetIDEQ(contestEditInfo.ID)).
		All(ctx)
	if err != nil {
		return false, err
	}

	currentProblemIDSet := make(map[int64]bool)
	for _, p := range currentProblems {
		currentProblemIDSet[p.ProblemID] = true
	}

	newProblemIDSet := make(map[int64]bool)
	for _, pid := range contestEditInfo.ProblemIDList {
		newProblemIDSet[pid] = true
	}

	tx, err := contestRepo.data.db.Tx(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for pid := range currentProblemIDSet {
		if !newProblemIDSet[pid] {
			_, err = tx.ProblemSet_Includes.
				Delete().
				Where(
					problemset_includes.ProblemSetIDEQ(contestEditInfo.ID),
					problemset_includes.ProblemIDEQ(pid),
				).
				Exec(ctx)
			if err != nil {
				return false, err
			}
		}
	}

	for idx, pid := range contestEditInfo.ProblemIDList {
		if !currentProblemIDSet[pid] {
			_, err = tx.ProblemSet_Includes.
				Create().
				SetProblemID(pid).
				SetProblemSetID(contestEditInfo.ID).
				SetProblemOrder(idx + 1).
				Save(ctx)
			if err != nil {
				return false, err
			}
		} else {
			err = tx.ProblemSet_Includes.
				Update().
				Where(
					problemset_includes.ProblemSetIDEQ(contestEditInfo.ID),
					problemset_includes.ProblemIDEQ(pid),
				).
				SetProblemOrder(idx + 1).
				Exec(ctx)
			if err != nil {
				return false, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	// delete user cache for contest
	err = contestRepo.data.publishContestInvalidation(ctx, &mq.ContestCacheMsg{
		ContestID: contestEditInfo.ID,
		Scale:     "single",
	})
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to publish contest cache invalidation message for contest single: %v", err)
	}
	// delete user cache for problem list for this contest, since problem changes may affect the problem list page
	err = contestRepo.data.publishProblemInvalidation(ctx, &mq.ProblemCacheMsg{
		ProblemID: contestEditInfo.ID, // use contestID as problemID in this case
		Scale:     "list",
	})
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to publish problem cache invalidation message for problem list: %v", err)
	}

	return true, nil
}

func (contestRepo *ContestRepo) DeleteContest(ctx context.Context, contestID int64) (bool, error) {
	err := contestRepo.data.db.ProblemSet.DeleteOneID(contestID).Exec(ctx)
	if err != nil {
		return false, err
	}

	// delete user cache for contest list
	err = contestRepo.data.publishContestInvalidation(ctx, &mq.ContestCacheMsg{
		ContestID: contestID,
		Scale:     "list",
	})
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to publish contest cache invalidation message for contest list: %v", err)
	}

	// delete user cache for contest
	err = contestRepo.data.publishContestInvalidation(ctx, &mq.ContestCacheMsg{
		ContestID: contestID,
		Scale:     "single",
	})
	if err != nil {
		contestRepo.log.Errorf("[cache] failed to publish contest cache invalidation message for contest single: %v", err)
	}

	return true, nil
}

func (contestRepo *ContestRepo) UpdateContestStatuses(ctx context.Context) error {
	now := time.Now()

	allContests, err := contestRepo.data.db.ProblemSet.
		Query().
		Select(
			problemset.FieldID,
			problemset.FieldStartTime,
			problemset.FieldEndTime,
			problemset.FieldStatus,
		).
		All(ctx)
	if err != nil {
		return err
	}

	for _, contest := range allContests {
		expectedStatus := getContestStatus(contest.StartTime, contest.EndTime, now)
		if string(contest.Status) != expectedStatus {
			err = contestRepo.data.db.ProblemSet.
				UpdateOneID(contest.ID).
				SetStatus(problemset.Status(expectedStatus)).
				Exec(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (contestRepo *ContestRepo) GetContestCompetitors(ctx context.Context, contestID int64) ([]*biz.UserInfo, error) {
	entCompetitors, err := contestRepo.data.db.Competitor_List.
		Query().
		WithUser(func(uq *ent.UserQuery) {
			uq.Select(
				user.FieldID,
				user.FieldUsername,
				user.FieldEmail,
			)
		}).
		Where(
			competitor_list.ProblemSetIDEQ(contestID),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}

	competitors := make([]*biz.UserInfo, 0, len(entCompetitors))
	for _, ec := range entCompetitors {
		competitors = append(competitors,
			&biz.UserInfo{
				UserID:   ec.Edges.User.ID,
				Username: ec.Edges.User.Username,
				Email:    ec.Edges.User.Email,
			})
	}

	return competitors, nil
}

func (contestRepo *ContestRepo) PutContestCompetitors(ctx context.Context, contestID int64, userIDs []int64) (bool, error) {
	currentCompetitors, err := contestRepo.data.db.Competitor_List.
		Query().
		Select(competitor_list.FieldUserID).
		Where(competitor_list.ProblemSetIDEQ(contestID)).
		All(ctx)
	if err != nil {
		return false, err
	}

	currentUserIDSet := make(map[int64]bool)
	for _, c := range currentCompetitors {
		currentUserIDSet[c.UserID] = true
	}

	newUserIDSet := make(map[int64]bool)
	for _, uid := range userIDs {
		newUserIDSet[uid] = true
	}

	tx, err := contestRepo.data.db.Tx(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for uid := range currentUserIDSet {
		if !newUserIDSet[uid] {
			_, err = tx.Competitor_List.
				Delete().
				Where(
					competitor_list.ProblemSetIDEQ(contestID),
					competitor_list.UserIDEQ(uid),
				).
				Exec(ctx)
			if err != nil {
				return false, err
			}
		}
	}

	for _, uid := range userIDs {
		if !currentUserIDSet[uid] {
			_, err = tx.Competitor_List.
				Create().
				SetUserID(uid).
				SetProblemSetID(contestID).
				SetTotalScore(0).
				Save(ctx)
			if err != nil {
				return false, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return true, nil
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
