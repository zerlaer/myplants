<template>
  <router-view v-slot="{ Component, route }">
    <component :is="Component" :key="route.path" />
  </router-view>
  <BottomNav v-if="showNav" />
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import BottomNav from './components/BottomNav.vue'

const route = useRoute()
// 详情/编辑页隐藏底部导航
const showNav = computed(() => {
  const hidePatterns = ['/plants/', '/new', '/edit', '/album']
  return !hidePatterns.some(p => route.path.includes(p))
})
</script>
