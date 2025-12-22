package data

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
	"cascade-oj/app/gateway/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type contestRepo struct {
	data *Data
	log  *log.Helper
}

func (repo *contestRepo) GetContests(ctx context.Context) ([]*biz.ContestMetadata, error) {
	contests, err := repo.data.grpcUserClient.GetContests(ctx, &pb.GetContestsRequest{})
	if err != nil {
		return nil, err
	}
	var res []*biz.ContestMetadata
	for _, v := range contests.Contests {
		res = append(res, &biz.ContestMetadata{
			ID:        v.Id,
			Title:     v.Title,
			StartTime: v.StartTime.AsTime(),
			EndTime:   v.EndTime.AsTime(),
			Status:    v.Status,
		})
	}
	return res, nil
}

func (repo *contestRepo) GetSingleContest(ctx context.Context, contestID int64) (*biz.ContestDetail, error) {
	contest, err := repo.data.grpcUserClient.GetSingleContest(ctx, &pb.GetSingleContestRequest{ContestId: contestID})
	if err != nil {
		return nil, err
	}
	return &biz.ContestDetail{
		Metadata: biz.ContestMetadata{
			ID:        contest.Metadata.Id,
			Title:     contest.Metadata.Title,
			StartTime: contest.Metadata.StartTime.AsTime(),
			EndTime:   contest.Metadata.EndTime.AsTime(),
			Status:    contest.Metadata.Status,
		},
		Description: contest.Description,
	}, nil
}

func (repo *contestRepo) JoinContest(ctx context.Context, contestID int64, userID int64) error {
	_, err := repo.data.grpcUserClient.JoinContest(ctx, &pb.JoinContestRequest{ContestId: contestID, UserId: userID})
	if err != nil {
		return err
	}
	return nil
}

func (repo *contestRepo) QuitContest(ctx context.Context, contestID int64, userID int64) error {
	_, err := repo.data.grpcUserClient.QuitContest(ctx, &pb.QuitContestRequest{ContestId: contestID, UserId: userID})
	if err != nil {
		return err
	}
	return nil
}

func (repo *contestRepo) GetProblems(ctx context.Context, contestID int64) ([]*biz.ProblemMetadata, error) {
	problems, err := repo.data.grpcUserClient.GetProblems(ctx, &pb.GetProblemsRequest{ContestId: contestID})
	if err != nil {
		return nil, err
	}
	var res []*biz.ProblemMetadata
	for _, v := range problems.Problems {
		res = append(res, &biz.ProblemMetadata{
			ID:            v.Id,
			Title:         v.Title,
			TimeLimitMs:   v.TimeLimitMs,
			MemoryLimitMb: v.MemoryLimitMb,
		})
	}
	return res, nil
}

func (repo *contestRepo) GetSingleProblem(ctx context.Context, contestID int64, problemID int64) (*biz.ProblemDetail, error) {
	problem, err := repo.data.grpcUserClient.GetSingleProblem(ctx, &pb.GetSingleProblemRequest{ProblemId: problemID})
	if err != nil {
		return nil, err
	}
	return &biz.ProblemDetail{
		Metadata: biz.ProblemMetadata{
			ID:            problem.Metadata.Id,
			Title:         problem.Metadata.Title,
			TimeLimitMs:   problem.Metadata.TimeLimitMs,
			MemoryLimitMb: problem.Metadata.MemoryLimitMb,
		},
		Creator:     problem.Creator,
		Description: problem.Description,
	}, nil
}

func NewContestRepo(data *Data, logger log.Logger) biz.ContestRepo {
	return &contestRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}
