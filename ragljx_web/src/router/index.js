import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/views/Register.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/',
      component: () => import('@/layouts/MainLayout.vue'),
      meta: { requiresAuth: true },
      redirect: '/knowledge',
      children: [
        {
          path: 'knowledge',
          name: 'Knowledge',
          component: () => import('@/views/Knowledge.vue'),
          meta: { title: '知识库管理' }
        },
        {
          path: 'knowledge/:id/documents',
          name: 'Documents',
          component: () => import('@/views/Documents.vue'),
          meta: { title: '文档管理' }
        },
        {
          path: 'chat',
          name: 'Chat',
          component: () => import('@/views/Chat.vue'),
          meta: { title: '智能对话' }
        },
        {
          path: 'users',
          name: 'Users',
          component: () => import('@/views/Users.vue'),
          meta: { title: '用户管理', requiresAdmin: true }
        },
        {
          path: 'profile',
          name: 'Profile',
          component: () => import('@/views/Profile.vue'),
          meta: { title: '个人中心' }
        }
      ]
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue')
    }
  ]
})

// 导航守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  const isLoggedIn = userStore.isLoggedIn

  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - RAG 知识库系统`
  } else {
    document.title = 'RAG 知识库系统'
  }

  // 如果路由需要认证
  if (to.meta.requiresAuth !== false) {
    if (!isLoggedIn) {
      // 未登录，重定向到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
    } else {
      // 已登录，检查是否需要管理员权限
      if (to.meta.requiresAdmin && !userStore.isAdmin) {
        // 需要管理员权限但用户不是管理员
        next({ path: '/' })
      } else {
        next()
      }
    }
  } else {
    // 不需要认证的路由
    if (isLoggedIn && (to.path === '/login' || to.path === '/register')) {
      // 已登录用户访问登录/注册页，重定向到首页
      next({ path: '/' })
    } else {
      next()
    }
  }
})

export default router
