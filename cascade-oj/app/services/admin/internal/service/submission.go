package service

import (
	pb "cascade-oj/api/cascade/admin/v1"
	"cascade-oj/app/services/admin/internal/biz"
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (adminService *AdminService) GetSubmissions(ctx context.Context, request *pb.GetSubmissionsRequest) (*pb.GetSubmissionsReply, error) {
	submissions, err := adminService.submissionUseCase.GetSubmissions(ctx,
		biz.SubmissionsRequestInfo{
			ProblemID: request.ProblemId,
			ContestID: request.ContestId,
			UserID:    request.UserId,
			Page:      request.Page,
			PageSize:  request.PageSize,
		})
	if err != nil {
		return nil, err
	}
	pbSubmissions := make([]*pb.SubmissionMetadata, 0, len(submissions))
	for _, submission := range submissions {
		pbSubmissions = append(pbSubmissions,
			&pb.SubmissionMetadata{
				SubmissionUuid: submission.SubmissionUUID,
				ProblemId:      submission.ProblemID,
				UserId:         submission.UserID,
				Status:         submission.Status,
				SubmitTime:     timestamppb.New(submission.SubmitTime),
				Score:          submission.Score,
			},
		)
	}
	return &pb.GetSubmissionsReply{
		Submissions: pbSubmissions,
	}, nil
}

func (adminService *AdminService) GetSingleSubmission(ctx context.Context, request *pb.GetSingleSubmissionRequest) (*pb.GetSingleSubmissionReply, error) {
	submission, err := adminService.submissionUseCase.GetSingleSubmission(ctx, request.SubmissionUuid)
	if err != nil {
		return nil, err
	}
	pbCaseResults := make([]*pb.CaseMetadata, 0, len(submission.TestCases))
	for _, testCase := range submission.TestCases {
		pbCaseResults = append(pbCaseResults,
			&pb.CaseMetadata{
				Score:      testCase.Score,
				Status:     testCase.Status,
				TimeCost:   testCase.TimeCost,
				MemoryCost: testCase.MemoryCost,
			},
		)
	}
	return &pb.GetSingleSubmissionReply{
		Metadata: &pb.SubmissionMetadata{
			SubmissionUuid: submission.Submission.SubmissionUUID,
			ProblemId:      submission.Submission.ProblemID,
			UserId:         submission.Submission.UserID,
			Status:         submission.Submission.Status,
			SubmitTime:     timestamppb.New(submission.Submission.SubmitTime),
			Score:          submission.Submission.Score,
		},
		Code:       submission.Code,
		Language:   submission.Language,
		TimeCost:   submission.TimeCost,
		MemoryCost: submission.MemoryCost,
		CaseResults: &pb.GetSingleSubmissionReply_CaseResults{
			Cases: pbCaseResults,
		},
	}, nil
}

func (adminService *AdminService) RejudgeSubmission(ctx context.Context, request *pb.RejudgeSubmissionRequest) (*pb.RejudgeSubmissionReply, error) {
	newSubmissionUUID, err := adminService.submissionUseCase.RejudgeSubmission(ctx, request.SubmissionUuid)
	if err != nil {
		return nil, err
	}
	return &pb.RejudgeSubmissionReply{
		NewSubmissionUuid: newSubmissionUUID,
	}, nil
}
