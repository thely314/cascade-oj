package data

import (
	"cascade-oj/app/services/user/internal/biz"
	"cascade-oj/ent"
	"cascade-oj/ent/problem"
	"cascade-oj/ent/problemset_problem"
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type ProblemRepo struct {
	data *Data
	log  *log.Helper
}

func (problemRepo *ProblemRepo) GetProblems(ctx context.Context, contestID int64) ([]*ent.Problem, error) {
	queryProblemsID, err := problemRepo.data.db.ProblemSet_Problem.Query().Select(problemset_problem.FieldProblemID).Where(problemset_problem.ProblemSetIDEQ(contestID)).All(ctx)
	if err != nil {
		return nil, err
	}

	problems := make([]*ent.Problem, 0, len(queryProblemsID))
	for i := 0; i < len(queryProblemsID); i++ {
		queryProblem, err := problemRepo.data.db.Problem.Query().Select(problem.FieldID, problem.FieldTitle, problem.FieldProblemType, problem.FieldTimeLimit, problem.FieldMemoryLimit).Where(problem.IDEQ(queryProblemsID[i].ProblemID)).Only(ctx)
		if err != nil {
			return nil, err
		}
		problems = append(problems, queryProblem)
	}
	return problems, nil
}

func (problemRepo *ProblemRepo) GetSingleProblem(ctx context.Context, problemID int64) (*ent.Problem, error) {
	queryProblem, err := problemRepo.data.db.Problem.Query().Select(problem.FieldID, problem.FieldTitle, problem.FieldProblemType, problem.FieldTimeLimit, problem.FieldMemoryLimit).Where(problem.IDEQ(problemID)).Only(ctx)
	if err != nil {
		return nil, err
	}
	return queryProblem, nil
}

func NewProblemRepo(data *Data, logger log.Logger) biz.ProblemRepo {
	return &ProblemRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
