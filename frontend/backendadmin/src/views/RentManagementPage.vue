<script setup>
import { ref, onMounted } from 'vue'

// 农机租赁申请列表数据
const rentApplications = ref([
  {
    id: 1,
    subject: '拖拉机X100',
    startDate: '2024-04-10 00:00:00',
    endDate: '2024-04-20 00:00:00',
    equipmentId: '推拉机X100',
    applicantId: 'xiaofeifei',
    status: '待审核'
  },
  {
    id: 2,
    subject: '收割机Y200',
    startDate: '2024-04-15 00:00:00',
    endDate: '2024-04-25 00:00:00',
    equipmentId: '收割机Y200',
    applicantId: 'wanglei',
    status: '申请通过'
  },
  {
    id: 3,
    subject: '播种机Z300',
    startDate: '2024-04-05 00:00:00',
    endDate: '2024-04-10 00:00:00',
    equipmentId: '播种机Z300',
    applicantId: 'zhanghua',
    status: '申请拒绝'
  },
  {
    id: 4,
    subject: '喷洒机A400',
    startDate: '2024-04-12 00:00:00',
    endDate: '2024-04-18 00:00:00',
    equipmentId: '喷洒机A400',
    applicantId: 'chenqiang',
    status: '待审核'
  },
  {
    id: 5,
    subject: '粮食烘干机B500',
    startDate: '2024-04-20 00:00:00',
    endDate: '2024-04-30 00:00:00',
    equipmentId: '粮食烘干机B500',
    applicantId: 'sunting',
    status: '申请通过'
  },
  {
    id: 6,
    subject: '起垄机C600',
    startDate: '2024-04-22 00:00:00',
    endDate: '2024-05-02 00:00:00',
    equipmentId: '起垄机C600',
    applicantId: 'liuxin',
    status: '待审核'
  }
])

// 当前页和页大小
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(6)

// 分页处理
const handleCurrentChange = (page) => {
  currentPage.value = page
  // 在实际应用中，这里需要重新加载数据
}

// 导出Excel
const exportToExcel = () => {
  // 实际应用中，这里应该调用导出Excel的API
  alert('导出Excel功能')
}

// 导入Excel
const importExcel = () => {
  // 实际应用中，这里应该调用导入Excel的API
  alert('导入Excel功能')
}

// 下载Excel模板
const downloadTemplate = () => {
  // 实际应用中，这里应该提供一个Excel模板下载
  alert('下载Excel模板')
}

// 查看详情
const viewDetails = (id) => {
  alert(`查看申请ID: ${id} 的详情`)
}

// 编辑申请
const editApplication = (id) => {
  alert(`编辑申请ID: ${id}`)
}

// 删除申请
const deleteApplication = (id) => {
  if (confirm(`确定要删除ID: ${id} 的申请吗？`)) {
    rentApplications.value = rentApplications.value.filter(app => app.id !== id)
    total.value--
  }
}
</script>

<template>
  <div class="rent-management">
    <div class="page-header">
      <h1 class="page-title">农机租赁申请</h1>
      <div class="header-desc">主题、设备ID、申请用户ID、租赁开始日期、租赁结束日期、申请状态。</div>
    </div>

    <div class="page-content">
      <div class="table-actions">
        <div class="left-actions">
          <button class="action-btn" @click="exportToExcel">
            <i class="icon-add"></i>
            <span>添加</span>
          </button>
          <button class="action-btn" @click="importExcel">
            <i class="icon-upload"></i>
            <span>导入Excel</span>
          </button>
          <button class="action-btn" @click="exportToExcel">
            <i class="icon-export"></i>
            <span>导出Excel</span>
          </button>
          <button class="action-btn" @click="downloadTemplate">
            <i class="icon-download"></i>
            <span>下载Excel模板</span>
          </button>
        </div>
        <div class="right-actions">
          <div class="search-box">
            <input type="text" placeholder="高级搜索" class="search-input" />
            <button class="search-btn">
              <i class="icon-search"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>主题</th>
              <th>租赁开始日期</th>
              <th>租赁结束日期</th>
              <th>设备ID</th>
              <th>申请用户ID</th>
              <th>申请状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in rentApplications" :key="item.id">
              <td>{{ item.subject }}</td>
              <td>{{ item.startDate }}</td>
              <td>{{ item.endDate }}</td>
              <td>{{ item.equipmentId }}</td>
              <td>{{ item.applicantId }}</td>
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
                <button class="table-btn view-btn" @click="viewDetails(item.id)">详情</button>
                <button class="table-btn edit-btn" @click="editApplication(item.id)">编辑</button>
                <button class="table-btn delete-btn" @click="deleteApplication(item.id)">删除</button>
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

<style scoped>
.rent-management {
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
  background-color: #e6f7ff;
  color: #1890ff;
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

.edit-btn {
  background-color: #f6ffed;
  color: #52c41a;
}

.delete-btn {
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
.icon-add, .icon-upload, .icon-export, .icon-download, .icon-search {
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

.icon-upload {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M9 16h6v-6h4l-7-7-7 7h4v6zm-4 2h14v2H5v-2z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M9 16h6v-6h4l-7-7-7 7h4v6zm-4 2h14v2H5v-2z'/%3E%3C/svg%3E");
}

.icon-export {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 12v7H5v-7H3v7c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2v-7h-2zm-6 .67l2.59-2.58L17 11.5l-5 5-5-5 1.41-1.41L11 12.67V3h2v9.67z'/%3E%3C/svg%3E");
}

.icon-download {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z'/%3E%3C/svg%3E");
}

.icon-search {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
}
</style> 