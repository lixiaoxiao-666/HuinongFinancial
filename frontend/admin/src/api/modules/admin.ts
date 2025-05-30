import request from '../request'
import type { 
  ApiResponse, 
  BaseQueryParams,
  SessionInfo,
  SessionStatistics,
  LoanApplication,
  MachineOrder,
  AppUser,
  ApprovalAction,
  SystemMetrics
} from '../types'

/**
 * 管理员会话管理API
 */
export const adminSessionApi = {
  /**
   * 获取系统所有活跃会话 (管理员权限)
   */
  getActiveSessions(params?: {
    user_id?: number
    user_id_type?: 'app_user' | 'oa_user'
    platform?: 'app' | 'web' | 'oa'
    ip_address?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<{
    total: number
    sessions: SessionInfo[]
  }>> {
    return request.get('/admin/sessions/active', { params })
  },

  /**
   * 管理员强制注销指定会话
   */
  revokeSession(session_id: string): Promise<ApiResponse<null>> {
    return request.post('/admin/sessions/revoke', {
      session_id_to_revoke: session_id
    })
  },

  /**
   * 管理员强制注销指定用户的所有会话
   */
  revokeUserSessions(data: {
    user_id_to_revoke: number
    user_id_type: 'app_user' | 'oa_user'
  }): Promise<ApiResponse<{ revoked_count: number }>> {
    return request.post('/admin/sessions/revoke-user', data)
  },

  /**
   * 获取会话统计信息
   */
  getSessionStatistics(): Promise<ApiResponse<SessionStatistics>> {
    return request.get('/admin/sessions/statistics')
  }
}

/**
 * 贷款管理API
 */
export const loanAdminApi = {
  /**
   * 获取贷款申请列表
   */
  getApplications(params?: BaseQueryParams & {
    status?: string
    risk_level?: string
    ai_recommendation?: string
    date_from?: string
    date_to?: string
  }): Promise<ApiResponse<{
    total: number
    applications: LoanApplication[]
  }>> {
    return request.get('/admin/loans/applications', { params })
  },

  /**
   * 获取贷款申请详情
   */
  getApplicationDetail(id: string): Promise<ApiResponse<LoanApplication & {
    user_info: AppUser
    documents: Array<{
      type: string
      url: string
      status: string
    }>
    credit_report: any
    ai_analysis: any
  }>> {
    return request.get(`/admin/loans/applications/${id}`)
  },

  /**
   * 审批贷款申请
   */
  approveApplication(id: string, action: ApprovalAction): Promise<ApiResponse<null>> {
    return request.post(`/admin/loans/applications/${id}/approve`, action)
  },

  /**
   * 批量审批
   */
  batchApprove(data: {
    application_ids: string[]
    action: ApprovalAction
  }): Promise<ApiResponse<{ success_count: number, failed_count: number }>> {
    return request.post('/admin/loans/applications/batch-approve', data)
  },

  /**
   * 获取审批统计
   */
  getApprovalStatistics(): Promise<ApiResponse<{
    total_pending: number
    total_today: number
    approval_rate: number
    average_processing_time: number
    risk_distribution: Record<string, number>
  }>> {
    return request.get('/admin/loans/statistics')
  }
}

/**
 * 农机管理API
 */
export const machineAdminApi = {
  /**
   * 获取农机订单列表
   */
  getOrders(params?: BaseQueryParams & {
    status?: string
    date_from?: string
    date_to?: string
  }): Promise<ApiResponse<{
    total: number
    orders: MachineOrder[]
  }>> {
    return request.get('/admin/machines/orders', { params })
  },

  /**
   * 获取订单详情
   */
  getOrderDetail(id: string): Promise<ApiResponse<MachineOrder & {
    user_info: AppUser
    machine_info: any
    payment_info: any
  }>> {
    return request.get(`/admin/machines/orders/${id}`)
  },

  /**
   * 确认订单
   */
  confirmOrder(id: string, notes?: string): Promise<ApiResponse<null>> {
    return request.post(`/admin/machines/orders/${id}/confirm`, { notes })
  },

  /**
   * 取消订单
   */
  cancelOrder(id: string, reason: string): Promise<ApiResponse<null>> {
    return request.post(`/admin/machines/orders/${id}/cancel`, { reason })
  }
}

/**
 * 用户管理API
 */
export const userAdminApi = {
  /**
   * 获取用户列表
   */
  getUsers(params?: BaseQueryParams & {
    user_type?: string
    status?: string
    verification_status?: string
  }): Promise<ApiResponse<{
    total: number
    users: AppUser[]
  }>> {
    return request.get('/admin/users', { params })
  },

  /**
   * 获取用户详情
   */
  getUserDetail(id: number): Promise<ApiResponse<AppUser & {
    verification_info: any
    loan_history: LoanApplication[]
    machine_orders: MachineOrder[]
    credit_info: any
  }>> {
    return request.get(`/admin/users/${id}`)
  },

  /**
   * 更新用户状态
   */
  updateUserStatus(id: number, status: string, reason?: string): Promise<ApiResponse<null>> {
    return request.post(`/admin/users/${id}/status`, { status, reason })
  },

  /**
   * 重置用户密码
   */
  resetUserPassword(id: number): Promise<ApiResponse<{ temporary_password: string }>> {
    return request.post(`/admin/users/${id}/reset-password`)
  }
}

/**
 * 系统监控API
 */
export const systemAdminApi = {
  /**
   * 获取系统指标
   */
  getSystemMetrics(): Promise<ApiResponse<SystemMetrics>> {
    return request.get('/admin/system/metrics')
  },

  /**
   * 获取操作日志
   */
  getOperationLogs(params?: BaseQueryParams & {
    action_type?: string
    date_from?: string
    date_to?: string
  }): Promise<ApiResponse<{
    total: number
    logs: Array<{
      id: number
      user_id: number
      user_name: string
      action: string
      target: string
      ip_address: string
      created_at: string
    }>
  }>> {
    return request.get('/admin/system/logs', { params })
  }
} 