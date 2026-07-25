<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title"><i v-icon="'calendar'"></i> 养护日历</div>
        <div class="page-subtitle">{{ monthLabel }}</div>
      </div>
      <div class="month-nav">
        <button @click="prevMonth"><i v-icon="'left'"></i></button>
        <button @click="goToday" class="today-btn">今天</button>
        <button @click="nextMonth"><i v-icon="'right'"></i></button>
      </div>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <template v-else>
      <!-- 日历主体 -->
      <div class="calendar">
        <div class="cal-weekday" v-for="w in weekdays" :key="w">{{ w }}</div>
        <div v-for="i in firstDayOffset" :key="'blank'+i" class="cal-day blank"></div>
        <div v-for="day in daysInMonth" :key="day"
             class="cal-day"
             :class="{ today: isToday(day), hasRecords: hasRecords(day), selected: selectedDay === day }"
             @click="selectDay(day)">
          <span class="day-num">{{ day }}</span>
          <div class="day-dots" v-if="recordsByDay[day]">
            <span v-if="recordsByDay[day].water" class="dot water"></span>
            <span v-if="recordsByDay[day].fertilize" class="dot fertilize"></span>
            <span v-if="recordsByDay[day].spray" class="dot spray"></span>
          </div>
        </div>
      </div>

      <!-- 图例 -->
      <div class="legend">
        <span class="legend-item"><span class="dot water"></span>浇水</span>
        <span class="legend-item"><span class="dot fertilize"></span>施肥</span>
        <span class="legend-item"><span class="dot spray"></span>打药</span>
      </div>

      <!-- 选中日的记录 -->
      <div v-if="selectedRecords.length" class="section-title">
        {{ formatDate(selectedDate) }} 的记录 ({{ selectedRecords.length }})
      </div>
      <div v-if="selectedRecords.length" class="card" style="padding:14px">
        <div class="timeline">
          <div v-for="r in selectedRecords" :key="r.id" class="timeline-item"
               @click="goPlant(r.plant_id)">
            <div class="timeline-dot" :class="careTypeMap[r.type].color"><i v-icon="careTypeMap[r.type].icon"></i></div>
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="plant-link">{{ getPlantName(r.plant_id) }}</span>
                <span class="tag" :class="'tag-' + careTypeMap[r.type].color">{{ careTypeMap[r.type].label }}</span>
                <span class="timeline-time">{{ formatTime(r.record_time) }}</span>
              </div>
              <div class="timeline-remark" v-if="r.remark">{{ r.remark }}</div>
            </div>
          </div>
        </div>
      </div>

      <div v-else-if="selectedDay" class="empty-state">
        <i class="emoji" v-icon="'cloud'"></i>
        <div class="text">这天没有养护记录</div>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { careApi, plantApi } from '../api'
import { toast, formatDate, careTypeMap } from '../utils'

const router = useRouter()
const loading = ref(true)
const records = ref([])
const plants = ref([])
const current = ref(new Date())
const selectedDay = ref(null)

const weekdays = ['日', '一', '二', '三', '四', '五', '六']

const monthLabel = computed(() => {
  const y = current.value.getFullYear()
  const m = current.value.getMonth() + 1
  return `${y}年${m}月`
})

const daysInMonth = computed(() => new Date(current.value.getFullYear(), current.value.getMonth() + 1, 0).getDate())
const firstDayOffset = computed(() => new Date(current.value.getFullYear(), current.value.getMonth(), 1).getDay())

const selectedDate = computed(() => {
  if (!selectedDay.value) return null
  return new Date(current.value.getFullYear(), current.value.getMonth(), selectedDay.value)
})

const recordsByDay = computed(() => {
  const map = {}
  const y = current.value.getFullYear()
  const m = current.value.getMonth()
  records.value.forEach(r => {
    const d = new Date(r.record_time)
    if (d.getFullYear() === y && d.getMonth() === m) {
      const day = d.getDate()
      if (!map[day]) map[day] = {}
      map[day][r.type] = true
    }
  })
  return map
})

const selectedRecords = computed(() => {
  if (!selectedDay.value) return []
  const y = current.value.getFullYear()
  const m = current.value.getMonth()
  return records.value
    .filter(r => {
      const d = new Date(r.record_time)
      return d.getFullYear() === y && d.getMonth() === m && d.getDate() === selectedDay.value
    })
    .sort((a, b) => new Date(a.record_time) - new Date(b.record_time))
})

