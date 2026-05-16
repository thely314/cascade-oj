export interface ContestMetadata {
    id: number;
    title: string;
    startTime: string; // ISO string or timestamp
    endTime: string;
    status: string;
}

export enum ProblemStatus {
    PROBLEM_STATUS_UNSPECIFIED = 0,
    PROBLEM_STATUS_DISABLED = 1,
    PROBLEM_STATUS_ENABLED = 2,
    PROBLEM_STATUS_DELETED = 3,
}

export interface ProblemMetadata {
    id: number;
    title: string;
    timeLimitMs: number;
    memoryLimitMb: number;
    Status: ProblemStatus;
}

export interface SubmissionMetadata {
    submissionUuid: string;
    problemId: number;
    userId: number;
    status: string;
    submitTime: string;
    score: number;
}

export interface CaseMetadata {
    score: number;
    status: string;
    timeCost: number;
    memoryCost: number;
}

export interface RankItem {
    userId: number;
    username: string;
    rank: number;
    score: number;
}

export interface Announcement {
    id: number;
    publisherName: string;
    title: string;
    content: string;
}

export interface UserInfo {
    userId: number;
    username: string;
    email: string;
}

export interface LogEntry {
    id: number;
    message: string;
    timestamp: string;
}

// Request/Response Interfaces

// Contests
export interface GetContestsReply {
    contests: ContestMetadata[];
}

export interface GetSingleContestReply {
    metadata: ContestMetadata;
    description: string;
}

export interface PostContestRequest {
    title: string;
    description: string;
    startTime: string;
    endTime: string;
    problems: {
        problemIds: number[];
    };
}

export interface PostContestReply {
    contestId: number;
}

export interface PutContestRequest {
    title: string;
    description: string;
    startTime: string;
    endTime: string;
}

export interface PutContestReply {
    isUpdated: boolean;
}

export interface DeleteContestReply {
    isDeleted: boolean;
}

// Problems
export interface GetProblemsReply {
    problems: ProblemMetadata[];
}

export interface GetSingleProblemReply {
    metadata: ProblemMetadata;
    creator: string;
    description: string;
    codeTemplate: string;
}

export interface PostProblemRequest {
    metadata: ProblemMetadata;
    creator: string;
    description: string;
    codeTemplate: string; 
}

export interface PostProblemReply {
    problemId: number;
}

export interface PutProblemRequest {
    metadata: ProblemMetadata;
    description: string;
    codeTemplate: string; 
}

export interface PutProblemReply {
    isUpdated: boolean;
}

export interface DeleteProblemReply {
    isDeleted: boolean;
}

export interface PublishProblemReply {
    isSuccess: boolean;
}

export interface DisableProblemRequest {
    ProblemId: number;
}

export interface DisableProblemReply {
    IsSuccess: boolean;
}

// Submissions
export interface GetSubmissionsRequest {
    problemId?: number;
    contestId?: number;
    userId?: number;
    page?: number;
    pageSize?: number;
}

export interface GetSubmissionsReply {
    submissions: SubmissionMetadata[];
}

export interface GetSingleSubmissionReply {
    metadata: SubmissionMetadata;
    code: string;
    language: string;
    timeCost: number;
    memoryCost: number;
    caseResults: {
        cases: CaseMetadata[];
    };
}

export interface RejudgeSubmissionReply {
    newSubmissionUuid: string;
}

// Ranks
export interface GetRanksReply {
    ranks: RankItem[];
}

// Announcements
export interface GetAnnouncementsReply {
    announcements: Announcement[];
}

export interface PostAnnouncementRequest {
    publisherName: string;
    title: string;
    content: string;
}

export interface PostAnnouncementReply {
    announcementId: number;
}

export interface PutAnnouncementRequest {
    title: string;
    content: string;
}

export interface PutAnnouncementReply {
    isUpdated: boolean;
}

export interface DeleteAnnouncementReply {
    isDeleted: boolean;
}

// Users
export interface GetUsersRequest {
    page?: number;
    pageSize?: number;
}

export interface GetUsersReply {
    users: UserInfo[];
}

export interface UpdateUserInfoRequest {
    username: string;
    email: string;
}

export interface UpdateUserInfoReply {
    isUpdated: boolean;
}

export interface DeleteUserReply {
    isDeleted: boolean;
}

export interface GetContestUsersReply {
    users: UserInfo[];
}

export interface AddContestUserReply {
    isJoined: boolean;
}

export interface RemoveContestUserReply {
    isJoined: boolean;
}

export interface UpdateUserPasswordRequest {
    password: string;
}

export interface UpdateUserPasswordReply {
    isUpdated: boolean;
}

// Statistics
export interface GetContestStatisticsReply {
    statisticsJson: string;
}

// Logs
export interface ListLogFilesRequest {
    pageSize?: number;
    offset?: number;
}

export interface ListLogFilesReply {
    filenames: string[];
    total: number;
}

export interface QueryLogContentRequest {
    filename: string;
    level?: string;
    timeStart?: string;
    timeEnd?: string;
    pageSize?: number;
    offset?: number;
}

export interface QueryLogContentReply {
    lines: string[];
    total: number;
}

export interface DownloadLogsRequest {
    filenames: string[];
}
