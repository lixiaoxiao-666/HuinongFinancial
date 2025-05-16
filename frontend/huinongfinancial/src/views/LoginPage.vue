<script setup lang="ts">
// 导入组合式API
import { reactive, ref } from 'vue'
// 导入 element-plus 表单实例和表单规则类型
import type { FormInstance, FormRules } from 'element-plus'

// 定义ref对象表单实例
const FormRef = ref<FormInstance>()

// 定义表单数据
const FormData = reactive({
  username: '',
  password: '',
})

// 定义表单验证
const validateUsername = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入账号'))
  } else {
    callback()
  }
}

// 定义表单验证
const validatePassword = (rule: any, value: any, callback: any) => {
  if (value === '') {
    callback(new Error('请输入密码'))
  } else {
    callback()
  }
}

// 定义表单规则
const rules = reactive<FormRules<typeof FormData>>({
    username: [{ validator: validateUsername, trigger: 'blur' }],
    password: [{ validator: validatePassword, trigger: 'blur' }],
})

// 定义表单提交
const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.validate((valid) => {
    if (valid) {
      console.log('submit!')
    } else {
      console.log('error submit!')
    }
  })
}

// 定义表单重置
const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return
  formEl.resetFields()
}
</script>

<template>
  <div class="login">
    <div class="box">
        <h2 class="title">登录</h2>
      <el-form
        size="small"
        ref="FormRef"
        style="max-width: 600px"
        :model="FormData"
        :rules="rules"
        label-width="40px"
    >
        <el-form-item label="账号" prop="username">
        <el-input v-model="FormData.username" type="password" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
        <el-input
            v-model="FormData.password"
            type="password"
            autocomplete="off"
        />
        </el-form-item>
        <el-form-item>
        <el-button type="primary" @click="submitForm(FormRef)">
            登录
        </el-button>
        <el-button type="primary" @click="resetForm(FormRef)">
            重置
        </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<style scoped>
.login {
  width: 100vw;
  height: 100vh;
  /* 设置 flex 布局 */
  display: flex;
  /* 设置水平居中 */
  justify-content: center;
  /* 设置垂直居中 */
  align-items: center;
  /* 设置背景颜色 */
  background: linear-gradient(to bottom, #16293f, rgb(128, 167, 211));
}
.box {
  width: 400px;
  /* 设置边框 */
  border: 2px solid #fff;
  /* 设置圆角 */
  border-radius: 10px;
  /* 设置内边距 */
  padding: 20px;
  /* ::deep 修改 element-plus 原始的 css 样式 */
  ::v-deep(.el-form-item__label) {
    color: #000000;
  }
}
.title {
  /* 设置水平居中 */
  text-align: center;
  /* 设置字体大小 */
  font-size: 20px;
}
</style>