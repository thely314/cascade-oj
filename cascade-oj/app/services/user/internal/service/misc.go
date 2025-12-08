package service

import (
	"context"

	pb "cascade-oj/api/cascade/user/v1"
)

func (userService *UserService) GetRanks(ctx context.Context, req *pb.GetRanksRequest) (*pb.GetRanksReply, error) {
	ret, err := userService.miscUsecase.GetRanks(ctx, req.ContestId)
	if err != nil {
		return nil, err
	}
	var replyRanks []*pb.GetRanksReply_RankItem
	for _, rank := range ret {
		replyRanks = append(replyRanks, &pb.GetRanksReply_RankItem{
			UserId:   rank.UserID,
			Username: rank.Username,
			Score:    int32(rank.Score),
			Rank:     rank.Rank,
		})
	}
	return &pb.GetRanksReply{
		Ranks: replyRanks,
	}, nil
}

func (userService *UserService) GetAnnouncements(ctx context.Context, req *pb.GetAnnouncementsRequest) (*pb.GetAnnouncementsReply, error) {
	ret, err := userService.miscUsecase.GetAnnouncements(ctx)
	if err != nil {
		return nil, err
	}
	var replyAnnouncements []*pb.GetAnnouncementsReply_Announcement
	for _, announcement := range ret {
		replyAnnouncement := &pb.GetAnnouncementsReply_Announcement{
			Id:    announcement.ID,
			Title: announcement.Title,
		}
		replyAnnouncements = append(replyAnnouncements, replyAnnouncement)
	}
	return &pb.GetAnnouncementsReply{
		Announcements: replyAnnouncements,
	}, nil
}

func (userService *UserService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoRequest) (*pb.GetUserInfoReply, error) {
	ret, err := userService.miscUsecase.GetUserInfoByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserInfoReply{
		UserId:   ret.UserID,
		Username: ret.Username,
		Email:    ret.Email,
	}, nil
}

func (userService *UserService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoReply, error) {
	err := userService.miscUsecase.UpdateUserInfo(ctx, req.UserId, req.Username, req.Email)
	return &pb.UpdateUserInfoReply{
		IsUpdated: err == nil,
	}, nil
}
