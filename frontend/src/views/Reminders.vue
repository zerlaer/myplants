<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title"><i v-icon="'remind'"></i> 养护提醒</div>
        <div class="page-subtitle">{{ overdueCount }} 项待处理 · {{ todayCount }} 项今日</div>
      </div>
    </div>

    <!-- 类型筛选 -->
    <div class="filter-chips">
      <button class="chip" :class="{ active: type === '' }" @click="setType('')">全部</button>
      <button class="chip water" :class="{ active: type === 'water' }" @click="setType('water')"><i v-icon="'kettle-one'"></i> 浇水</button>
      <button class="chip fertilize" :class="{ active: type === 'fertilize' }" @click="setType('fertilize')"><i v-icon="'pills'"></i> 施肥</button>
      <button class="chip spray" :class="{ active: type === 'spray' }" @click="setType('spray')"><i v-icon="'medicine-bottle-one'"></i> 打药</button>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <template v-else>
      <!-- 批量操作 -->
      <div v-if="overdueReminders.length" class="batch-card">
        <div class="batch-info">
          <i class="batch-icon" v-icon="'flash'"></i>
          <span>{{ overdueReminders.length }} 项逾期，一键全部处理</span>
        </div>
        <button class="btn btn-primary btn-sm" @click="batchCare">一键处理</button>
      </div>

      <div v-if="filtered.length" class="reminder-list">
        <div v-for="r in filtered" :key="r.plant_id + '-' + r.type" class="reminder-card list-item"
             :class="{ overdue: r.overdue }">
          <div class="reminder-avatar" @click="$router.push(`/plants/${r.plant_id}`)">
            <img v-if="r.avatar" :src="r.avatar" />
            <i v-else v-icon="'natural-mode'"></i>
          </div>
          <div class="reminder-main" @click="$router.push(`/plants/${r.plant_id}`)">
            <div class="reminder-name">{{ r.plant_name }}</div>
            <div class="reminder-tag">
              <span class="tag" :class="'tag-' + careTypeMap[r.type].color">
                <i v-icon="careTypeMap[r.type].icon"></i> {{ careTypeMap[r.type].label }}
              </span>
              <span class="reminder-cycle">每 {{ r.cycle_days }} 天</span>
            </div>
            <div class="reminder-time" v-if="r.last_time">上次: {{ timeAgo(r.last_time) }}</div>
            <div class="reminder-time" v-else>还未进行过</div>
          </div>
          <div class="reminder-status">
            <div class="status-badge" :class="statusClass(r)">
              {{ statusText(r) }}
            </div>
            <button class="btn-care btn-sm" :class="'btn-' + careTypeMap[r.type].color"
                    @click="quickCare(r)" :disabled="caring[r.plant_id+'-'+r.type]">
              <i class="icon" v-icon="careTypeMap[r.type].icon"></i>
              <span>立即</span>
            </button>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="emoji" v-icon="'smile'"></i>
        <div class="text">全部养护到位，植物很健康！</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { reminderApi, careApi } from '../api'
import { toast, timeAgo, careTypeMap } from '../utils'

const reminders = ref([])
const loading = ref(true)
const type = ref('')
const caring = ref({})

const setType = (t) => { type.value = t }

const filtered = computed(() => {
  let list = reminders.value
  if (type.value) list = list.filter(r => r.type === type.value)
  // 浇水提醒常驻显示;施肥/打药仅在需要操作时(今日及逾期)显示
  list = list.filter(r => {
    if (r.type === 'water') return true
    // 非浇水类型:仅显示今日及逾期
    return r.overdue || r.days_left <= 0
  })
  // 逾期在前
  return [...list].sort((a, b) => {
    if (a.overdue && !b.overdue) return -1
    if (!a.overdue && b.overdue) return 1
    return a.days_left - b.days_left
  })
})

const overdueCount = computed(() => reminders.value.filter(r => r.overdue).length)
const todayCount = computed(() => reminders.value.filter(r => r.days_left >= 0 && r.days_left <= 1).length)
const overdueReminders = computed(() => reminders.value.filter(r => r.overdue))

