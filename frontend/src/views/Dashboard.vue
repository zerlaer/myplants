<template>
  <div class="page">
    <!-- 头部 -->
    <div class="hero-header">
      <div class="hero-bg"></div>
      <div class="hero-content">
        <div class="hero-greeting">
          <h1><i v-icon="'natural-mode'"></i> 我的花园</h1>
          <p>{{ greeting }}，愿你与绿植一同生长</p>
        </div>
        <div class="hero-weather" v-if="data">
          <span class="big">{{ data.plant_count }}</span>
          <span class="label">棵植物</span>
        </div>
      </div>
    </div>

    <!-- 待办提醒横幅 -->
    <div v-if="data && data.overdue_count > 0" class="alert-banner" @click="$router.push('/reminders')">
      <i class="icon" v-icon="'caution'"></i>
      <div class="alert-text">
        <div class="alert-title">{{ data.overdue_count }} 棵植物需要养护</div>
        <div class="alert-sub">点击查看详情</div>
      </div>
      <i class="arrow" v-icon="'right'"></i>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <template v-else-if="data">
      <!-- 总价值卡片 -->
      <div class="total-value-card">
        <div class="tv-left">
          <div class="tv-icon"><i v-icon="'tag'"></i></div>
          <div class="tv-info">
            <div class="tv-label">植物总价值</div>
            <div class="tv-amount">¥{{ formatPrice(data.total_price) }}</div>
          </div>
        </div>
        <div class="tv-right">
          <div class="tv-num">{{ data.plant_count }}</div>
          <div class="tv-unit">棵植物</div>
        </div>
      </div>

      <!-- 统计卡片 -->
      <div class="stats-grid">
        <div class="stat-card stat-water">
          <i class="stat-icon" v-icon="'kettle-one'"></i>
          <div class="stat-info">
            <div class="stat-num">{{ data.water_count }}</div>
            <div class="stat-label">浇水次数</div>
          </div>
          <div class="stat-today">今日 {{ data.today_water }}</div>
        </div>
        <div class="stat-card stat-fertilize">
          <i class="stat-icon" v-icon="'pills'"></i>
          <div class="stat-info">
            <div class="stat-num">{{ data.fertilize_count }}</div>
            <div class="stat-label">施肥次数</div>
          </div>
          <div class="stat-today">今日 {{ data.today_fertilize }}</div>
        </div>
        <div class="stat-card stat-spray">
          <i class="stat-icon" v-icon="'medicine-bottle-one'"></i>
          <div class="stat-info">
            <div class="stat-num">{{ data.spray_count }}</div>
            <div class="stat-label">打药次数</div>
          </div>
          <div class="stat-today">今日 {{ data.today_spray }}</div>
        </div>
        <div class="stat-card stat-photo">
          <i class="stat-icon" v-icon="'picture-album'"></i>
          <div class="stat-info">
            <div class="stat-num">{{ data.photo_count }}</div>
            <div class="stat-label">成长照片</div>
          </div>
          <div class="stat-today">{{ data.pot_count }} 花盆</div>
        </div>
      </div>

      <!-- 近7天养护趋势 -->
      <div class="section-title"><i v-icon="'chart-line'"></i> 近7天养护趋势</div>
      <div class="card trend-card">
        <div class="trend-chart">
          <div v-for="item in data.trend" :key="item.date" class="trend-bar-wrap">
            <div class="trend-bar" :style="{ height: barHeight(item.count) + 'px' }">
              <span class="trend-num" v-if="item.count > 0">{{ item.count }}</span>
            </div>
            <div class="trend-label">{{ formatMonthDay(item.date) }}</div>
          </div>
        </div>
      </div>

      <!-- 分类统计(数量+价格) -->
      <div class="section-title"><i v-icon="'grid-four'"></i> 分类统计</div>
      <div class="card" v-if="data.category_price_stats && data.category_price_stats.length">
        <div class="cat-price-list">
          <div v-for="item in data.category_price_stats" :key="item.category" class="cat-price-item">
            <i class="cp-emoji" v-icon="categoryMap[item.category] || 'natural-mode'"></i>
            <span class="cp-name">{{ item.category || '未分类' }}</span>
            <span class="cp-count">{{ item.count }} 棵</span>
            <span class="cp-price">¥{{ formatPrice(item.total_price) }}</span>
          </div>
        </div>
      </div>
      <div class="card" v-else>
        <div class="empty-hint">暂无植物，快去添加第一棵吧～</div>
      </div>

      <!-- 我的植物 -->
      <div class="section-title">
        <i v-icon="'natural-mode'"></i> 我的植物
        <router-link to="/plants" class="section-more">查看全部 <i v-icon="'right'"></i></router-link>
      </div>
      <div class="plant-name-grid" v-if="data.plants && data.plants.length">
        <router-link v-for="p in data.plants" :key="p.id" :to="`/plants/${p.id}`" class="plant-name-card">
          <div class="pnc-avatar">
            <img v-if="p.avatar" :src="p.avatar" />
            <i v-else v-icon="'natural-mode'"></i>
          </div>
          <div class="pnc-name">{{ p.name }}</div>
          <div class="pnc-price" v-if="p.price > 0">¥{{ formatPrice(p.price) }}</div>
        </router-link>
      </div>
      <div class="card" v-else>
        <div class="empty-hint">还没有植物，快去添加吧～</div>
      </div>

      <!-- 植物状态 -->
      <div class="section-title"><i v-icon="'like'"></i> 植物状态</div>
      <div class="card" v-if="data.status_stats && data.status_stats.length">
        <div class="status-bar">
          <div v-for="item in data.status_stats" :key="item.status"
               class="status-seg"
               :class="statusClass(item.status)"
               :style="{ flex: item.count }">
            {{ item.status }} {{ item.count }}
          </div>
        </div>
      </div>

      <!-- 快捷入口 -->
      <div class="section-title"><i v-icon="'flash'"></i> 快捷入口</div>
      <div class="quick-grid">
        <router-link to="/plants" class="quick-item">
          <i class="emoji" v-icon="'natural-mode'"></i><span>我的植物</span>
        </router-link>
        <router-link to="/reminders" class="quick-item">
          <i class="emoji" v-icon="'remind'"></i><span>养护提醒</span>
        </router-link>
        <router-link to="/calendar" class="quick-item">
          <i class="emoji" v-icon="'calendar'"></i><span>养护日历</span>
        </router-link>
        <router-link to="/pots" class="quick-item">
          <i class="emoji" v-icon="'box'"></i><span>花盆管理</span>
        </router-link>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { dashboardApi } from '../api'
