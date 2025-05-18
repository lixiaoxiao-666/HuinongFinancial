<template>
  <div class="machinery-management">
    <div class="page-header">
      <h1 class="page-title">农机设备管理</h1>
      <div class="header-desc">设备编号、设备类型、状态、创建日期、操作</div>
    </div>

    <div class="page-content">
      <div class="table-actions">
        <div class="left-actions">
          <button class="action-btn">
            <i class="icon-add"></i>
            <span>添加设备</span>
          </button>
          <button class="action-btn">
            <i class="icon-upload"></i>
            <span>导入设备</span>
          </button>
          <button class="action-btn">
            <i class="icon-export"></i>
            <span>导出Excel</span>
          </button>
        </div>
        <div class="right-actions">
          <div class="search-box">
            <input type="text" placeholder="搜索设备..." class="search-input" />
            <button class="search-btn">
              <i class="icon-search"></i>
            </button>
          </div>
        </div>
      </div>

      <div class="machinery-grid">
        <div v-for="item in machineryList" :key="item.id" class="machinery-card">
          <div class="machinery-image">
            <img :src="item.imageUrl" :alt="item.name" />
          </div>
          <div class="machinery-details">
            <h3 class="machinery-name">{{ item.name }}</h3>
            <div class="machinery-info">
              <span class="info-label">编号:</span>
              <span class="info-value">{{ item.id }}</span>
            </div>
            <div class="machinery-info">
              <span class="info-label">状态:</span>
              <span class="info-value">
                <span 
                  class="status-tag"
                  :class="{
                    'status-available': item.status === '可用',
                    'status-inuse': item.status === '使用中',
                    'status-maintenance': item.status === '维修中'
                  }"
                >
                  {{ item.status }}
                </span>
              </span>
            </div>
            <div class="machinery-actions">
              <button class="table-btn view-btn">详情</button>
              <button class="table-btn edit-btn">编辑</button>
              <button class="table-btn delete-btn">删除</button>
            </div>
          </div>
        </div>
      </div>

      <div class="pagination">
        <div class="pagination-info">
          共 {{ total }} 条
        </div>
        <div class="pagination-selector">
          <select v-model="pageSize">
            <option value="8">8条/页</option>
            <option value="16">16条/页</option>
            <option value="24">24条/页</option>
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

// 农机设备数据
const machineryList = ref([
  {
    id: 'TL-X100',
    name: '拖拉机X100',
    status: '可用',
    createDate: '2023-05-10',
    imageUrl: '/assets/machinery/tractor.jpg'
  },
  {
    id: 'SG-Y200',
    name: '收割机Y200',
    status: '使用中',
    createDate: '2023-06-15',
    imageUrl: '/assets/machinery/harvester.jpg'
  },
  {
    id: 'BZ-Z300',
    name: '播种机Z300',
    status: '可用',
    createDate: '2023-07-20',
    imageUrl: '/assets/machinery/seeder.jpg'
  },
  {
    id: 'PS-A400',
    name: '喷洒机A400',
    status: '维修中',
    createDate: '2023-08-25',
    imageUrl: '/assets/machinery/sprayer.jpg'
  },
  {
    id: 'HG-B500',
    name: '粮食烘干机B500',
    status: '可用',
    createDate: '2023-09-30',
    imageUrl: '/assets/machinery/dryer.jpg'
  },
  {
    id: 'QL-C600',
    name: '起垄机C600',
    status: '使用中',
    createDate: '2023-10-05',
    imageUrl: '/assets/machinery/ridger.jpg'
  },
  {
    id: 'NY-D700',
    name: '农药喷雾机D700',
    status: '可用',
    createDate: '2023-11-10',
    imageUrl: '/assets/machinery/pesticide_sprayer.jpg'
  },
  {
    id: 'GW-E800',
    name: '耕田机E800',
    status: '可用',
    createDate: '2023-12-15',
    imageUrl: '/assets/machinery/plough.jpg'
  }
])

// 分页相关
const currentPage = ref(1)
const pageSize = ref(8)
const total = ref(8)

// 分页处理
const handleCurrentChange = (page) => {
  currentPage.value = page
}
</script>

<style scoped>
.machinery-management {
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
  margin-bottom: 24px;
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

/* 农机设备网格布局 */
.machinery-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 24px;
}

.machinery-card {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  overflow: hidden;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.machinery-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
}

.machinery-image {
  height: 160px;
  overflow: hidden;
  background-color: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
}

.machinery-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s;
}

.machinery-card:hover .machinery-image img {
  transform: scale(1.05);
}

.machinery-details {
  padding: 16px;
}

.machinery-name {
  margin: 0 0 12px;
  font-size: 16px;
  font-weight: 500;
  color: #333;
}

.machinery-info {
  display: flex;
  margin-bottom: 8px;
  font-size: 14px;
  color: #666;
}

.info-label {
  width: 45px;
  color: #888;
}

.status-tag {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.status-available {
  background-color: #f6ffed;
  color: #52c41a;
}

.status-inuse {
  background-color: #e6f7ff;
  color: #1890ff;
}

.status-maintenance {
  background-color: #fff2f0;
  color: #ff4d4f;
}

.machinery-actions {
  display: flex;
  margin-top: 12px;
  gap: 8px;
}

.table-btn {
  padding: 4px 8px;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
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
.icon-add, .icon-upload, .icon-export, .icon-search {
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

.icon-search {
  mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 5 1.5-1.5-5-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
  -webkit-mask-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24'%3E%3Cpath d='M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 5 1.5-1.5-5-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z'/%3E%3C/svg%3E");
}
</style> 