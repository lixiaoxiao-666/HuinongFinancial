<template>
  <div class="rent-application">
    <div class="page-header">
      <h1 class="page-title">租赁申请管理</h1>
      <div class="header-desc">申请编号、申请人、联系方式、申请设备、申请日期、归还日期、状态</div>
    </div>

    <div class="page-content">
      <div class="table-actions">
        <div class="left-actions">
          <button class="action-btn" @click="exportToExcel">
            <i class="icon-export"></i>
            <span>导出Excel</span>
          </button>
        </div>
        <div class="right-actions">
          <div class="search-box">
            <input type="text" placeholder="搜索申请..." class="search-input" v-model="searchQuery" />
            <button class="search-btn" @click="searchApplications">
              <i class="icon-search"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>申请编号</th>
              <th>申请人</th>
              <th>联系方式</th>
              <th>申请设备</th>
              <th>申请日期</th>
              <th>归还日期</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in applicationList" :key="item.id">
              <td>{{ item.id }}</td>
              <td>{{ item.applicant }}</td>
              <td>{{ item.phone }}</td>
              <td>{{ item.equipment }}</td>
              <td>{{ item.applyDate }}</td>
              <td>{{ item.returnDate }}</td>
              <td>
                <span 
                  class="status-tag"
                  :class="{
                    'status-pending': item.status === '待审核',
                    'status-approved': item.status === '申请通过',
                    'status-rejected': item.status === '申请拒绝'
                  }"
                >
                  {{ item.status }}
                </span>
              </td>
              <td class="action-column">
                <button class="table-btn view-btn" @click="viewDetails(item.id)">查看</button>
                <button class="table-btn approve-btn" @click="approveApplication(item.id)" v-if="item.status === '待审核'">审核</button>
                <button class="table-btn reject-btn" @click="rejectApplication(item.id)" v-if="item.status === '待审核'">拒绝</button>
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

// 租赁申请数据
const applicationList = ref([
  {
    id: 'RA-20240601-001',
    applicant: '张三',
    phone: '13812345678',
    equipment: '拖拉机X100',
    applyDate: '2024-06-01',
    returnDate: '2024-06-15',
    status: '待审核'
  },
  {
    id: 'RA-20240529-002',
    applicant: '李四',
    phone: '13898765432',
    equipment: '收割机Y200',
    applyDate: '2024-05-29',
    returnDate: '2024-06-10',
    status: '申请通过'
  },
  {
    id: 'RA-20240527-003',
    applicant: '王五',
    phone: '13756781234',
    equipment: '播种机Z300',
    applyDate: '2024-05-27',
    returnDate: '2024-06-05',
    status: '申请拒绝'
  },
  {
    id: 'RA-20240526-004',
    applicant: '赵六',
    phone: '13678901234',
    equipment: '喷洒机A400',
    applyDate: '2024-05-26',
    returnDate: '2024-06-08',
    status: '待审核'
  },
  {
    id: 'RA-20240525-005',
    applicant: '孙七',
    phone: '13567890123',
    equipment: '粮食烘干机B500',
    applyDate: '2024-05-25',
    returnDate: '2024-06-12',
    status: '申请通过'
  }
])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(5)
const searchQuery = ref('')

// 分页处理
const handleCurrentChange = (page) => {
  currentPage.value = page
}

// 查看详情
const viewDetails = (id) => {
  alert(`查看申请ID: ${id} 的详情`)
}

// 审核通过
const approveApplication = (id) => {
  const application = applicationList.value.find(item => item.id === id)
  if (application) {
    application.status = '申请通过'
  }
}

// 拒绝申请
const rejectApplication = (id) => {
  const application = applicationList.value.find(item => item.id === id)
  if (application) {
    application.status = '申请拒绝'
  }
}

// 导出Excel
const exportToExcel = () => {
  alert('导出Excel功能')
}

// 搜索申请
const searchApplications = () => {
  alert(`搜索关键词: ${searchQuery.value}`)
}
</script>

<style scoped>
.rent-application {
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

.status-pending {
  background-color: #fff7e6;
  color: #fa8c16;
}

.status-approved {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-rejected {
  background-color: #fff2f0;
  color: #ff4d4f;
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

.approve-btn {
  background-color: #f6ffed;
  color: #52c41a;
}

.reject-btn {
  background-color: #fff2f0;
  color: #ff4d4f;
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
.icon-export, .icon-search {
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

.icon-export {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
}

.icon-search {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
}
</style>
 