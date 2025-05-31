import type { 
  LoginResponse, 
  DashboardData, 
  PaginationResponse, 
  LoanApplication, 
  ApplicationDetail, 
  AdminUser, 
  OperationLog, 
  SystemConfig 
} from '@/types'

// 模拟延迟
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

// 模拟登录数据
export const mockLogin = async (data: { username: string; password: string }): Promise<LoginResponse> => {
  await delay(1000) // 模拟网络延迟
  
  if ((data.username === 'admin' && data.password === 'admin123') || 
      (data.username === 'reviewer' && data.password === 'reviewer123')) {
    return {
      token: 'mock_token_' + Date.now(),
      admin_user_id: 'admin_001',
      username: data.username,
      role: data.username === 'admin' ? 'ADMIN' : '审批员'
    }
  }
  
  throw new Error('用户名或密码错误')
}

// 模拟工作台数据
export const mockDashboard = async (): Promise<DashboardData> => {
  await delay(800)
  
  return {
    pending_count: 8,
    approved_count: 25,
    rejected_count: 3,
    ai_processing_count: 5,
    ai_enabled: true,
    ai_processing_rate: 85,
    pending_tasks: [
      {
        task_id: 'T001',
        task_type: '贷款审批',
        title: '张三的贷款申请需要审核',
        priority: 'high',
        created_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(), // 2小时前
        application_id: 'APP001'
      },
      {
        task_id: 'T002',
        task_type: '贷款审批',
        title: '李四的补充材料已提交',
        priority: 'medium',
        created_at: new Date(Date.now() - 4 * 60 * 60 * 1000).toISOString(), // 4小时前
        application_id: 'APP002'
      },
      {
        task_id: 'T003',
        task_type: '贷款审批',
        title: '王五的贷款申请等待处理',
        priority: 'low',
        created_at: new Date(Date.now() - 6 * 60 * 60 * 1000).toISOString(), // 6小时前
        application_id: 'APP003'
      }
    ],
    recent_activities: [
      {
        activity_id: 'A001',
        activity_type: '贷款审批',
        description: '批准了张三的贷款申请',
        timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString(), // 30分钟前
        operator: '管理员'
      },
      {
        activity_id: 'A002',
        activity_type: '用户管理',
        description: '创建了新用户 李审批员',
        timestamp: new Date(Date.now() - 60 * 60 * 1000).toISOString(), // 1小时前
        operator: '管理员'
      },
      {
        activity_id: 'A003',
        activity_type: 'AI设置',
        description: '调整了AI风险阈值为70分',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString(), // 2小时前
        operator: '管理员'
      }
    ]
  }
}

// 模拟贷款申请列表
export const mockApplications = async (params: any): Promise<PaginationResponse<LoanApplication>> => {
  await delay(600)
  
  const allApplications: LoanApplication[] = [
    {
      application_id: 'APP001',
      applicant_name: '张三',
      amount: 50000,
      status: '待人工复核',
      ai_risk_score: 25,
      ai_suggestion: '建议批准，申请人信用良好，收入稳定',
      submission_time: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      application_id: 'APP002',
      applicant_name: '李四',
      amount: 80000,
      status: 'AI_审批中',
      ai_risk_score: 45,
      ai_suggestion: '需要进一步核实收入证明',
      submission_time: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      application_id: 'APP003',
      applicant_name: '王五',
      amount: 30000,
      status: '已批准',
      ai_risk_score: 15,
      ai_suggestion: '低风险客户，建议快速批准',
      submission_time: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      application_id: 'APP004',
      applicant_name: '赵六',
      amount: 120000,
      status: '已拒绝',
      ai_risk_score: 85,
      ai_suggestion: '高风险申请，建议拒绝',
      submission_time: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
    },
    {
      application_id: 'APP005',
      applicant_name: '孙七',
      amount: 65000,
      status: '待人工复核',
      ai_risk_score: 55,
      ai_suggestion: '中等风险，建议人工审核',
      submission_time: new Date(Date.now() - 12 * 60 * 60 * 1000).toISOString()
    }
  ]
  
  // 简单的筛选逻辑
  let filteredApplications = allApplications
  if (params.status_filter) {
    filteredApplications = filteredApplications.filter(app => app.status === params.status_filter)
  }
  if (params.applicant_name) {
    filteredApplications = filteredApplications.filter(app => 
      app.applicant_name.includes(params.applicant_name)
    )
  }
  if (params.application_id) {
    filteredApplications = filteredApplications.filter(app => 
      app.application_id.includes(params.application_id)
    )
  }
  
  const page = params.page || 1
  const limit = params.limit || 20
  const start = (page - 1) * limit
  const end = start + limit
  
  return {
    data: filteredApplications.slice(start, end),
    total: filteredApplications.length,
    page,
    limit
  }
}

