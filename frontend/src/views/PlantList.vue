<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title"><i v-icon="'natural-mode'"></i> 我的植物</div>
        <div class="page-subtitle" v-if="plants.length">共 {{ plants.length }} 棵植物</div>
      </div>
      <router-link to="/plants/new" class="btn btn-primary btn-sm"><i v-icon="'plus'"></i> 添加</router-link>
    </div>

    <!-- 搜索 -->
    <div class="search-bar">
      <i class="search-icon" v-icon="'search'"></i>
      <input v-model="keyword" placeholder="搜索植物名称或品种..." class="search-input" @input="onSearch" />
    </div>

    <!-- 分类筛选 -->
    <div class="filter-chips">
      <button class="chip" :class="{ active: category === '' }" @click="setCategory('')">全部</button>
      <button v-for="cat in categories" :key="cat" class="chip" :class="{ active: category === cat }" @click="setCategory(cat)">
        <i v-icon="categoryMap[cat]"></i> {{ cat }}
      </button>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <div v-else-if="filtered.length" class="plant-list">
      <div v-for="(plant, i) in filtered" :key="plant.id" class="plant-card list-item" :style="{ animationDelay: i * 0.05 + 's' }" @click="$router.push(`/plants/${plant.id}`)">
        <div class="plant-avatar">
          <img v-if="plant.avatar" :src="plant.avatar" :alt="plant.name" />
          <i v-else v-icon="'natural-mode'"></i>
        </div>
        <div class="plant-info">
          <div class="plant-name-row">
            <span class="plant-name">{{ plant.name }}</span>
            <span class="tag" :class="healthTagClass(plant.health_status)">{{ plant.health_status }}</span>
          </div>
          <div class="plant-meta">
            <span v-if="plant.species">{{ plant.species }}</span>
            <span v-if="plant.location"><i v-icon="'local'"></i> {{ plant.location }}</span>
            <span class="care-days" v-if="careDays(plant) !== null">
              <i v-icon="'calendar'"></i> 养护 {{ careDays(plant) }} 天
            </span>
            <span class="plant-price" v-if="plant.price > 0">
              <i v-icon="'tag'"></i> ¥{{ formatPrice(plant.price) }}
            </span>
          </div>
          <div class="plant-care-info">
            <span class="care-tag water" v-if="plant.last_watered_at">
              <i v-icon="'kettle-one'"></i> {{ timeAgo(plant.last_watered_at) }}
            </span>
            <span class="care-tag water overdue" v-else><i v-icon="'kettle-one'"></i> 待浇水</span>
          </div>
        </div>
        <div class="plant-actions" @click.stop>
          <button class="mini-btn water" @click="quickCare(plant, 'water')" :disabled="caring[plant.id+'-water']"><i v-icon="'kettle-one'"></i></button>
          <button class="mini-btn fertilize" @click="quickCare(plant, 'fertilize')" :disabled="caring[plant.id+'-fertilize']"><i v-icon="'pills'"></i></button>
          <button class="mini-btn spray" @click="quickCare(plant, 'spray')" :disabled="caring[plant.id+'-spray']"><i v-icon="'medicine-bottle-one'"></i></button>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <i class="emoji" v-icon="'natural-mode'"></i>
      <div class="text">还没有植物，快添加第一棵吧～</div>
      <router-link to="/plants/new" class="btn btn-primary"><i v-icon="'plus'"></i> 添加植物</router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { plantApi, careApi } from '../api'
import { toast, timeAgo, categoryMap } from '../utils'

const plants = ref([])
const loading = ref(true)
const keyword = ref('')
const category = ref('')
const caring = ref({})
const categories = ['绿植', '多肉', '花卉', '草本', '木本', '果树']

const filtered = computed(() => {
  if (!keyword.value) return plants.value
  const k = keyword.value.toLowerCase()
  return plants.value.filter(p =>
    p.name.toLowerCase().includes(k) ||
    (p.species || '').toLowerCase().includes(k)
  )
})

const setCategory = (cat) => { category.value = cat; loadPlants() }

const onSearch = () => {}

