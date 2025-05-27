<template>
  <div class="users-management">
    <div class="page-header">
      <h1 class="page-title">用户管理</h1>
      <div class="header-desc">用户ID、用户名、手机号、注册时间、账户状态、操作</div>
    </div>

    <div class="page-content">
      <div class="table-actions">
        <div class="left-actions">
          <button class="action-btn" @click="addUser">
            <i class="icon-add"></i>
            <span>添加用户</span>
          </button>
          <button class="action-btn" @click="exportUserList">
            <i class="icon-export"></i>
            <span>导出Excel</span>
          </button>
        </div>
        <div class="right-actions">
          <div class="search-box">
            <input type="text" placeholder="搜索用户..." class="search-input" v-model="searchQuery" />
            <button class="search-btn" @click="searchUsers">
              <i class="icon-search"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>用户ID</th>
              <th>用户名</th>
              <th>手机号</th>
              <th>用户类型</th>
              <th>注册时间</th>
              <th>账户状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="user in userList" :key="user.id">
              <td>{{ user.id }}</td>
              <td>{{ user.username }}</td>
              <td>{{ user.phone }}</td>
              <td>{{ user.type }}</td>
              <td>{{ user.registerDate }}</td>
              <td>
                <span 
                  class="status-tag"
                  :class="{
                    'status-active': user.status === '正常',
                    'status-locked': user.status === '已锁定',
                    'status-inactive': user.status === '未激活'
                  }"
                >
                  {{ user.status }}
                </span>
              </td>
              <td class="action-column">
                <button class="table-btn view-btn" @click="viewUserDetails(user.id)">查看</button>
                <button class="table-btn edit-btn" @click="editUser(user.id)">编辑</button>
                <button class="table-btn" 
                  :class="user.status === '正常' ? 'lock-btn' : 'unlock-btn'"
                  @click="toggleUserStatus(user.id)">
                  {{ user.status === '正常' ? '锁定' : '解锁' }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <div class="pagination-info">
          共 {{ total }} 条
        </div>
        <div class="pagination-selector">
          <select v-model="pageSize">
            <option value="10">10条/页</option>
            <option value="20">20条/页</option>
            <option value="50">50条/页</option>
          </select>
        </div>
        <div class="pagination-pages">
          <button 
            class="page-btn"
            :disabled="currentPage === 1"
            @click="handleCurrentChange(currentPage - 1)"
          >
            &lt;
          </button>
          <button
            class="page-btn"
            :class="{ 'active': currentPage === 1 }"
            @click="handleCurrentChange(1)"
          >
            1
          </button>
          <button
            class="page-btn"
            @click="handleCurrentChange(currentPage + 1)"
            :disabled="currentPage * pageSize >= total"
          >
            &gt;
          </button>
          <div class="pagination-goto">
            跳转
            <input type="number" v-model="currentPage" min="1" :max="Math.ceil(total / pageSize)" class="goto-input" />
            页
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 用户数据
const userList = ref([
  {
    id: 'U0001',
    username: '张三',
    phone: '13812345678',
    type: '农户',
    registerDate: '2023-01-15',
    status: '正常'
  },
  {
    id: 'U0002',
    username: '李四',
    phone: '13876543210',
    type: '农户',
    registerDate: '2023-02-20',
    status: '正常'
  },
  {
    id: 'U0003',
    username: '王五',
    phone: '13598765432',
    type: '设备商',
    registerDate: '2023-03-10',
    status: '已锁定'
  },
  {
    id: 'U0004',
    username: '赵六',
    phone: '13765432109',
    type: '农户',
    registerDate: '2023-04-05',
    status: '未激活'
  },
  {
    id: 'U0005',
    username: '孙七',
    phone: '13987654321',
    type: '农技专家',
    registerDate: '2023-05-12',
    status: '正常'
  },
  {
    id: 'U0006',
    username: '周八',
    phone: '13654321098',
    type: '农户',
    registerDate: '2023-06-18',
    status: '正常'
  },
  {
    id: 'U0007',
    username: '吴九',
    phone: '13543210987',
    type: '金融顾问',
    registerDate: '2023-07-22',
    status: '正常'
  },
  {
    id: 'U0008',
    username: '郑十',
    phone: '13432109876',
    type: '农户',
    registerDate: '2023-08-30',
    status: '已锁定'
  }
])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(8)
const searchQuery = ref('')

// 分页处理
const handleCurrentChange = (page) => {
  currentPage.value = page
}

// 查看用户详情
const viewUserDetails = (id) => {
  alert(`查看用户ID: ${id} 的详情`)
}

// 编辑用户
const editUser = (id) => {
  alert(`编辑用户ID: ${id} 的信息`)
}

// 切换用户状态
const toggleUserStatus = (id) => {
  const user = userList.value.find(user => user.id === id)
  if (user) {
    if (user.status === '正常') {
      user.status = '已锁定'
      alert(`用户 ${user.username} 已被锁定`)
    } else {
      user.status = '正常'
      alert(`用户 ${user.username} 已解锁`)
    }
  }
}

// 添加用户
const addUser = () => {
  alert('添加新用户')
}

// 导出用户列表
const exportUserList = () => {
  alert('导出用户列表')
}

// 搜索用户
const searchUsers = () => {
  alert(`搜索关键词: ${searchQuery.value}`)
}
</script>

<style scoped>
.users-management {
  background-color: #ffffff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.page-header {
  padding: 16px 24px;
  background: linear-gradient(90deg, #4285f4 0%, #34a853 100%);
  color: #fff;
  border-radius: 4px 4px 0 0;
}

.page-title {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.header-desc {
  margin-top: 8px;
  font-size: 14px;
  opacity: 0.8;
}

.page-content {
  padding: 24px;
}

.table-actions {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
}

.left-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 8px 16px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.3s;
}

.action-btn:hover {
  color: #4285f4;
  border-color: #4285f4;
}

.action-btn span {
  margin-left: 4px;
}

.search-box {
  display: flex;
  align-items: center;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

.search-input {
  padding: 8px 12px;
  border: none;
  outline: none;
  width: 200px;
  font-size: 14px;
}

.search-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 36px;
  background-color: #f5f7fa;
  border: none;
  cursor: pointer;
}

.table-container {
  width: 100%;
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  border-spacing: 0;
  font-size: 14px;
}

.data-table th, .data-table td {
  padding: 12px 16px;
  text-align: left;
  border-bottom: 1px solid #ebeef5;
}

.data-table th {
  background-color: #f5f7fa;
  color: #606266;
  font-weight: 500;
}

.status-tag {
  display: inline-block;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-active {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-locked {
  background-color: #fff2f0;
  color: #ff4d4f;
}

.status-inactive {
  background-color: #f4f4f5;
  color: #909399;
}

.action-column {
  white-space: nowrap;
}

.table-btn {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
  margin-right: 4px;
}

.view-btn {
  background-color: #e6f7ff;
  color: #1890ff;
}

.edit-btn {
  background-color: #f6ffed;
  color: #52c41a;
}

.lock-btn {
  background-color: #fff2f0;
  color: #ff4d4f;
}

.unlock-btn {
  background-color: #f4f4f5;
  color: #909399;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-top: 16px;
  font-size: 14px;
}

.pagination-info {
  margin-right: 16px;
  color: #606266;
}

.pagination-selector {
  margin-right: 16px;
}

.pagination-selector select {
  padding: 4px 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: #fff;
  outline: none;
}

.pagination-pages {
  display: flex;
  align-items: center;
}

.page-btn {
  min-width: 32px;
  height: 32px;
  padding: 0 4px;
  margin: 0 4px;
  border: 1px solid #dcdfe6;
  background-color: #fff;
  color: #606266;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.page-btn.active {
  color: #1890ff;
  border-color: #1890ff;
}

.page-btn:disabled {
  color: #c0c4cc;
  cursor: not-allowed;
}

.pagination-goto {
  margin-left: 16px;
  display: flex;
  align-items: center;
}

.goto-input {
  width: 40px;
  height: 32px;
  margin: 0 8px;
  padding: 0 8px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  text-align: center;
}

/* 图标样式 */
.icon-add, .icon-export, .icon-search {
  display: inline-block;
  width: 16px;
  height: 16px;
  background-color: currentColor;
  mask-repeat: no-repeat;
  mask-position: center;
  mask-size: contain;
  -webkit-mask-repeat: no-repeat;
  -webkit-mask-position: center;
  -webkit-mask-size: contain;
}

.icon-add {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z'/%3E%3C/svg%3E");
}

.icon-export {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
}

.icon-search {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
}
</style> 