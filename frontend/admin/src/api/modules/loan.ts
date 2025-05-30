import { request } from '../request'
import type { ApiResponse } from '@/types/api'

/**
 * 贷款管理相关接口类型定义
 */
export interface LoanApplication {
  id: string
  user_info: {
    user_id: number
    real_name: string
    phone: string
    avatar?: string
    is_verified: boolean
    bank_verified: boolean
    id_card: string
    user_type: string
    created_at: string
    credit_score: number
  }
  product_name: string
  amount: number
  term: number
  status: string
  applied_at: string
  ai_assessment?: {
    suggestion: string
    risk_score: number
    confidence: number
    risk_factors: Array<{
      name: string
      score: number
    }>
    analysis_text: string
  }
  materials: Array<{
    id: string
    name: string
    file_url: string
    verified: boolean
  }>
  approval_records: Array<{
    id: string
    action: string
    operator_name: string
    comments?: string
    created_at: string
  }>
  reviewer?: {
    id: number
    name: string
  }
}

export interface LoanProduct {
  id: number
  product_name: string
  description: string
  min_amount: number
  max_amount: number
  interest_rate: number
  max_term: number
}

export interface LoanStatistics {
  total_applications: number
  pending_applications: number
  approved_applications: number
  rejected_applications: number
  total_loan_amount: number
  approval_rate: number
  risk_alerts: number
}

export interface Reviewer {
  id: number
  name: string
  avatar?: string
  department: string
  position: string
  current_cases: number
  expertise: string[]
  efficiency_rating: number
  status: 'online' | 'offline'
}

/**
 * 贷款管理API接口
 */
