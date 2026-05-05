package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (adminService *AdminService) GetContests(ctx context.Context, request *pb.GetContestsRequest) (*pb.GetContestsReply, error) {
	contests, err := adminService.contestUseCase.GetContests(ctx)
	if err != nil {
		return nil, err
	}
	pbContests := make([]*pb.ContestMetadata, 0, len(contests))
	for _, contest := range contests {
		pbContests = append(pbContests,
			&pb.ContestMetadata{
				Id:        contest.ID,
				Title:     contest.Title,
				StartTime: timestamppb.New(contest.StartTime),
				EndTime:   timestamppb.New(contest.EndTime),
				Status:    contest.Status,
			},
		)
	}
	return &pb.GetContestsReply{
		Contests: pbContests,
	}, nil
}

func (adminService *AdminService) GetSingleContest(ctx context.Context, request *pb.GetSingleContestRequest) (*pb.GetSingleContestReply, error) {
	contest, err := adminService.contestUseCase.GetSingleContest(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}
	return &pb.GetSingleContestReply{
		Metadata: &pb.ContestMetadata{
			Id:        contest.Contest.ID,
			Title:     contest.Contest.Title,
			StartTime: timestamppb.New(contest.Contest.StartTime),
			EndTime:   timestamppb.New(contest.Contest.EndTime),
			Status:    contest.Contest.Status,
		},
		Description: contest.Description,
	}, nil
}

func (adminService *AdminService) PostContest(ctx context.Context, request *pb.PostContestRequest) (*pb.PostContestReply, error) {
	contestID, err := adminService.contestUseCase.PostContest(
		ctx,
		biz.ContestCreateInfo{
			Title:         request.Title,
			StartTime:     request.StartTime.AsTime(),
			EndTime:       request.EndTime.AsTime(),
			ProblemIDList: request.Problems.ProblemIds,
		},
	)
	if err != nil {
		return nil, err
	}
	return &pb.PostContestReply{
		ContestId: contestID,
	}, nil
}

func (adminService *AdminService) PutContest(ctx context.Context, request *pb.PutContestRequest) (*pb.PutContestReply, error) {
	isUpdated, err := adminService.contestUseCase.PutContest(
		ctx,
		biz.ContestEditInfo{
			ID:          request.ContestId,
			Title:       request.Title,
			Description: request.Description,
			StartTime:   request.StartTime.AsTime(),
			EndTime:     request.EndTime.AsTime(),
		},
	)
	if err != nil {
		return nil, err
	}
	return &pb.PutContestReply{
		IsUpdated: isUpdated,
	}, nil
}

func (adminService *AdminService) DeleteContest(ctx context.Context, request *pb.DeleteContestRequest) (*pb.DeleteContestReply, error) {
	isDeleted, err := adminService.contestUseCase.DeleteContest(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteContestReply{
		IsDeleted: isDeleted,
	}, nil
}

func (adminService *AdminService) GetRanks(ctx context.Context, request *pb.GetRanksRequest) (*pb.GetRanksReply, error) {
	ranks, err := adminService.contestUseCase.GetRanks(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}
	pbRanks := make([]*pb.GetRanksReply_RankItem, 0, len(ranks))
	for _, rank := range ranks {
		pbRanks = append(pbRanks,
			&pb.GetRanksReply_RankItem{
				UserId:   rank.UserID,
				Username: rank.Username,
				Rank:     rank.Rank,
				Score:    rank.Score,
			},
		)
	}
	return &pb.GetRanksReply{
		Ranks: pbRanks,
	}, nil
}

func (adminService *AdminService) GetContestUsers(ctx context.Context, request *pb.GetContestUsersRequest) (*pb.GetContestUsersReply, error) {
	users, err := adminService.contestUseCase.GetContestUsers(ctx, request.ContestId)
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*pb.UserInfo, 0, len(users))
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.UserInfo{
			UserId:   u.UserID,
			Username: u.Username,
			Email:    u.Email,
		})
	}
	return &pb.GetContestUsersReply{Users: pbUsers}, nil
}

func (adminService *AdminService) AddContestUser(ctx context.Context, request *pb.AddContestUserRequest) (*pb.AddContestUserReply, error) {
	isJoined, err := adminService.contestUseCase.AddContestUser(ctx, request.ContestId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.AddContestUserReply{IsJoined: isJoined}, nil
}

func (adminService *AdminService) RemoveContestUser(ctx context.Context, request *pb.RemoveContestUserRequest) (*pb.RemoveContestUserReply, error) {
	isJoined, err := adminService.contestUseCase.RemoveContestUser(ctx, request.ContestId, request.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.RemoveContestUserReply{IsJoined: isJoined}, nil
}
