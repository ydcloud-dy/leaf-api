import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '@/stores/user'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/',
    component: () => import('@/layouts/AdminLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Index.vue'),
        meta: { title: '控制台' }
      },
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/Index.vue'),
        meta: { title: '用户管理' }
      },
      {
        path: 'articles',
        name: 'Articles',
        component: () => import('@/views/article/Index.vue'),
        meta: { title: '文章管理' }
      },
      {
        path: 'articles/create',
        name: 'CreateArticle',
        component: () => import('@/views/article/Edit.vue'),
        meta: { title: '创建文章' }
      },
      {
        path: 'articles/:id/edit',
        name: 'EditArticle',
        component: () => import('@/views/article/Edit.vue'),
        meta: { title: '编辑文章' }
      },
      {
        path: 'comments',
        name: 'Comments',
        component: () => import('@/views/comment/Index.vue'),
        meta: { title: '评论管理' }
      },
      {
        path: 'tags',
        name: 'Tags',
        component: () => import('@/views/tag/Index.vue'),
        meta: { title: '标签管理' }
      },
      {
        path: 'categories',
        name: 'Categories',
        component: () => import('@/views/category/Index.vue'),
        meta: { title: '分类管理' }
      },
      {
        path: 'chapters',
        name: 'Chapters',
        component: () => import('@/views/chapter/Index.vue'),
        meta: { title: '章节管理' }
      },
      {
        path: 'files',
        name: 'Files',
        component: () => import('@/views/file/Index.vue'),
        meta: { title: '文件管理' }
      },
      {
        path: 'settings',
        name: 'Settings',
        component: () => import('@/views/settings/Index.vue'),
        meta: { title: '系统设置' }
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/Index.vue'),
        meta: { title: '个人信息' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()

  if (to.meta.requiresAuth !== false && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.path === '/login' && userStore.isLoggedIn) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
