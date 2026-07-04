// TriggerAnalysis
export interface TriggerAnalysisRequest {
  timeRangeMinutes: number;
  services: string[];
}

export interface TriggerAnalysisReply {
  reportId: number;
  title: string;
  summary: string;
  noiseCount: number;
  realRiskCount: number;
  fullReportMarkdown: string;
  createdAt: string;
}

// GetReport
export interface GetReportRequest {
  reportId: number;
}

export interface GetReportReply {
  id: number;
  title: string;
  summary: string;
  rootCausesJson: string;
  suggestedActionsJson: string;
  noiseCount: number;
  realRiskCount: number;
  status: string;
  fullReportMarkdown: string;
  createdAt: string;
  acknowledgedAt: string;
}

// ListReports
export interface ListReportsRequest {
  page: number;
  pageSize: number;
  statusFilter: string;
}

export interface ReportSummary {
  id: number;
  title: string;
  noiseCount: number;
  realRiskCount: number;
  status: string;
  createdAt: string;
}

export interface ListReportsReply {
  reports: ReportSummary[];
  total: number;
}

// AckReport
export interface AckReportRequest {
  reportId: number;
}

export interface AckReportReply {
}

// GetMetricsSummary
export interface GetMetricsSummaryRequest {
  services: string[];
}

export interface GetMetricsSummaryReply {
  prometheusStatus: string;
  metricsSnapshotJson: string;
  queriedAt: string;
}
