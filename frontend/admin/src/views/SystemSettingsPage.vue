<template>
  <div class="system-settings">
    <div class="page-header">
      <h1 class="page-title">系统设置</h1>
      <div class="header-desc">基本设置、安全设置、通知设置、系统参数</div>
    </div>

    <div class="page-content">
      <!-- 基本设置 -->
      <div class="settings-card">
        <div class="card-header">
          <h2 class="card-title">基本设置</h2>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">系统名称</label>
            <input type="text" class="form-input" v-model="basicSettings.systemName" />
          </div>
          <div class="form-group">
            <label class="form-label">系统描述</label>
            <textarea class="form-textarea" v-model="basicSettings.systemDesc" rows="3"></textarea>
          </div>
          <div class="form-group">
            <label class="form-label">管理员邮箱</label>
            <input type="email" class="form-input" v-model="basicSettings.adminEmail" />
          </div>
          <div class="form-group">
            <label class="form-label">默认语言</label>
            <select class="form-select" v-model="basicSettings.language">
              <option value="zh-CN">简体中文</option>
              <option value="en-US">English</option>
            </select>
          </div>
          <div class="form-actions">
            <button class="save-btn" @click="saveBasicSettings">保存设置</button>
          </div>
        </div>
      </div>

      <!-- 安全设置 -->
      <div class="settings-card">
        <div class="card-header">
          <h2 class="card-title">安全设置</h2>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">密码最小长度</label>
            <div class="slider-group">
              <input 
                type="range" 
                min="6" 
                max="20" 
                step="1" 
                class="form-slider" 
                v-model="securitySettings.minPasswordLength" 
              />
              <span class="slider-value">{{ securitySettings.minPasswordLength }} 个字符</span>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">密码复杂度要求</label>
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="securitySettings.requireUppercase" />
                <span>必须包含大写字母</span>
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="securitySettings.requireLowercase" />
                <span>必须包含小写字母</span>
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="securitySettings.requireNumbers" />
                <span>必须包含数字</span>
              </label>
              <label class="checkbox-label">
                <input type="checkbox" v-model="securitySettings.requireSpecialChars" />
                <span>必须包含特殊字符</span>
              </label>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">登录失败尝试次数</label>
            <div class="slider-group">
              <input 
                type="range" 
                min="3" 
                max="10" 
                step="1" 
                class="form-slider" 
                v-model="securitySettings.maxLoginAttempts" 
              />
              <span class="slider-value">{{ securitySettings.maxLoginAttempts }} 次</span>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">双因素认证</label>
            <div class="switch-container">
              <label class="switch">
                <input type="checkbox" v-model="securitySettings.enable2FA" />
                <span class="slider round"></span>
              </label>
              <span class="switch-label">{{ securitySettings.enable2FA ? '已开启' : '已关闭' }}</span>
            </div>
          </div>
          <div class="form-actions">
            <button class="save-btn" @click="saveSecuritySettings">保存设置</button>
          </div>
        </div>
      </div>

      <!-- 通知设置 -->
      <div class="settings-card">
        <div class="card-header">
          <h2 class="card-title">通知设置</h2>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">邮件通知</label>
            <div class="switch-container">
              <label class="switch">
                <input type="checkbox" v-model="notificationSettings.enableEmailNotification" />
                <span class="slider round"></span>
              </label>
              <span class="switch-label">{{ notificationSettings.enableEmailNotification ? '已开启' : '已关闭' }}</span>
            </div>
          </div>
          <div class="form-group" v-if="notificationSettings.enableEmailNotification">
            <label class="form-label">SMTP服务器</label>
            <input type="text" class="form-input" v-model="notificationSettings.smtpServer" />
          </div>
          <div class="form-group" v-if="notificationSettings.enableEmailNotification">
            <label class="form-label">SMTP端口</label>
            <input type="number" class="form-input" v-model="notificationSettings.smtpPort" />
          </div>
          <div class="form-group">
            <label class="form-label">短信通知</label>
            <div class="switch-container">
              <label class="switch">
                <input type="checkbox" v-model="notificationSettings.enableSmsNotification" />
                <span class="slider round"></span>
              </label>
              <span class="switch-label">{{ notificationSettings.enableSmsNotification ? '已开启' : '已关闭' }}</span>
            </div>
          </div>
          <div class="form-group" v-if="notificationSettings.enableSmsNotification">
            <label class="form-label">短信服务提供商</label>
            <select class="form-select" v-model="notificationSettings.smsProvider">
              <option value="aliyun">阿里云</option>
              <option value="tencent">腾讯云</option>
              <option value="netease">网易云信</option>
            </select>
          </div>
          <div class="form-actions">
            <button class="save-btn" @click="saveNotificationSettings">保存设置</button>
          </div>
        </div>
      </div>

      <!-- 系统参数 -->
      <div class="settings-card">
        <div class="card-header">
          <h2 class="card-title">系统参数</h2>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">默认每页记录数</label>
            <select class="form-select" v-model="systemParams.defaultPageSize">
              <option value="10">10条/页</option>
              <option value="20">20条/页</option>
              <option value="50">50条/页</option>
              <option value="100">100条/页</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">系统日志保留天数</label>
            <input type="number" class="form-input" v-model="systemParams.logRetentionDays" />
          </div>
          <div class="form-group">
            <label class="form-label">自动备份</label>
            <div class="switch-container">
              <label class="switch">
                <input type="checkbox" v-model="systemParams.enableAutoBackup" />
                <span class="slider round"></span>
              </label>
              <span class="switch-label">{{ systemParams.enableAutoBackup ? '已开启' : '已关闭' }}</span>
            </div>
          </div>
          <div class="form-group" v-if="systemParams.enableAutoBackup">
            <label class="form-label">备份频率</label>
            <select class="form-select" v-model="systemParams.backupFrequency">
              <option value="daily">每日</option>
              <option value="weekly">每周</option>
              <option value="monthly">每月</option>
            </select>
          </div>
          <div class="form-actions">
            <button class="save-btn" @click="saveSystemParams">保存设置</button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// 基本设置
