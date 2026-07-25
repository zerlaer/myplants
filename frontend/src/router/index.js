import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  { path: '/', redirect: '/dashboard' },
  { path: '/dashboard', name: 'Dashboard', component: () => import('../views/Dashboard.vue'), meta: { title: '首页' } },
  { path: '/plants', name: 'PlantList', component: () => import('../views/PlantList.vue'), meta: { title: '我的植物' } },
  { path: '/plants/new', name: 'PlantNew', component: () => import('../views/PlantForm.vue'), meta: { title: '添加植物' } },
  { path: '/plants/:id/edit', name: 'PlantEdit', component: () => import('../views/PlantForm.vue'), meta: { title: '编辑植物' } },
  { path: '/plants/:id', name: 'PlantDetail', component: () => import('../views/PlantDetail.vue'), meta: { title: '植物详情' } },
  { path: '/reminders', name: 'Reminders', component: () => import('../views/Reminders.vue'), meta: { title: '养护提醒' } },
  { path: '/pots', name: 'Pots', component: () => import('../views/Pots.vue'), meta: { title: '花盆管理' } },
  { path: '/calendar', name: 'Calendar', component: () => import('../views/Calendar.vue'), meta: { title: '养护日历' } },
  { path: '/plants/:id/album', name: 'Album', component: () => import('../views/Album.vue'), meta: { title: '成长相册' } }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  }
})

router.afterEach((to) => {
  if (to.meta.title) {
    document.title = to.meta.title + ' - 我的花园'
  }
})

export default router
