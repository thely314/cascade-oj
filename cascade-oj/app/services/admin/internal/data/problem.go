package data

import (
	"cascade-oj/app/services/admin/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/problemset_includes"
	"cascade-oj/ent/user"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ProblemRepo struct {
	data *Data
	log  *log.Helper
}

func NewProblemRepo(data *Data, logger log.Logger) biz.ProblemRepo {
	return &ProblemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func mapEntStatusToBizStatus(entStatus problem.UseStatus) biz.ProblemStatus {
	switch entStatus {
	case problem.UseStatusDisabled:
		return biz.ProblemStatusDisabled
	case problem.UseStatusAvailable, problem.UseStatusUsing:
		return biz.ProblemStatusEnabled
	case problem.UseStatusDeleted:
		return biz.ProblemStatusDeleted
	default:
		return biz.ProblemStatusUnspecified
	}
}

func (problemRepo *ProblemRepo) GetProblems(ctx context.Context, contestID int64) ([]*biz.Problem, error) {
	//管理端默认
	if contestID == 0 {
		allProblems, err := problemRepo.data.db.Problem.Query().
			Where(problem.UseStatusNEQ(problem.UseStatusDeleted)). // 不查已删除的
			All(ctx)
		if err != nil {
			return nil, err
		}
		res := make([]*biz.Problem, 0, len(allProblems))
		for _, p := range allProblems {
			res = append(res, &biz.Problem{
				ID:            p.ID,
				Title:         p.Title,
				TimeLimitMs:   int32(p.TimeLimitMs),
				MemoryLimitKB: int32(p.MemoryLimitKB),
				Status:        mapEntStatusToBizStatus(p.UseStatus),
			})
		}
		return res, nil
	}
	//比赛视角
	entProblemIncludes, err := problemRepo.data.db.ProblemSet_Includes.
		Query().
		Select().
		WithProblem(func(pq *ent.ProblemQuery) {
			pq.Select(problem.FieldTitle, problem.FieldTimeLimitMs, problem.FieldMemoryLimitKB, problem.FieldUseStatus)
		}).
		Where(
			problemset_includes.ProblemSetIDEQ(contestID),
		).All(ctx)
	if err != nil {
		return nil, err
	}
	problems := make([]*biz.Problem, 0, len(entProblemIncludes))
	for i := 0; i < len(entProblemIncludes); i++ {
		problems = append(problems,
			&biz.Problem{
				ID:            entProblemIncludes[i].ProblemID,
				Title:         entProblemIncludes[i].Edges.Problem.Title,
				TimeLimitMs:   int32(entProblemIncludes[i].Edges.Problem.TimeLimitMs),
				MemoryLimitKB: int32(entProblemIncludes[i].Edges.Problem.MemoryLimitKB),
				Status:        mapEntStatusToBizStatus(entProblemIncludes[i].Edges.Problem.UseStatus),
			})
	}
	return problems, nil
}

func (problemRepo *ProblemRepo) GetSingleProblem(ctx context.Context, problemID int64) (*biz.DetailedProblem, error) {
	entProblem, err := problemRepo.data.db.Problem.
		Query().
		Select(problem.FieldID,
			problem.FieldTitle,
			problem.FieldTimeLimitMs,
			problem.FieldMemoryLimitKB,
			problem.FieldDescription,
			problem.FieldUseStatus,
			problem.FieldCodeTemplate).
		WithCreator(func(uq *ent.UserQuery) {
			uq.Select(
				user.FieldUsername,
			)
		}).
		Where(
			problem.IDEQ(problemID),
		).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return &biz.DetailedProblem{
		Problem: biz.Problem{
			ID:            entProblem.ID,
			Title:         entProblem.Title,
			TimeLimitMs:   int32(entProblem.TimeLimitMs),
			MemoryLimitKB: int32(entProblem.MemoryLimitKB),
			Status:        mapEntStatusToBizStatus(entProblem.UseStatus),
			CodeTemplate:  entProblem.CodeTemplate,
		},
		CreatorUsername: entProblem.Edges.Creator.Username,
		Description:     entProblem.Description,
	}, nil
}

func (problemRepo *ProblemRepo) PostProblem(ctx context.Context, problemCreateInfo biz.ProblemCreateInfo) (int64, error) {
	creator, err := problemRepo.data.db.User.
		Query().
		Select(
			user.FieldID,
		).
		Where(
			user.UsernameEQ(problemCreateInfo.CreatorUsername),
		).
		Only(ctx)
	if err != nil {
		return -1, err
	}
	newProblem, err := problemRepo.data.db.Problem.
		Create().
		SetTitle(problemCreateInfo.Title).
		SetTimeLimitMs(int(problemCreateInfo.TimeLimitMs)).
		SetMemoryLimitKB(int(problemCreateInfo.MemoryLimitKB)).
		SetJudgeConfigID(1). // TODO default go-judge engine config
		SetDescription(problemCreateInfo.Description).
		SetCreatorID(creator.ID).
		SetUseStatus(problem.UseStatusDisabled).
		SetCodeTemplate(problemCreateInfo.CodeTemplate).
		Save(ctx)
	if err != nil {
		return -1, err
	}
	return newProblem.ID, nil
}

func (problemRepo *ProblemRepo) PutProblem(ctx context.Context, problemEditInfo biz.ProblemEditInfo) (bool, error) {
	err := problemRepo.data.db.Problem.
		UpdateOneID(problemEditInfo.ID).
		SetTitle(problemEditInfo.Title).
		SetTimeLimitMs(int(problemEditInfo.TimeLimitMs)).
		SetMemoryLimitKB(int(problemEditInfo.MemoryLimitKB)).
		SetDescription(problemEditInfo.Description).
		SetCodeTemplate(problemEditInfo.CodeTemplate).
		Exec(ctx)
	return err == nil, err
}

func (problemRepo *ProblemRepo) DeleteProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.
		UpdateOneID(problemID).
		SetUseStatus(problem.UseStatusDeleted).
		Exec(ctx)
	return err == nil, err
}

func (problemRepo *ProblemRepo) PublishProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.
		UpdateOneID(problemID).
		SetUseStatus(problem.UseStatusAvailable).
		Exec(ctx)
	return err == nil, err
}

func (problemRepo *ProblemRepo) DisableProblem(ctx context.Context, problemID int64) (bool, error) {
	err := problemRepo.data.db.Problem.
		UpdateOneID(problemID).
		SetUseStatus(problem.UseStatusDisabled).
		Exec(ctx)
	return err == nil, err
}
