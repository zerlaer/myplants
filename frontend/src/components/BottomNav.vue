<template>
  <nav class="bottom-nav">
    <router-link to="/dashboard" class="bottom-nav-item" :class="{ active: isActive('/dashboard') }">
      <i class="icon" v-icon="'home'"></i>
      <span>首页</span>
    </router-link>
    <router-link to="/plants" class="bottom-nav-item" :class="{ active: isActive('/plants') }">
      <i class="icon" v-icon="'natural-mode'"></i>
      <span>植物</span>
    </router-link>
    <router-link to="/plants/new" class="bottom-nav-item center">
      <i class="icon" v-icon="'plus'"></i>
      <span>添加</span>
    </router-link>
    <router-link to="/reminders" class="bottom-nav-item" :class="{ active: isActive('/reminders') }">
      <i class="icon" v-icon="'remind'"></i>
      <span>提醒</span>
    </router-link>
    <router-link to="/pots" class="bottom-nav-item" :class="{ active: isActive('/pots') }">
      <i class="icon" v-icon="'box'"></i>
      <span>花盆</span>
    </router-link>
  </nav>
</template>

<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()
const isActive = (path) => route.path === path || route.path.startsWith(path + '/')
</script>
