// 通用API响应类型
export interface ApiResponse<T = any> {
  code: number
  message: string
  data: T
  meta?: {
    timestamp?: string
    request_id?: string
    pagination?: {
      page: number
      limit: number
      total: number
    }
  }
}

// 分页请求参数
export interface PaginationParams {
  page?: number
  limit?: number
  total?: number
}

// 排序参数
export interface SortParams {
  sortBy?: string
  sortOrder?: 'ASC' | 'DESC'
}

// 基础查询参数
export interface BaseQueryParams extends PaginationParams, SortParams {
  keyword?: string
}

// OA用户信息类型
export interface OAUser {
  id: number
  username: string
  real_name: string
  email: string
  phone?: string
  avatar?: string
  role_id: number
  role_name: string // 角色名称，如 "系统管理员"
  department: string
  position: string
  status: 'active' | 'inactive' | 'suspended'
  permissions: string[]
  created_at: string
  updated_at: string
  last_login_at?: string
}

// 设备信息类型
export interface DeviceInfo {
  device_id?: string
  device_type?: 'web' | 'ios' | 'android'
  device_name?: string
  user_agent?: string
  app_version?: string
}

// 登录请求参数
export interface LoginCredentials {
  username: string
  password: string
  platform?: 'oa' | 'app' | 'web'
  device_info?: DeviceInfo
}

// 登录响应数据
export interface LoginResponse {
  user: OAUser
  session: {
    access_token: string
    refresh_token: string
    expires_in: number
  }
}

// 会话信息类型
export interface SessionInfo {
  session_id: string
  user_id: number
  user_real_name: string
  platform: 'app' | 'web' | 'oa'
  device_name: string
  ip_address: string
  location: string
  login_time: string
  last_active_at: string
  is_current?: boolean
}

// 会话统计信息
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

// 惠农用户类型
export interface AppUser {
  id: number
  uuid: string
  phone: string
  real_name: string
  user_type: 'farmer' | 'farm_owner' | 'cooperative' | 'enterprise'
  id_card: string
  status: 'active' | 'inactive' | 'suspended'
  is_real_name_verified: boolean
  is_bank_card_verified: boolean
  is_credit_verified: boolean
  credit_score?: number
  province?: string
  city?: string
  county?: string
  address?: string
  created_at: string
  updated_at: string
  last_login_time?: string
}

// 贷款申请类型
export interface LoanApplication {
  id: string
  user_id: number
  user_name: string
  product_name: string
  amount: number
  term_months: number
  purpose: string
  status: 'pending' | 'reviewing' | 'approved' | 'rejected' | 'disbursed'
  risk_level: 'low' | 'medium' | 'high'
  ai_recommendation: 'approve' | 'reject' | 'manual_review'
  ai_score: number
  ai_reason: string
  created_at: string
  updated_at: string
  reviewed_at?: string
  reviewer_id?: number
  reviewer_name?: string
  review_notes?: string
}

// 农机订单类型
export interface MachineOrder {
  id: string
  user_id: number
  user_name: string
  machine_name: string
  machine_type: string
  rental_date: string
  duration_days: number
  total_amount: number
  status: 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled'
  created_at: string
  updated_at: string
}

// 审批操作类型
export interface ApprovalAction {
  action: 'approve' | 'reject' | 'suspend'
  notes?: string
  amount?: number // 批准金额可能与申请金额不同
}

// 系统监控数据类型
export interface SystemMetrics {
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  active_users: number
  pending_applications: number
  daily_transactions: number
  error_rate: number
  response_time: number
}

// 错误响应类型
export interface ErrorResponse {
  code: number
  message: string
  errors?: Record<string, string[]>
  timestamp?: string
  request_id?: string
} 