const healthTagClass = (status) => {
  return {
    '优秀': 'tag-primary',
    '良好': 'tag-primary',
    '一般': 'tag-warning',
    '需关注': 'tag-danger'
  }[status] || 'tag-grey'
}

const careDays = (plant) => {
  const base = plant.acquired_at || plant.created_at
  if (!base) return null
  const d = new Date(base)
  if (isNaN(d.getTime())) return null
  const diff = Date.now() - d.getTime()
  return Math.max(0, Math.floor(diff / (1000 * 60 * 60 * 24)))
}

const formatPrice = (v) => {
  const n = Number(v) || 0
  return n.toFixed(2)
}

const loadPlants = async () => {
  loading.value = true
  try {
    const params = {}
    if (category.value) params.category = category.value
    const res = await plantApi.list(params)
    plants.value = res.data || []
  } catch (e) {
    toast('加载失败')
  } finally {
    loading.value = false
  }
}

const quickCare = async (plant, type) => {
  const key = plant.id + '-' + type
  caring.value[key] = true
  try {
    await careApi.oneClick(plant.id, type)
    const labels = { water: '浇水', fertilize: '施肥', spray: '打药' }
    toast(`「${plant.name}」${labels[type]}完成 ✅`)
    loadPlants()
  } catch (e) {
    toast(e.message || '操作失败')
  } finally {
    caring.value[key] = false
  }
}

onMounted(loadPlants)
</script>

<style scoped>
.search-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  background: var(--card);
  border-radius: var(--radius);
  padding: 10px 14px;
  margin-bottom: 14px;
  box-shadow: var(--shadow-sm);
}
.search-icon { font-size: 16px; color: var(--text-light); display: inline-flex; }
.search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  font-size: 15px;
}

.filter-chips {
  display: flex;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  margin-bottom: 16px;
  scrollbar-width: none;
}
.filter-chips::-webkit-scrollbar { display: none; }
.chip {
  flex-shrink: 0;
  padding: 7px 16px;
  border-radius: 20px;
  background: var(--card);
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 600;
  border: 1.5px solid var(--border);
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
.chip i { font-size: 14px; }
.chip.active {
  background: var(--primary);
  color: #fff;
  border-color: var(--primary);
}

.plant-list { display: flex; flex-direction: column; gap: 12px; }

.plant-card {
  display: flex;
  align-items: center;
  gap: 14px;
  background: var(--card);
  border-radius: var(--radius-lg);
  padding: 14px;
  box-shadow: var(--shadow);
  cursor: pointer;
  transition: transform 0.2s;
}
.plant-card:active { transform: scale(0.98); }

.plant-avatar {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  overflow: hidden;
  flex-shrink: 0;
  background: var(--primary-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--primary);
}
.plant-avatar img { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { font-size: 32px; }
.plant-avatar i { font-size: 32px; color: var(--primary); }

.plant-info { flex: 1; min-width: 0; }
.plant-name-row { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.plant-name { font-size: 16px; font-weight: 700; color: var(--text); }
.plant-meta {
  display: flex;
  gap: 10px;
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
  flex-wrap: wrap;
  align-items: center;
}
.plant-meta i { font-size: 13px; }
.care-days {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  color: var(--primary);
  font-weight: 600;
}
.plant-price {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  color: var(--accent);
  font-weight: 700;
}
.plant-care-info { display: flex; gap: 8px; flex-wrap: wrap; }
.care-tag {
  font-size: 11px;
  padding: 2px 8px;
  border-radius: 10px;
  background: var(--water-soft);
  color: var(--water);
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 3px;
}
.care-tag i { font-size: 12px; }
.care-tag.overdue { background: #fdeaea; color: var(--danger); }

.plant-actions {
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.mini-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  transition: transform 0.15s;
}
.mini-btn:active:not(:disabled) { transform: scale(0.85); }
.mini-btn:disabled { opacity: 0.5; }
.mini-btn.water { background: var(--water-soft); }
.mini-btn.fertilize { background: var(--fertilize-soft); }
.mini-btn.spray { background: var(--spray-soft); }
</style>