export const loanApi = {
  /**
   * OA管理员 - 获取贷款申请列表
   */
  getApplicationList(params: {
    status?: string
    user_id?: number
    product_id?: number
    risk_level?: string
    date_range_start?: string
    date_range_end?: string
    applicant?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<{
    total: number
    applications: LoanApplication[]
  }>> {
    return request.get('/api/oa/admin/loans/applications', { params })
  },

  /**
   * OA管理员 - 获取申请详情
   */
  getApplicationDetail(applicationId: string): Promise<ApiResponse<LoanApplication>> {
    return request.get(`/api/oa/admin/loans/applications/${applicationId}`)
  },

  /**
   * OA管理员 - 批准申请
   */
  approveApplication(applicationId: string, data: {
    approved_amount: number
    approved_term: number
    interest_rate: number
    repayment_start_date: string
    comments?: string
    notify_user: boolean
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/oa/admin/loans/applications/${applicationId}/approve`, data)
  },

  /**
   * OA管理员 - 拒绝申请
   */
  rejectApplication(applicationId: string, data: {
    rejection_reason_code: string
    comments: string
    notify_user: boolean
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/oa/admin/loans/applications/${applicationId}/reject`, data)
  },

  /**
   * OA管理员 - 退回申请
   */
  returnApplication(applicationId: string, data: {
    return_reason: string
    required_materials: string[]
    comments?: string
    notify_user: boolean
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/oa/admin/loans/applications/${applicationId}/return`, data)
  },

  /**
   * OA管理员 - 开始人工审核
   */
  startReview(applicationId: string, data: {
    reviewer_id: number
    review_department: string
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/oa/admin/loans/applications/${applicationId}/start-review`, data)
  },

  /**
   * OA管理员 - 重试AI评估
   */
  retryAIAssessment(applicationId: string): Promise<ApiResponse<any>> {
    return request.post(`/api/oa/admin/loans/applications/${applicationId}/retry-ai`)
  },

  /**
   * OA管理员 - 获取统计数据
   */
  getStatistics(): Promise<ApiResponse<LoanStatistics>> {
    return request.get('/api/oa/admin/loans/statistics')
  },

  /**
   * 获取贷款产品列表
   */
  getLoanProducts(): Promise<ApiResponse<LoanProduct[]>> {
    return request.get('/api/user/loan/products')
  },

  /**
   * 获取审核员列表
   */
  getReviewers(): Promise<ApiResponse<Reviewer[]>> {
    return request.get('/api/oa/admin/loans/reviewers')
  },

  /**
   * 分配审核员
   */
  assignReviewer(data: {
    application_ids: string[]
    reviewer_id: number
    assign_reason: string
    priority: string
    deadline?: string
    notes?: string
    send_email: boolean
    send_sms: boolean
    create_reminder: boolean
  }): Promise<ApiResponse<any>> {
    return request.post('/api/oa/admin/loans/assign-reviewer', data)
  },

  /**
   * 获取审核员推荐
   */
  getReviewerRecommendations(data: {
    application_ids: string[]
  }): Promise<ApiResponse<Array<{
    id: number
    name: string
    avatar?: string
    recommendation_reason: string
    match_score: number
  }>>> {
    return request.post('/api/oa/admin/loans/reviewer-recommendations', data)
  },

  /**
   * 批量操作
   */
  batchApprove(applicationIds: string[]): Promise<ApiResponse<any>> {
    return request.post('/api/oa/admin/loans/batch-approve', {
      application_ids: applicationIds
    })
  },

  batchReject(applicationIds: string[]): Promise<ApiResponse<any>> {
    return request.post('/api/oa/admin/loans/batch-reject', {
      application_ids: applicationIds
    })
  },

  /**
   * 导出数据
   */
  exportApplications(params: any): Promise<ApiResponse<any>> {
    return request.get('/api/oa/admin/loans/export', { 
      params,
      responseType: 'blob'
    })
  },

  // ===== 用户端接口 =====

  /**
   * 用户端 - 获取贷款产品列表
   */
  getUserLoanProducts(params?: {
    user_type?: string
    amount_min?: number
    amount_max?: number
  }): Promise<ApiResponse<LoanProduct[]>> {
    return request.get('/api/user/loan/products', { params })
  },

  /**
   * 用户端 - 获取产品详情
   */
  getLoanProductDetail(productId: number): Promise<ApiResponse<LoanProduct>> {
    return request.get(`/api/user/loan/products/${productId}`)
  },

  /**
   * 用户端 - 提交贷款申请
   */
  submitApplication(data: {
    product_id: number
    amount: number
    term: number
    purpose: string
    materials: Array<{
      type: string
      file_url: string
    }>
  }): Promise<ApiResponse<{ application_id: string }>> {
    return request.post('/api/user/loan/applications', data)
  },

  /**
   * 用户端 - 获取我的申请列表
   */
  getMyApplications(params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<{
    total: number
    applications: LoanApplication[]
  }>> {
    return request.get('/api/user/loan/applications', { params })
  },

  /**
   * 用户端 - 获取申请详情
   */
  getMyApplicationDetail(applicationId: string): Promise<ApiResponse<LoanApplication>> {
    return request.get(`/api/user/loan/applications/${applicationId}`)
  },

  /**
   * 用户端 - 取消申请
   */
  cancelApplication(applicationId: string): Promise<ApiResponse<any>> {
    return request.delete(`/api/user/loan/applications/${applicationId}`)
  },

  // ===== 贷款合同管理 =====

  /**
   * 生成贷款合同
   */
  generateContract(applicationId: string): Promise<ApiResponse<{
    contract_id: string
    contract_url: string
    expires_at: string
    signing_required: boolean
  }>> {
    return request.post(`/api/loans/contracts/${applicationId}`)
  },

  /**
   * 签署合同
   */
  signContract(contractId: string, data: {
    signature: string
    signing_location: string
    ip_address: string
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/loans/contracts/${contractId}/sign`, data)
  },

  /**
   * 获取合同详情
   */
  getContractDetail(contractId: string): Promise<ApiResponse<any>> {
    return request.get(`/api/loans/contracts/${contractId}`)
  },

  // ===== 额度管理 =====

  /**
   * 获取信用额度
   */
  getCreditLimit(): Promise<ApiResponse<{
    total_limit: number
    available_limit: number
    used_limit: number
    credit_score: number
    credit_level: string
    limit_details: Array<{
      product_code: string
      product_name: string
      limit: number
      available: number
      interest_rate: number
    }>
    factors: {
      income_stability: number
      credit_history: number
      business_scale: number
      collateral_value: number
    }
    suggestions: string[]
  }>> {
    return request.get('/api/loans/credit-limit')
  },

  /**
   * 申请额度提升
   */
  requestLimitIncrease(data: {
    requested_limit: number
    reason: string
    additional_materials: Array<{
      type: string
      file_url: string
    }>
  }): Promise<ApiResponse<any>> {
    return request.post('/api/loans/credit-limit/increase', data)
  },

  // ===== 还款管理 =====

  /**
   * 获取还款计划
   */
  getRepaymentSchedule(loanId: string): Promise<ApiResponse<{
    loan_id: string
    total_amount: number
    remaining_amount: number
    total_installments: number
    remaining_installments: number
    next_payment: {
      installment_no: number
      due_date: string
      principal: number
      interest: number
      total_amount: number
      status: string
    }
    schedule: Array<{
      installment_no: number
      due_date: string
      principal: number
      interest: number
      total_amount: number
      paid_amount?: number
      paid_date?: string
      status: string
    }>
  }>> {
    return request.get(`/api/loans/${loanId}/repayment-schedule`)
  },

  /**
   * 主动还款
   */
  makeRepayment(loanId: string, data: {
    payment_method: string
    amount: number
    installment_nos: number[]
    bank_card_id?: string
  }): Promise<ApiResponse<{
    payment_id: string
    paid_amount: number
    remaining_principal: number
    next_due_date: string
    transaction_time: string
  }>> {
    return request.post(`/api/loans/${loanId}/repayment`, data)
  },

  /**
   * 提前还款计算
   */
  calculatePrepayment(loanId: string, data: {
    prepayment_type: 'full' | 'partial'
    amount?: number
  }): Promise<ApiResponse<{
    prepayment_amount: number
    interest_saved: number
    remaining_principal: number
    new_monthly_payment?: number
    new_final_payment_date?: string
    prepayment_fee: number
    total_payment_required: number
  }>> {
    return request.post(`/api/loans/${loanId}/prepayment-calculate`, data)
  },

  /**
   * 获取还款记录
   */
  getPaymentHistory(loanId: string, params?: {
    page?: number
    limit?: number
  }): Promise<ApiResponse<{
    total: number
    payments: Array<{
      payment_id: string
      payment_date: string
      amount: number
      principal: number
      interest: number
      penalty?: number
      payment_method: string
      status: string
    }>
  }>> {
    return request.get(`/api/loans/${loanId}/payments`, { params })
  },

  // ===== 放款管理 =====

  /**
   * 获取放款信息
   */
  getDisbursementInfo(loanId: string): Promise<ApiResponse<{
    loan_id: string
    approved_amount: number
    disbursement_amount: number
    service_fee: number
    bank_account: {
      account_number: string
      bank_name: string
      account_holder: string
    }
    estimated_arrival_time: string
    status: string
  }>> {
    return request.get(`/api/loans/${loanId}/disbursement-info`)
  },

  /**
   * 申请放款
   */
  requestDisbursement(loanId: string, data: {
    confirm_amount: number
    bank_account_id: string
  }): Promise<ApiResponse<any>> {
    return request.post(`/api/loans/${loanId}/disburse`, data)
  },

  // ===== 统计查询 =====

  /**
   * 获取贷款概览
   */
  getLoanOverview(): Promise<ApiResponse<{
    active_loans: number
    total_borrowed: number
    total_repaid: number
    outstanding_balance: number
    credit_limit_used: number
    next_payment: {
      amount: number
      due_date: string
      days_until_due: number
    }
    payment_status: string
  }>> {
    return request.get('/api/loans/overview')
  },

  /**
   * 获取历史贷款
   */
  getLoanHistory(params?: {
    status?: string
    page?: number
    limit?: number
  }): Promise<ApiResponse<{
    total: number
    loans: Array<{
      loan_id: string
      product_name: string
      amount: number
      status: string
      start_date: string
      end_date?: string
      remaining_amount: number
    }>
  }>> {
    return request.get('/api/loans/history', { params })
  },

  // ===== 提醒服务 =====

  /**
   * 获取还款提醒设置
   */
  getReminderSettings(): Promise<ApiResponse<{
    sms_reminder: boolean
    email_reminder: boolean
    push_notification: boolean
    reminder_days: number[]
    reminder_time: string
  }>> {
    return request.get('/api/loans/reminder-settings')
  },

  /**
   * 更新提醒设置
   */
  updateReminderSettings(data: {
    sms_reminder: boolean
    email_reminder?: boolean
    push_notification: boolean
    reminder_days: number[]
    reminder_time: string
  }): Promise<ApiResponse<any>> {
    return request.put('/api/loans/reminder-settings', data)
  }
} 