<template>
  <div class="page">
    <div class="page-header">
      <div>
        <div class="page-title"><i v-icon="'box'"></i> 花盆管理</div>
        <div class="page-subtitle">{{ pots.length }} 个花盆 · 使用中 {{ inUseCount }}</div>
      </div>
      <button class="btn btn-primary btn-sm" @click="openForm()"><i v-icon="'plus'"></i> 添加</button>
    </div>

    <!-- 状态筛选 -->
    <div class="filter-chips">
      <button class="chip" :class="{ active: status === '' }" @click="status = ''">全部</button>
      <button class="chip" :class="{ active: status === '使用中' }" @click="status = '使用中'">使用中</button>
      <button class="chip" :class="{ active: status === '空闲' }" @click="status = '空闲'">空闲</button>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <div v-else-if="filtered.length" class="pot-grid">
      <div v-for="pot in filtered" :key="pot.id" class="pot-card list-item" @click="openForm(pot)">
        <i class="pot-emoji" v-icon="materialIcon(pot.material)"></i>
        <div class="pot-info">
          <div class="pot-name">{{ pot.name || `花盆#${pot.id}` }}</div>
          <div class="pot-meta">
            <span v-if="pot.diameter"><i v-icon="'circle'"></i> {{ pot.diameter }}cm</span>
            <span v-if="pot.height"><i v-icon="'vertical-align'"></i> {{ pot.height }}cm</span>
            <span v-if="pot.material">{{ pot.material }}</span>
          </div>
          <div class="pot-meta">
            <span v-if="pot.color"><i v-icon="'palette'"></i> {{ pot.color }}</span>
            <span v-if="pot.plant_id">已关联植物</span>
          </div>
          <div class="pot-status">
            <span class="tag" :class="pot.status === '使用中' ? 'tag-primary' : 'tag-grey'">{{ pot.status }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="empty-state">
      <i class="emoji" v-icon="'box'"></i>
      <div class="text">还没有花盆记录</div>
      <button class="btn btn-primary" @click="openForm()"><i v-icon="'plus'"></i> 添加花盆</button>
    </div>

    <!-- 换盆记录 -->
    <div v-if="repottings.length" class="section-title"><i v-icon="'timer'"></i> 换盆记录</div>
    <div v-if="repottings.length" class="card" style="padding:14px">
      <div class="timeline">
        <div v-for="r in repottings" :key="r.id" class="timeline-item">
          <div class="timeline-dot fertilize"><i v-icon="'refresh'"></i></div>
          <div class="timeline-content">
            <div class="timeline-header">
              <span class="tag tag-fertilize">换盆</span>
              <span class="timeline-time">{{ formatDate(r.repot_time) }}</span>
            </div>
            <div class="timeline-remark">{{ r.from_pot_name || '原花盆' }} → {{ r.to_pot_name || '新花盆' }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 表单弹窗 -->
    <div v-if="showForm" class="modal-mask" @click.self="showForm = false">
      <div class="modal">
        <div class="modal-handle"></div>
        <div class="modal-title">{{ editing ? '编辑花盆' : '添加花盆' }}</div>
        <div class="form-group">
          <label class="form-label">名称</label>
          <input v-model="form.name" class="form-input" placeholder="给花盆起个名字" />
        </div>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">直径 (cm)</label>
            <input v-model.number="form.diameter" type="number" class="form-input" />
          </div>
          <div class="form-group">
            <label class="form-label">高度 (cm)</label>
            <input v-model.number="form.height" type="number" class="form-input" />
          </div>
        </div>
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">材质</label>
            <select v-model="form.material" class="form-select">
              <option value="">请选择</option>
              <option v-for="m in materials" :key="m" :value="m">{{ materialMap[m] }} {{ m }}</option>
            </select>
          </div>
          <div class="form-group">
            <label class="form-label">颜色</label>
            <input v-model="form.color" class="form-input" placeholder="如：白色" />
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">状态</label>
          <select v-model="form.status" class="form-select">
            <option value="使用中">使用中</option>
            <option value="空闲">空闲</option>
            <option value="已弃用">已弃用</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">关联植物</label>
          <select v-model="form.plant_id" class="form-select">
            <option :value="null">不关联</option>
            <option v-for="p in plants" :key="p.id" :value="p.id">{{ p.name }}</option>
          </select>
        </div>
        <div class="form-group">
          <label class="form-label">备注</label>
          <input v-model="form.remark" class="form-input" />
        </div>
        <div class="modal-actions">
          <button class="btn btn-danger" v-if="editing" @click="onDelete">删除</button>
          <button class="btn btn-ghost" @click="showForm = false">取消</button>
          <button class="btn btn-primary" @click="onSave">保存</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { potApi, plantApi, repotApi } from '../api'
import { toast, formatDate, materialMap } from '../utils'

const pots = ref([])
const plants = ref([])
const repottings = ref([])
const loading = ref(true)
const status = ref('')
const showForm = ref(false)
const editing = ref(null)

const materials = ['塑料', '陶土', '陶瓷', '水泥', '木质', '其他']

const materialIconMap = {
  '塑料': 'coffee',
  '陶土': 'cube',
  '陶瓷': 'coffee',
  '水泥': 'grid-four',
  '木质': 'package',
  '其他': 'box'
}
const materialIcon = (m) => materialIconMap[m] || 'box'

const form = reactive({
  name: '', diameter: 0, height: 0, material: '', color: '',
  status: '使用中', plant_id: null, remark: ''
})

const filtered = computed(() => {
  if (!status.value) return pots.value
  return pots.value.filter(p => p.status === status.value)
})

const inUseCount = computed(() => pots.value.filter(p => p.status === '使用中').length)

const resetForm = () => {
  Object.assign(form, { name: '', diameter: 0, height: 0, material: '', color: '', status: '使用中', plant_id: null, remark: '' })
}

const openForm = (pot = null) => {
  resetForm()
  if (pot) {
    editing.value = pot
    Object.assign(form, pot)
  } else {
    editing.value = null
  }
  showForm.value = true
}

const loadAll = async () => {
  loading.value = true
  try {
    const [p, pl, rp] = await Promise.all([potApi.list(), plantApi.list(), repotApi.list()])
    pots.value = p.data || []
    plants.value = pl.data || []
    repottings.value = rp.data || []
  } catch (e) {
    toast('加载失败')
  } finally {
    loading.value = false
  }
}

const onSave = async () => {
  try {
    const data = { ...form, size: form.diameter ? `⌀${form.diameter}cm` : '' }
    if (editing.value) {
      await potApi.update(editing.value.id, data)
      toast('更新成功 ✅')
    } else {
      await potApi.create(data)
      toast('添加成功 🏺')
    }
    showForm.value = false
    loadAll()
  } catch (e) {
    toast(e.message || '保存失败')
  }
}

const onDelete = async () => {
  if (!confirm('删除这个花盆？')) return
  try {
    await potApi.delete(editing.value.id)
    toast('已删除')
    showForm.value = false
    loadAll()
  } catch (e) { toast('删除失败') }
}

onMounted(loadAll)
</script>

<style scoped>
.filter-chips {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}
.chip {
  padding: 7px 16px;
  border-radius: 20px;
  background: var(--card);
  color: var(--text-secondary);
  font-size: 13px;
  font-weight: 600;
  border: 1.5px solid var(--border);
}
.chip.active { background: var(--primary); color: #fff; border-color: var(--primary); }

.pot-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}
.pot-card {
  background: var(--card);
  border-radius: var(--radius-lg);
  padding: 16px;
  box-shadow: var(--shadow);
  cursor: pointer;
  transition: transform 0.2s;
}
.pot-card:active { transform: scale(0.97); }
.pot-emoji { font-size: 36px; margin-bottom: 8px; color: var(--primary); display: inline-flex; }
.pot-name { font-weight: 700; font-size: 15px; margin-bottom: 6px; }
.pot-meta { display: flex; gap: 8px; font-size: 12px; color: var(--text-secondary); flex-wrap: wrap; margin-bottom: 4px; align-items: center; }
.pot-meta i { font-size: 14px; }
.pot-status { margin-top: 6px; }

.timeline { display: flex; flex-direction: column; }
.timeline-item {
  display: flex;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
}
.timeline-item:last-child { border-bottom: none; }
.timeline-dot {
  width: 32px; height: 32px;
  border-radius: 50%;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  background: var(--fertilize-soft);
  font-size: 16px;
  color: #5a9216;
}
.timeline-content { flex: 1; }
.timeline-header { display: flex; gap: 8px; align-items: center; margin-bottom: 4px; }
.timeline-time { font-size: 12px; color: var(--text-light); }
.timeline-remark { font-size: 13px; color: var(--text-secondary); }
</style>
