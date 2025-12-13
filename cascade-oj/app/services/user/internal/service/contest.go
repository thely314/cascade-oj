package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (userService *UserService) GetContests(ctx context.Context, req *pb.GetContestsRequest) (*pb.GetContestsReply, error) {
	contests, err := userService.contestUseCase.GetContests(ctx)
	if err != nil {
		return nil, err
	}
	contestsMetaData := make([]*pb.ContestMetadata, 0, len(contests))
	for i := 0; i < len(contests); i++ {
		contestsMetaData = append(contestsMetaData,
			&pb.ContestMetadata{
				Id:        contests[i].ID,
				Title:     contests[i].Title,
				StartTime: timestamppb.New(contests[i].StartTime),
				EndTime:   timestamppb.New(contests[i].EndTime),
				Status:    contests[i].Status,
			},
		)
	}
	return &pb.GetContestsReply{
		Contests: contestsMetaData,
	}, nil
}

func (userService *UserService) GetSingleContest(ctx context.Context, req *pb.GetSingleContestRequest) (*pb.GetSingleContestReply, error) {
	detailedContest, err := userService.contestUseCase.GetSingleContest(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleContestReply{
		Metadata: &pb.ContestMetadata{
			Id:        detailedContest.Contest.ID,
			Title:     detailedContest.Contest.Title,
			StartTime: timestamppb.New(detailedContest.Contest.StartTime),
			EndTime:   timestamppb.New(detailedContest.Contest.EndTime),
			Status:    detailedContest.Contest.Status,
		}, Description: detailedContest.Description,
	}, nil
}

func (userService *UserService) JoinContest(ctx context.Context, req *pb.JoinContestRequest) (*pb.JoinContestReply, error) {
	result, err := userService.contestUseCase.JoinContest(ctx, req.ContestId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.JoinContestReply{
		IsJoin: result,
	}, nil
}

func (userService *UserService) QuitContest(ctx context.Context, req *pb.QuitContestRequest) (*pb.QuitContestReply, error) {
	result, err := userService.contestUseCase.QuitContest(ctx, req.ContestId, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.QuitContestReply{
		IsJoin: result,
	}, nil
}
