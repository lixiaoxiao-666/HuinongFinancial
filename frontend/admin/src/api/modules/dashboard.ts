import { request } from '../request'

// 仪表盘概览数据接口
export interface DashboardOverview {
  total_applications: number
  pending_applications: number
  approved_applications: number
  rejected_applications: number
  total_loan_amount: number
  average_loan_amount: number
  default_rate: number
  total_users: number
  active_users: number
  today_new_users: number
  total_machines: number
  rented_machines: number
  maintenance_machines: number
  total_articles: number
  published_articles: number
  draft_articles: number
}

// 风险监控数据接口
export interface RiskMonitoring {
  high_risk_applications: number
  ai_approval_rate: number
  manual_review_rate: number
  avg_processing_time: number
  system_alerts: Array<{
    id: number
    type: 'warning' | 'error' | 'info'
    message: string
    created_at: string
  }>
}

// 会话统计接口
export interface SessionStatistics {
  total_active_sessions: number
  platform_distribution: {
    app: number
    web: number
    oa: number
  }
  daily_peak_users: number
  average_session_duration_minutes: number
}

// 待处理任务接口
export interface PendingTask {
  id: string
  title: string
  description: string
  type: string
  priority: 'high' | 'medium' | 'low'
  created_at: string
  user_info?: {
    user_id: number
    real_name: string
    phone: string
  }
}

// 最新申请接口
export interface RecentApplication {
  id: string
  user_info: {
    user_id: number
    real_name: string
    phone: string
    user_avatar?: string
  }
  product_name: string
  amount: number
  status: string
  status_text: string
  created_at: string
}

/**
 * 获取仪表盘概览数据
 */
export const getDashboardOverview = () => {
  return request<DashboardOverview>({
    url: '/api/oa/admin/dashboard/overview',
    method: 'GET'
  })
}

/**
 * 获取风险监控数据
 */
export const getRiskMonitoring = () => {
  return request<RiskMonitoring>({
    url: '/api/oa/admin/dashboard/risk-monitoring',
    method: 'GET'
  })
}

/**
 * 获取会话统计数据
 */
export const getSessionStatistics = () => {
  return request<SessionStatistics>({
    url: '/api/oa/admin/sessions/statistics',
    method: 'GET'
  })
}

/**
 * 获取待处理任务列表
 */
export const getPendingTasks = (params?: {
  status?: string
  type?: string
  limit?: number
}) => {
  return request<{
    total: number
    tasks: PendingTask[]
  }>({
    url: '/api/oa/admin/tasks/pending',
    method: 'GET',
    params
  })
}

/**
 * 获取最新申请列表
 */
export const getRecentApplications = (params?: {
  status?: string
  limit?: number
}) => {
  return request<{
    total: number
    applications: RecentApplication[]
  }>({
    url: '/api/oa/admin/loans/applications',
    method: 'GET',
    params: {
      ...params,
      sort: 'newest'
    }
  })
}

/**
 * 处理任务
 */
export const handleTask = (taskId: string, action: string) => {
  return request<any>({
    url: `/api/oa/admin/tasks/${taskId}/handle`,
    method: 'POST',
    data: { action }
  })
}

/**
 * 获取系统公告
 */
export const getSystemAnnouncements = (params?: {
  limit?: number
  status?: 'published' | 'draft'
}) => {
  return request<{
    total: number
    announcements: Array<{
      id: number
      title: string
      content: string
      created_at: string
      updated_at: string
    }>
  }>({
    url: '/api/oa/admin/content/announcements',
    method: 'GET',
    params
  })
} 