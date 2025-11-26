<template>
  <div class="login-container">
    <div class="bg-shape blur-one"></div>
    <div class="bg-shape blur-two"></div>
    <div class="login-box glass">
      <div class="login-header">
        <div class="logo-pill">
          <span>R</span>
        </div>
        <h1>RAG 知识库系统</h1>
        <p>智能文档管理与对话平台</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
        
        <div class="login-footer">
          <span>还没有账号？</span>
          <el-link type="primary" @click="goToRegister">立即注册</el-link>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度在 3 到 50 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await userStore.login(loginForm)
        ElMessage.success('登录成功')
        
        // 跳转到之前的页面或首页
        const redirect = route.query.redirect || '/'
        router.push(redirect)
      } catch (error) {
        console.error('Login failed:', error)
        ElMessage.error(error.response?.data?.message || '登录失败，请检查用户名和密码')
      } finally {
        loading.value = false
      }
    }
  })
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  background: radial-gradient(120% 120% at 20% 20%, rgba(91, 140, 255, 0.25), transparent 40%),
    radial-gradient(120% 120% at 80% 0%, rgba(17, 207, 161, 0.25), transparent 35%),
    linear-gradient(135deg, #2b68ff 0%, #6c5ce7 50%, #5f8bff 100%);
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.8;
}

.blur-one {
  width: 380px;
  height: 380px;
  background: rgba(255, 255, 255, 0.25);
  top: -80px;
  left: 10%;
}

.blur-two {
  width: 420px;
  height: 420px;
  background: rgba(17, 207, 161, 0.2);
  bottom: -120px;
  right: 5%;
}

.login-box {
  width: 460px;
  padding: 42px;
  border-radius: 18px;
  position: relative;
  z-index: 1;
}

.glass {
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.18);
  backdrop-filter: blur(14px);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 28px;
  color: #fff;
  margin: 8px 0 6px 0;
}

.login-header p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.86);
  margin: 0;
}

.login-form {
  margin-top: 20px;
}

.login-button {
  width: 100%;
  box-shadow: 0 10px 24px rgba(43, 104, 255, 0.35);
}

.login-footer {
  text-align: center;
  margin-top: 20px;
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
}

.login-footer span {
  margin-right: 8px;
}

.logo-pill {
  width: 54px;
  height: 54px;
  border-radius: 16px;
  background: linear-gradient(135deg, rgba(43, 104, 255, 0.65), rgba(17, 207, 161, 0.7));
  display: grid;
  place-items: center;
  color: #fff;
  font-weight: 800;
  font-size: 24px;
  margin: 0 auto 6px auto;
  box-shadow: 0 10px 24px rgba(17, 95, 226, 0.35);
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.18);
  border-radius: 10px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

:deep(.el-input__inner) {
  color: #fff;
}

:deep(.el-input__prefix i) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-form-item__error) {
  color: #ffe1e1;
}

:deep(.el-link) {
  color: #cfe2ff;
}
</style>