// 模拟申请详情
export const mockApplicationDetail = async (applicationId: string): Promise<ApplicationDetail> => {
  await delay(800)
  
  return {
    application_id: applicationId,
    amount: 50000,
    term_months: 12,
    purpose: '购买农业设备',
    status: '待人工复核',
    submitted_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    updated_at: new Date(Date.now() - 30 * 60 * 1000).toISOString(),
    approved_amount: undefined,
    applicant_details: {
      user_id: 'U001',
      real_name: '张三',
      id_card_number: '320123199001011234',
      phone: '13812345678',
      address: '江苏省南京市玄武区某某街道123号',
      bank_account: '6222081234567890123'
    },
    ai_analysis_report: {
      overall_risk_score: 25,
      risk_factors: ['征信记录良好', '收入稳定'],
      data_verification_results: [
        { item: '身份证验证', result: '通过' },
        { item: '银行流水验证', result: '通过' },
        { item: '收入证明验证', result: '通过' }
      ],
      suggestion: '建议批准，申请人信用良好，收入稳定，风险较低'
    },
    uploaded_documents_details: [
      {
        file_id: 'F001',
        doc_type: 'id_card_front',
        file_url: '/api/files/F001',
        ocr_result: '身份证信息已提取'
      },
      {
        file_id: 'F002',
        doc_type: 'income_proof',
        file_url: '/api/files/F002',
        ocr_result: '月收入8000元'
      }
    ],
    history: [
      {
        status: '已提交',
        operator: '张三',
        timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
      },
      {
        status: 'AI_审批中',
        operator: 'AI系统',
        timestamp: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000 + 5 * 60 * 1000).toISOString()
      },
      {
        status: '待人工复核',
        operator: 'AI系统',
        timestamp: new Date(Date.now() - 30 * 60 * 1000).toISOString()
      }
    ]
  }
}

// 模拟用户列表
export const mockUsers = async (params: any): Promise<PaginationResponse<AdminUser>> => {
  await delay(600)
  
  const allUsers: AdminUser[] = [
    {
      admin_user_id: 'U001',
      username: 'admin',
      role: 'ADMIN',
      display_name: '系统管理员',
      email: 'admin@example.com',
      status: 0,
      created_at: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString()
    },
    {
      admin_user_id: 'U002',
      username: 'reviewer',
      role: '审批员',
      display_name: '李审批员',
      email: 'reviewer@example.com',
      status: 0,
      created_at: new Date(Date.now() - 15 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString()
    },
    {
      admin_user_id: 'U003',
      username: 'reviewer2',
      role: '审批员',
      display_name: '王审批员',
      email: 'reviewer2@example.com',
      status: 1,
      created_at: new Date(Date.now() - 10 * 24 * 60 * 60 * 1000).toISOString(),
      updated_at: new Date().toISOString()
    }
  ]
  
  let filteredUsers = allUsers
  if (params.role) {
    filteredUsers = filteredUsers.filter(user => user.role === params.role)
  }
  
  const page = params.page || 1
  const limit = params.limit || 20
  const start = (page - 1) * limit
  const end = start + limit
  
  return {
    data: filteredUsers.slice(start, end),
    total: filteredUsers.length,
    page,
    limit
  }
}

// 模拟操作日志
export const mockLogs = async (params: any): Promise<PaginationResponse<OperationLog>> => {
  await delay(600)
  
  const allLogs: OperationLog[] = [
    {
      id: '1',
      operator_id: 'U001',
      operator_name: '系统管理员',
      action: '登录',
      target: '管理后台',
      result: '成功',
      ip_address: '192.168.1.100',
      user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      occurred_at: new Date(Date.now() - 30 * 60 * 1000).toISOString()
    },
    {
      id: '2',
      operator_id: 'U001',
      operator_name: '系统管理员',
      action: '审批申请',
      target: '申请编号: APP001',
      result: '成功',
      ip_address: '192.168.1.100',
      user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      occurred_at: new Date(Date.now() - 60 * 60 * 1000).toISOString()
    },
    {
      id: '3',
      operator_id: 'U002',
      operator_name: '李审批员',
      action: '查看申请详情',
      target: '申请编号: APP002',
      result: '成功',
      ip_address: '192.168.1.101',
      user_agent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36',
      occurred_at: new Date(Date.now() - 2 * 60 * 60 * 1000).toISOString()
    }
  ]
  
  const page = params.page || 1
  const limit = params.limit || 20
  const start = (page - 1) * limit
  const end = start + limit
  
  return {
    data: allLogs.slice(start, end),
    total: allLogs.length,
    page,
    limit
  }
}

// 模拟系统配置
export const mockConfigs = async (): Promise<SystemConfig[]> => {
  await delay(500)
  
  return [
    {
      config_key: 'ai_risk_threshold',
      config_value: '70',
      description: 'AI风险评分阈值',
      updated_at: new Date().toISOString()
    },
    {
      config_key: 'auto_approval_limit',
      config_value: '50000',
      description: '自动批准金额上限',
      updated_at: new Date().toISOString()
    },
    {
      config_key: 'max_application_amount',
      config_value: '500000',
      description: '最大申请金额',
      updated_at: new Date().toISOString()
    }
  ]
}

// 模拟提交审批
export const mockSubmitReview = async (applicationId: string, data: any): Promise<void> => {
  await delay(1500)
  console.log(`提交审批: ${applicationId}`, data)
  // 模拟成功
}

// 模拟其他操作
export const mockCreateUser = async (data: any): Promise<AdminUser> => {
  await delay(1000)
  return {
    admin_user_id: 'U' + Date.now(),
    username: data.username,
    role: data.role,
    display_name: data.display_name,
    email: data.email,
    status: 0,
    created_at: new Date().toISOString(),
    updated_at: new Date().toISOString()
  }
}

export const mockUpdateUserStatus = async (userId: string, status: number): Promise<void> => {
  await delay(800)
  console.log(`更新用户状态: ${userId} -> ${status}`)
}

export const mockToggleAI = async (enabled: boolean): Promise<void> => {
  await delay(1000)
  console.log(`切换AI审批: ${enabled}`)
}

export const mockUpdateConfig = async (key: string, value: string): Promise<void> => {
  await delay(800)
  console.log(`更新配置: ${key} = ${value}`)
} 