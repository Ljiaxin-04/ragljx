<template>
  <div class="login-container">
    <div class="bg-shape blur-one"></div>
    <div class="bg-shape blur-two"></div>
    <div class="bg-shape blur-three"></div>
    <div class="bg-noise"></div>
    <div class="login-shell">
      <div class="brand-panel glass hover-tilt">
        <h1>RAG 知识库系统</h1>
        <p class="lead">智能文档管理 · 检索增强对话 · 多知识库协同</p>
        <div class="feature-chips">
          <span class="chip">向量检索</span>
          <span class="chip">多文件格式</span>
          <span class="chip">流式对话</span>
          <span class="chip">权限与角色</span>
        </div>
      </div>
      <div class="login-box glass hover-tilt">
        <div class="login-header">
          <h2>欢迎回来</h2>
          <p>请登录以继续使用</p>
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
              clearable
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
    radial-gradient(140% 140% at 50% 90%, rgba(255, 255, 255, 0.08), transparent 40%),
    linear-gradient(135deg, #1c2b4a 0%, #2b68ff 45%, #6c5ce7 100%);
}

.bg-noise {
  position: absolute;
  inset: 0;
  pointer-events: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='160' height='160' viewBox='0 0 160 160'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.8' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='160' height='160' filter='url(%23n)' opacity='0.06'/%3E%3C/svg%3E");
  z-index: 0;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.9;
}

.blur-one {
  width: 380px;
  height: 380px;
  background: rgba(255, 255, 255, 0.25);
  top: -80px;
  left: 8%;
}

.blur-two {
  width: 420px;
  height: 420px;
  background: rgba(17, 207, 161, 0.2);
  bottom: -120px;
  right: 5%;
}

.blur-three {
  width: 300px;
  height: 300px;
  background: rgba(255, 255, 255, 0.16);
  top: 30%;
  right: 55%;
}

.login-shell {
  width: 960px;
  max-width: 94vw;
  display: grid;
  grid-template-columns: 1.1fr 0.9fr;
  gap: 18px;
  align-items: stretch;
  position: relative;
  z-index: 1;
}

.glass {
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.18);
  backdrop-filter: blur(16px);
  position: relative;
  overflow: hidden;
}

.glass::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(120deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0));
  opacity: 0.8;
  pointer-events: none;
}

.hover-tilt {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateZ(0);
  transition: transform 0.45s ease, box-shadow 0.45s ease, border-color 0.45s ease;
}

.hover-tilt:hover {
  transform: perspective(1000px) rotateX(3deg) rotateY(-3deg) translateY(-4px);
  box-shadow: 0 30px 80px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 0.32);
}

.brand-panel {
  padding: 42px 36px;
  border-radius: 18px;
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 14px;
}

.brand-panel h1 {
  margin: 0;
  font-size: 30px;
  letter-spacing: 0.6px;
}

.lead {
  margin: 0;
  color: rgba(255, 255, 255, 0.86);
  line-height: 1.6;
}

.feature-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 6px;
}

.chip {
  padding: 8px 12px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.26);
  color: #fff;
  font-size: 13px;
}

.login-box {
  padding: 36px;
  border-radius: 18px;
  position: relative;
}

.login-header {
  text-align: left;
  margin-bottom: 28px;
  color: #fff;
}

.login-header h2 {
  margin: 0 0 6px 0;
  font-size: 24px;
}

.login-header p {
  margin: 0;
  color: rgba(255, 255, 255, 0.78);
  font-size: 14px;
}

.login-form {
  margin-top: 10px;
}

.login-button {
  width: 100%;
  box-shadow: 0 10px 24px rgba(43, 104, 255, 0.35);
}

.login-footer {
  text-align: center;
  margin-top: 16px;
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
  margin: 0 0 6px 0;
  box-shadow: 0 10px 24px rgba(17, 95, 226, 0.35);
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.18);
  border-radius: 10px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.28);
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

@media (max-width: 900px) {
  .login-shell {
    grid-template-columns: 1fr;
  }
  .brand-panel {
    display: none;
  }
  .login-box {
    width: 100%;
  }
}
</style>
