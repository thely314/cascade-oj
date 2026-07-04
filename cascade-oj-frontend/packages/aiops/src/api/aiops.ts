import service from './request';
import type {
  TriggerAnalysisRequest,
  TriggerAnalysisReply,
  GetReportReply,
  ListReportsRequest,
  ListReportsReply,
  GetMetricsSummaryRequest,
  GetMetricsSummaryReply,
  AckReportRequest,
  AckReportReply,
} from './types';

export async function triggerAnalysis(
  data: TriggerAnalysisRequest,
): Promise<TriggerAnalysisReply> {
  return service.post('/aiops/v1/analysis/trigger', data);
}

export async function getReport(
  reportId: number,
): Promise<GetReportReply> {
  return service.get(`/aiops/v1/reports/${reportId}`);
}

export async function listReports(
  params: ListReportsRequest,
): Promise<ListReportsReply> {
  const query: Record<string, string> = {
    page: String(params.page),
    page_size: String(params.pageSize),
  };
  let q = new URLSearchParams(query);
  if (params.statusFilter) {
    q.set('status_filter', params.statusFilter);
  }
  const queryString = q.toString();
  return service.get(`/aiops/v1/reports?${queryString}`);
}

export async function ackReport(
  reportId: number,
): Promise<AckReportReply> {
  return service.post(`/aiops/v1/reports/${reportId}/ack`, {});
}

export async function getMetricsSummary(
  params: GetMetricsSummaryRequest,
): Promise<GetMetricsSummaryReply> {
  const searchParams = new URLSearchParams();
  for (const svc of params.services) {
    searchParams.append('services', svc);
  }
  const qs = searchParams.toString();
  return service.get(`/aiops/v1/metrics${qs ? '?' + qs : ''}`);
}
