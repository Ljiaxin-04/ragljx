<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <div class="card-header">
          <h2>个人中心</h2>
        </div>
      </template>
      
      <div class="profile-content">
        <div class="avatar-section">
          <el-avatar :size="100" :src="userStore.avatar">
            {{ userStore.nickname.charAt(0) }}
          </el-avatar>
          <h3>{{ userStore.nickname }}</h3>
          <p>@{{ userStore.username }}</p>
        </div>
        
        <el-divider />
        
        <el-form
          ref="formRef"
          :model="form"
          :rules="formRules"
          label-width="100px"
          class="profile-form"
        >
          <el-form-item label="用户名">
            <el-input v-model="userStore.username" disabled />
          </el-form-item>
          
          <el-form-item label="昵称" prop="nickname">
            <el-input v-model="form.nickname" placeholder="请输入昵称" />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" :loading="updating" @click="updateProfile">
              保存修改
            </el-button>
          </el-form-item>
        </el-form>
        
        <el-divider />
        
        <h3>修改密码</h3>
        <el-form
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          label-width="100px"
          class="password-form"
        >
          <el-form-item label="原密码" prop="old_password">
            <el-input
              v-model="passwordForm.old_password"
              type="password"
              placeholder="请输入原密码"
              show-password
            />
          </el-form-item>
          
          <el-form-item label="新密码" prop="new_password">
            <el-input
              v-model="passwordForm.new_password"
              type="password"
              placeholder="请输入新密码"
              show-password
            />
          </el-form-item>
          
          <el-form-item label="确认密码" prop="confirm_password">
            <el-input
              v-model="passwordForm.confirm_password"
              type="password"
              placeholder="请确认新密码"
              show-password
            />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" :loading="updatingPassword" @click="updatePassword">
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { updateUser, updatePassword as updatePasswordApi } from '@/api/user'

const userStore = useUserStore()

const formRef = ref(null)
const passwordFormRef = ref(null)
const updating = ref(false)
const updatingPassword = ref(false)

const form = reactive({
  nickname: '',
  email: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const formRules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== passwordForm.new_password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const passwordRules = {
  old_password: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 50, message: '密码长度在 6 到 50 个字符', trigger: 'blur' }
  ],
  confirm_password: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const updateProfile = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      updating.value = true
      try {
        await updateUser(userStore.userInfo.id, form)
        await userStore.getUserInfo()
        ElMessage.success('更新成功')
      } catch (error) {
        console.error('Update profile failed:', error)
        ElMessage.error(error.response?.data?.message || '更新失败')
      } finally {
        updating.value = false
      }
    }
  })
}

const updatePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      updatingPassword.value = true
      try {
        await updatePasswordApi({
          old_password: passwordForm.old_password,
          new_password: passwordForm.new_password
        })
        ElMessage.success('密码修改成功，请重新登录')
        
        // 清空表单
        passwordForm.old_password = ''
        passwordForm.new_password = ''
        passwordForm.confirm_password = ''
        
        // 登出
        setTimeout(() => {
          userStore.logout()
        }, 1500)
      } catch (error) {
        console.error('Update password failed:', error)
        ElMessage.error(error.response?.data?.message || '修改密码失败')
      } finally {
        updatingPassword.value = false
      }
    }
  })
}

onMounted(() => {
  if (userStore.userInfo) {
    form.nickname = userStore.userInfo.nickname || userStore.userInfo.username
    form.email = userStore.userInfo.email || ''
  }
})
</script>

<style scoped>
.profile-container {
  max-width: 800px;
  margin: 0 auto;
}

.card-header h2 {
  margin: 0;
  font-size: 20px;
}

.profile-content {
  padding: 20px;
}

.avatar-section {
  text-align: center;
  margin-bottom: 30px;
}

.avatar-section h3 {
  margin: 15px 0 5px 0;
  font-size: 20px;
  color: #333;
}

.avatar-section p {
  margin: 0;
  color: #999;
  font-size: 14px;
}

.profile-form,
.password-form {
  max-width: 500px;
  margin: 20px auto;
}
</style>

