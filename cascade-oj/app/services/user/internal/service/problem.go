package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
)

func (userService *UserService) GetProblems(ctx context.Context, req *pb.GetProblemsRequest) (*pb.GetProblemsReply, error) {
	problems, err := userService.problemUseCase.GetProblems(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	problemsMetaData := make([]*pb.ProblemMetadata, 0, len(problems))
	for i := 0; i < len(problems); i++ {
		problemsMetaData = append(problemsMetaData,
			&pb.ProblemMetadata{
				Id:            problems[i].ID,
				Title:         problems[i].Title,
				TimeLimitMs:   problems[i].TimeLimitMs,
				MemoryLimitMb: problems[i].MemoryLimitMb,
			},
		)
	}
	return &pb.GetProblemsReply{
		Problems: problemsMetaData,
	}, nil
}

func (userService *UserService) GetSingleProblem(ctx context.Context, req *pb.GetSingleProblemRequest) (*pb.GetSingleProblemReply, error) {
	detailedProblem, err := userService.problemUseCase.GetSingleProblem(ctx, req.ProblemId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleProblemReply{
		Metadata: &pb.ProblemMetadata{
			Id:            detailedProblem.Problem.ID,
			Title:         detailedProblem.Problem.Title,
			TimeLimitMs:   detailedProblem.Problem.TimeLimitMs,
			MemoryLimitMb: detailedProblem.Problem.MemoryLimitMb,
		},
		Creator:     detailedProblem.CreatorUsername,
		Description: detailedProblem.Description,
	}, nil
}
