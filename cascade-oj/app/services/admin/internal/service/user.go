package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"context"
)

func (adminService *AdminService) GetUsers(ctx context.Context, request *pb.GetUsersRequest) (*pb.GetUsersReply, error) {
	users, err := adminService.userUseCase.GetUsers(ctx,
		biz.UserRequestInfo{
			Page:     request.Page,
			PageSize: request.PageSize,
		},
	)
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*pb.UserInfo, 0, len(users))
	for _, user := range users {
		pbUsers = append(pbUsers,
			&pb.UserInfo{
				UserId:   user.UserID,
				Username: user.Username,
				Email:    user.Email,
			},
		)
	}
	return &pb.GetUsersReply{
		Users: pbUsers,
	}, nil
}

func (adminService *AdminService) UpdateUserInfo(ctx context.Context, request *pb.UpdateUserInfoRequest) (*pb.UpdateUserInfoReply, error) {
	isUpdated, err := adminService.userUseCase.UpdateUserInfo(ctx,
		biz.UserEditInfo{
			UserID:   request.UserId,
			Username: request.Username,
			Email:    request.Email,
		},
	)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserInfoReply{
		IsUpdated: isUpdated,
	}, nil
}

func (adminService *AdminService) UpdateUserPassword(ctx context.Context, request *pb.UpdateUserPasswordRequest) (*pb.UpdateUserPasswordReply, error) {
	isUpdated, err := adminService.userUseCase.UpdateUserPassword(ctx, request.UserId, request.Password)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserPasswordReply{
		IsUpdated: isUpdated,
	}, nil
}

func (adminService *AdminService) DeleteUser(ctx context.Context, request *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	isDeleted, err := adminService.userUseCase.DeleteUser(ctx, request.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserReply{
		IsDeleted: isDeleted,
	}, nil
}