const basicSettings = ref({
  systemName: '数字惠农后台管理系统',
  systemDesc: '为农村金融服务提供智能化、数字化管理平台',
  adminEmail: 'admin@huinongfinancial.com',
  language: 'zh-CN'
})

// 安全设置
const securitySettings = ref({
  minPasswordLength: 8,
  requireUppercase: true,
  requireLowercase: true,
  requireNumbers: true,
  requireSpecialChars: false,
  maxLoginAttempts: 5,
  enable2FA: false
})

// 通知设置
const notificationSettings = ref({
  enableEmailNotification: true,
  smtpServer: 'smtp.example.com',
  smtpPort: 587,
  enableSmsNotification: true,
  smsProvider: 'aliyun'
})

// 系统参数
const systemParams = ref({
  defaultPageSize: 20,
  logRetentionDays: 90,
  enableAutoBackup: true,
  backupFrequency: 'daily'
})

// 保存基本设置
const saveBasicSettings = () => {
  alert('基本设置已保存')
}

// 保存安全设置
const saveSecuritySettings = () => {
  alert('安全设置已保存')
}

// 保存通知设置
const saveNotificationSettings = () => {
  alert('通知设置已保存')
}

// 保存系统参数
const saveSystemParams = () => {
  alert('系统参数已保存')
}
</script>

<style scoped>
.system-settings {
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

.settings-card {
  background-color: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  margin-bottom: 24px;
}

.card-header {
  padding: 16px;
  border-bottom: 1px solid #f0f0f0;
}

.card-title {
  margin: 0;
  font-size: 16px;
  color: #333;
}

.card-body {
  padding: 16px;
}

.form-group {
  margin-bottom: 16px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #606266;
}

.form-input,
.form-textarea,
.form-select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  font-size: 14px;
  color: #606266;
}

.form-input:focus,
.form-textarea:focus,
.form-select:focus {
  border-color: #4285f4;
  outline: none;
}

.form-textarea {
  resize: vertical;
}

.form-actions {
  margin-top: 24px;
  text-align: right;
}

.save-btn {
  padding: 8px 16px;
  background-color: #4285f4;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.save-btn:hover {
  background-color: #3367d6;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.checkbox-label input {
  margin-right: 8px;
}

.slider-group {
  display: flex;
  align-items: center;
}

.form-slider {
  flex: 1;
  margin-right: 16px;
}

.slider-value {
  width: 80px;
  text-align: right;
}

.switch-container {
  display: flex;
  align-items: center;
}

.switch {
  position: relative;
  display: inline-block;
  width: 48px;
  height: 24px;
  margin-right: 12px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #ccc;
  transition: .4s;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: .4s;
}

input:checked + .slider {
  background-color: #4285f4;
}

input:focus + .slider {
  box-shadow: 0 0 1px #4285f4;
}

input:checked + .slider:before {
  transform: translateX(24px);
}

.slider.round {
  border-radius: 24px;
}

.slider.round:before {
  border-radius: 50%;
}

.switch-label {
  font-size: 14px;
  color: #606266;
}
</style> 