const hasRecords = (day) => !!recordsByDay.value[day]
const isToday = (day) => {
  const now = new Date()
  return current.value.getFullYear() === now.getFullYear() &&
    current.value.getMonth() === now.getMonth() && day === now.getDate()
}

const selectDay = (day) => { selectedDay.value = selectedDay.value === day ? null : day }

const prevMonth = () => { current.value = new Date(current.value.getFullYear(), current.value.getMonth() - 1, 1); selectedDay.value = null }
const nextMonth = () => { current.value = new Date(current.value.getFullYear(), current.value.getMonth() + 1, 1); selectedDay.value = null }
const goToday = () => { current.value = new Date(); selectedDay.value = new Date().getDate() }

const getPlantName = (id) => {
  const p = plants.value.find(x => x.id === id)
  return p ? p.name : '未知植物'
}
const goPlant = (id) => router.push(`/plants/${id}`)
const formatTime = (t) => {
  const d = new Date(t)
  return `${String(d.getHours()).padStart(2,'0')}:${String(d.getMinutes()).padStart(2,'0')}`
}

const loadAll = async () => {
  loading.value = true
  try {
    const [c, p] = await Promise.all([careApi.list(), plantApi.list()])
    records.value = c.data || []
    plants.value = p.data || []
    selectedDay.value = new Date().getDate()
    if (current.value.getMonth() !== new Date().getMonth()) {
      // do nothing
    }
  } catch (e) {
    toast('加载失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)
</script>

<style scoped>
.month-nav { display: flex; gap: 8px; align-items: center; }
.month-nav button {
  width: 32px; height: 32px;
  border-radius: 50%;
  background: var(--card);
  font-size: 18px;
  color: var(--text);
  box-shadow: var(--shadow-sm);
  display: inline-flex;
  align-items: center;
  justify-content: center;
}
.month-nav button i { font-size: 22px; }
.today-btn {
  width: auto !important;
  padding: 0 14px;
  border-radius: 16px;
  font-size: 13px;
  font-weight: 600;
}

.calendar {
  background: var(--card);
  border-radius: var(--radius-lg);
  padding: 12px 8px;
  box-shadow: var(--shadow);
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 4px;
  margin-bottom: 14px;
}
.cal-weekday {
  text-align: center;
  font-size: 12px;
  font-weight: 700;
  color: var(--text-secondary);
  padding: 6px 0;
}
.cal-day {
  aspect-ratio: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-sm);
  cursor: pointer;
  position: relative;
  transition: all 0.2s;
}
.cal-day.blank { cursor: default; }
.cal-day:not(.blank):active { background: var(--primary-soft); }
.cal-day.today .day-num {
  background: var(--primary);
  color: #fff;
  border-radius: 50%;
  width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
}
.cal-day.hasRecords { background: var(--primary-soft); }
.cal-day.selected { background: var(--primary); }
.cal-day.selected .day-num { color: #fff; }
.cal-day.selected.today .day-num { background: #fff; color: var(--primary); }
.day-num { font-size: 14px; font-weight: 600; }
.day-dots {
  display: flex;
  gap: 2px;
  margin-top: 3px;
  height: 6px;
}
.dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  display: inline-block;
}
.dot.water { background: var(--water); }
.dot.fertilize { background: var(--fertilize); }
.dot.spray { background: var(--spray); }

.legend {
  display: flex;
  gap: 16px;
  justify-content: center;
  margin-bottom: 16px;
}
.legend-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-secondary);
}

.timeline { display: flex; flex-direction: column; }
.timeline-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
}
.timeline-item:last-child { border-bottom: none; }
.timeline-dot {
  width: 32px; height: 32px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
}
.timeline-dot.water { background: var(--water-soft); }
.timeline-dot.fertilize { background: var(--fertilize-soft); }
.timeline-dot.spray { background: var(--spray-soft); }
.timeline-content { flex: 1; }
.timeline-header { display: flex; gap: 8px; align-items: center; margin-bottom: 4px; flex-wrap: wrap; }
.plant-link { font-weight: 700; }
.timeline-time { font-size: 12px; color: var(--text-light); margin-left: auto; }
.timeline-remark { font-size: 13px; color: var(--text-secondary); }
</style>
