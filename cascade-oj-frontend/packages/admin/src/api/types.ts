export interface ContestMetadata {
    id: number;
    title: string;
    start_time: string; // ISO string or timestamp
    end_time: string;
    status: string;
}

export interface ProblemMetadata {
    id: number;
    title: string;
    time_limit_ms: number;
    memory_limit_mb: number;
}

export interface SubmissionMetadata {
    submission_uuid: string;
    problem_id: number;
    user_id: number;
    status: string;
    submit_time: string;
    score: number;
}

export interface CaseMetadata {
    score: number;
    status: string;
    time_cost: number;
    memory_cost: number;
}

export interface RankItem {
    user_id: number;
    username: string;
    rank: number;
    score: number;
}

export interface Announcement {
    id: number;
    publisher_name: string;
    title: string;
    content: string;
}

export interface UserInfo {
    user_id: number;
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
    start_time: string;
    end_time: string;
    problems: {
        problem_ids: number[];
    };
}

export interface PostContestReply {
    contest_id: number;
}

export interface PutContestRequest {
    title: string;
    description: string;
    start_time: string;
    end_time: string;
}

export interface PutContestReply {
    is_updated: boolean;
}

export interface DeleteContestReply {
    is_deleted: boolean;
}

// Problems
export interface GetProblemsReply {
    problems: ProblemMetadata[];
}

export interface GetSingleProblemReply {
    metadata: ProblemMetadata;
    creator: string;
    description: string;
}

export interface PostProblemRequest {
    metadata: ProblemMetadata;
    creator: string;
    description: string;
}

export interface PostProblemReply {
    problem_id: number;
}

export interface PutProblemRequest {
    metadata: ProblemMetadata;
    description: string;
}

export interface PutProblemReply {
    is_updated: boolean;
}

export interface DeleteProblemReply {
    is_deleted: boolean;
}

export interface PublishProblemReply {
    is_success: boolean;
}

// Submissions
export interface GetSubmissionsRequest {
    problem_id?: number;
    contest_id?: number;
    user_id?: number;
    page?: number;
    page_size?: number;
}

export interface GetSubmissionsReply {
    submissions: SubmissionMetadata[];
}

export interface GetSingleSubmissionReply {
    metadata: SubmissionMetadata;
    code: string;
    language: string;
    time_cost: number;
    memory_cost: number;
    case_results: {
        cases: CaseMetadata[];
    };
}

export interface RejudgeSubmissionReply {
    new_submission_uuid: string;
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
    publisher_name: string;
    title: string;
    content: string;
}

export interface PostAnnouncementReply {
    announcement_id: number;
}

export interface PutAnnouncementRequest {
    title: string;
    content: string;
}

export interface PutAnnouncementReply {
    is_updated: boolean;
}

export interface DeleteAnnouncementReply {
    is_deleted: boolean;
}

// Users
export interface GetUsersRequest {
    page?: number;
    page_size?: number;
}

export interface GetUsersReply {
    users: UserInfo[];
}

export interface UpdateUserInfoRequest {
    username: string;
    email: string;
}

export interface UpdateUserInfoReply {
    is_updated: boolean;
}

export interface DeleteUserReply {
    is_deleted: boolean;
}

// Statistics
export interface GetContestStatisticsReply {
    statistics_json: string;
}

// Logs
export interface GetLogsRequest {
    page?: number;
    page_size?: number;
}

export interface GetLogsReply {
    logs: LogEntry[];
}