const statusClass = (r) => {
  if (r.overdue) return 'badge-danger'
  if (r.days_left <= 1) return 'badge-warning'
  return 'badge-ok'
}

const statusText = (r) => {
  if (r.overdue) return `逾期${Math.abs(r.days_left)}天`
  if (r.days_left === 0) return '今天'
  if (r.days_left === 1) return '明天'
  return `${r.days_left}天后`
}

const loadReminders = async () => {
  loading.value = true
  try {
    const res = await reminderApi.list()
    reminders.value = res.data || []
  } catch (e) {
    toast('加载失败')
  } finally {
    loading.value = false
  }
}

const quickCare = async (r) => {
  const key = r.plant_id + '-' + r.type
  caring.value[key] = true
  try {
    await careApi.oneClick(r.plant_id, r.type)
    toast(`${r.plant_name} ${careTypeMap[r.type].label}完成 ✅`)
    await loadReminders()
  } catch (e) {
    toast(e.message || '操作失败')
  } finally {
    caring.value[key] = false
  }
}

const batchCare = async () => {
  // 收集逾期项,按类型分组批量处理(简化:逐个处理浇水)
  const items = overdueReminders.value
  if (!items.length) return
  if (!confirm(`一键处理 ${items.length} 项逾期提醒？`)) return
  let ok = 0
  for (const r of items) {
    try {
      await careApi.oneClick(r.plant_id, r.type)
      ok++
    } catch (e) {}
  }
  toast(`已处理 ${ok} 项 ✅`)
  loadReminders()
}

onMounted(loadReminders)
</script>

<style scoped>
.filter-chips {
  display: flex;
  gap: 8px;
  overflow-x: auto;
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
.chip.active { background: var(--primary); color: #fff; border-color: var(--primary); }
.chip.water.active { background: var(--water); border-color: var(--water); }
.chip.fertilize.active { background: #5a9216; border-color: #5a9216; }
.chip.spray.active { background: var(--spray); border-color: var(--spray); }

.batch-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: linear-gradient(135deg, #fff4e6, #ffe9d6);
  border-radius: var(--radius);
  padding: 12px 16px;
  margin-bottom: 14px;
  border-left: 4px solid var(--accent);
}
.batch-info { display: flex; align-items: center; gap: 8px; font-size: 13px; font-weight: 600; color: #c46a00; }
.batch-icon { font-size: 18px; }

.reminder-list { display: flex; flex-direction: column; gap: 12px; }
.reminder-card {
  display: flex;
  align-items: center;
  gap: 12px;
  background: var(--card);
  border-radius: var(--radius-lg);
  padding: 14px;
  box-shadow: var(--shadow);
  border-left: 4px solid transparent;
  transition: transform 0.2s;
}
.reminder-card:active { transform: scale(0.98); }
.reminder-card.overdue { border-left-color: var(--danger); }
.reminder-avatar {
  width: 52px;
  height: 52px;
  border-radius: 16px;
  background: var(--primary-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
  font-size: 26px;
  color: var(--primary);
  cursor: pointer;
}
.reminder-avatar img { width: 100%; height: 100%; object-fit: cover; }
.reminder-main { flex: 1; min-width: 0; cursor: pointer; }
.reminder-name { font-weight: 700; font-size: 15px; margin-bottom: 4px; }
.reminder-tag { display: flex; align-items: center; gap: 6px; margin-bottom: 2px; }
.reminder-cycle { font-size: 11px; color: var(--text-light); }
.reminder-time { font-size: 12px; color: var(--text-light); }
.reminder-status { display: flex; flex-direction: column; align-items: flex-end; gap: 8px; }
.status-badge {
  font-size: 11px;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 12px;
}
.badge-danger { background: #fdeaea; color: var(--danger); }
.badge-warning { background: var(--accent-soft); color: var(--accent); }
.badge-ok { background: var(--primary-soft); color: var(--primary); }

.btn-care.btn-sm { padding: 6px 12px; flex-direction: row; gap: 4px; }
.btn-care.btn-sm .icon { font-size: 16px; }
</style>