import { categoryMap } from '../utils'

const data = ref(null)
const loading = ref(true)

const greeting = computed(() => {
  const h = new Date().getHours()
  if (h < 6) return '夜深了'
  if (h < 11) return '早上好'
  if (h < 14) return '中午好'
  if (h < 18) return '下午好'
  return '晚上好'
})

const barHeight = (count) => {
  if (!data.value || !data.value.trend) return 4
  const max = Math.max(...data.value.trend.map(t => t.count), 1)
  return Math.max(4, (count / max) * 80)
}

const statusClass = (status) => {
  return {
    '旺盛': 'seg-excellent',
    '健康': 'seg-good',
    '欠佳': 'seg-normal',
    '萎蔫': 'seg-warning'
  }[status] || 'seg-good'
}

const formatPrice = (v) => {
  const n = Number(v) || 0
  return n.toFixed(2)
}

const formatMonthDay = (dateStr) => {
  if (!dateStr) return ''
  const parts = String(dateStr).split('-')
  if (parts.length >= 3) {
    return `${parseInt(parts[1])}月${parseInt(parts[2])}日`
  }
  return dateStr
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await dashboardApi.get()
    data.value = res.data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.hero-header {
  position: relative;
  border-radius: var(--radius-xl);
  overflow: hidden;
  margin-bottom: 16px;
  color: #fff;
}
.hero-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #2d8659 0%, #4ba878 50%, #7fc89e 100%);
}
.hero-bg::before {
  content: '';
  position: absolute;
  right: -30px;
  top: -30px;
  width: 160px;
  height: 160px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.12);
}
.hero-bg::after {
  content: '';
  position: absolute;
  right: 40px;
  bottom: -40px;
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}
.hero-content {
  position: relative;
  padding: 24px 22px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.hero-greeting h1 { font-size: 26px; font-weight: 800; margin-bottom: 6px; }
.hero-greeting p { font-size: 13px; opacity: 0.92; }
.hero-weather { text-align: center; }
.hero-weather .big { font-size: 38px; font-weight: 800; display: block; }
.hero-weather .label { font-size: 12px; opacity: 0.9; }

.alert-banner {
  display: flex;
  align-items: center;
  gap: 12px;
  background: linear-gradient(135deg, #fff4e6, #ffe9d6);
  border-radius: var(--radius);
  padding: 14px 18px;
  margin-bottom: 16px;
  border-left: 4px solid var(--accent);
}
.alert-banner .icon { font-size: 28px; color: #c46a00; }
.alert-banner .arrow { font-size: 24px; color: #c46a00; }
.alert-text { flex: 1; }
.alert-title { font-weight: 700; color: #c46a00; }
.alert-sub { font-size: 12px; color: #b88555; }

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
  margin-bottom: 8px;
}
.stat-card {
  background: var(--card);
  border-radius: var(--radius);
  padding: 14px;
  display: flex;
  align-items: center;
  gap: 10px;
  position: relative;
  box-shadow: var(--shadow-sm);
  overflow: hidden;
}
.stat-card::before {
  content: '';
  position: absolute;
  left: 0; top: 0; bottom: 0;
  width: 4px;
}
.stat-water::before { background: var(--water); }
.stat-fertilize::before { background: var(--fertilize); }
.stat-spray::before { background: var(--spray); }
.stat-photo::before { background: var(--accent); }
.stat-icon {
  font-size: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
}
.stat-info { flex: 1; }
.stat-num { font-size: 22px; font-weight: 800; color: var(--text); }
.stat-label { font-size: 12px; color: var(--text-secondary); }
.stat-today {
  position: absolute;
  top: 8px;
  right: 10px;
  font-size: 11px;
  color: var(--text-light);
  background: var(--bg);
  padding: 2px 8px;
  border-radius: 10px;
}

.trend-card { padding: 20px 14px 14px; }
.trend-chart {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  height: 110px;
  gap: 6px;
}
.trend-bar-wrap { flex: 1; display: flex; flex-direction: column; align-items: center; gap: 6px; }
.trend-bar {
  width: 100%;
  max-width: 32px;
  background: linear-gradient(180deg, var(--primary-light), var(--primary));
  border-radius: 6px 6px 0 0;
  min-height: 4px;
  position: relative;
  transition: height 0.4s ease;
  display: flex;
  justify-content: center;
}
.trend-num {
  position: absolute;
  top: -18px;
  font-size: 11px;
  font-weight: 700;
  color: var(--primary);
}
.trend-label { font-size: 10px; color: var(--text-light); }

.category-list { display: flex; flex-direction: column; gap: 10px; }
.category-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  background: var(--bg);
  border-radius: var(--radius-sm);
}
.cat-emoji { font-size: 22px; }
.cat-name { flex: 1; font-weight: 600; }
.cat-count {
  background: var(--primary);
  color: #fff;
  padding: 2px 12px;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 700;
}

.empty-hint { text-align: center; color: var(--text-light); padding: 12px; font-size: 13px; }

/* 总价值卡片 */
.total-value-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #2d8659 0%, #4ba878 100%);
  border-radius: var(--radius-lg);
  padding: 18px 20px;
  margin-bottom: 14px;
  color: #fff;
  box-shadow: 0 6px 20px rgba(45, 134, 89, 0.25);
}
.tv-left { display: flex; align-items: center; gap: 14px; }
.tv-icon {
  width: 48px; height: 48px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.2);
  display: flex; align-items: center; justify-content: center;
  font-size: 26px;
}
.tv-label { font-size: 13px; opacity: 0.9; }
.tv-amount { font-size: 26px; font-weight: 800; }
.tv-right { text-align: center; }
.tv-num { font-size: 32px; font-weight: 800; }
.tv-unit { font-size: 12px; opacity: 0.9; }

/* 分类统计列表 */
.cat-price-list { display: flex; flex-direction: column; gap: 10px; }
.cat-price-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: var(--bg);
  border-radius: var(--radius-sm);
}
.cp-emoji { font-size: 22px; }
.cp-name { flex: 1; font-weight: 600; }
.cp-count { font-size: 13px; color: var(--text-secondary); }
.cp-price {
  font-size: 15px;
  font-weight: 700;
  color: var(--primary);
}

/* 植物名称小卡片 */
.section-more {
  margin-left: auto;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 2px;
}
.section-more i { font-size: 16px; }
.plant-name-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 10px;
}
.plant-name-card {
  background: var(--card);
  border-radius: var(--radius);
  padding: 12px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  box-shadow: var(--shadow-sm);
  transition: transform 0.2s;
}
.plant-name-card:active { transform: scale(0.95); }
.pnc-avatar {
  width: 48px; height: 48px;
  border-radius: 50%;
  background: var(--primary-soft);
  display: flex; align-items: center; justify-content: center;
  overflow: hidden;
  font-size: 24px;
  color: var(--primary);
}
.pnc-avatar img { width: 100%; height: 100%; object-fit: cover; }
.pnc-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
  text-align: center;
  max-width: 100%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.pnc-price {
  font-size: 12px;
  color: var(--primary);
  font-weight: 700;
}

.status-bar {
  display: flex;
  height: 36px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  gap: 2px;
}
.status-seg {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 700;
  color: #fff;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
}
.seg-excellent { background: var(--success); }
.seg-good { background: var(--primary-light); }
.seg-normal { background: var(--accent); }
.seg-warning { background: var(--danger); }

.quick-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
}
.quick-item {
  background: var(--card);
  border-radius: var(--radius);
  padding: 16px 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  box-shadow: var(--shadow-sm);
  color: var(--text);
  font-size: 12px;
  font-weight: 600;
  transition: transform 0.2s;
}
.quick-item:active { transform: scale(0.95); }
.quick-item .emoji { font-size: 28px; display: inline-flex; }
</style>
