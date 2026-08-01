<template>
  <div class="page detail-page" v-if="plant">
    <!-- 顶部封面 -->
    <div class="detail-cover">
      <div class="cover-bg" :style="coverStyle"></div>
      <button class="cover-back" @click="$router.back()"><i v-icon="'left'"></i></button>
      <button class="cover-edit" @click="$router.push(`/plants/${plant.id}/edit`)"><i v-icon="'editor'"></i></button>
      <div class="cover-content">
        <div class="cover-avatar">
          <img v-if="plant.avatar" :src="plant.avatar" />
          <i v-else v-icon="'natural-mode'"></i>
        </div>
        <h1 class="cover-name">{{ plant.name }}</h1>
        <div class="cover-meta">
          <span v-if="plant.species">{{ plant.species }}</span>
          <span class="tag" :class="healthTagClass(plant.health_status)">{{ plant.health_status }}</span>
        </div>
      </div>
    </div>

    <!-- 一键养护 -->
    <div class="quick-care-card">
      <div class="quick-care-title">一键养护</div>
      <div class="quick-care-actions">
        <button class="btn-care btn-water" @click="oneClick('water')" :disabled="caring.water">
          <i class="icon" v-icon="'kettle-one'"></i>
          <span>浇水</span>
          <span class="last-time" v-if="plant.last_watered_at">{{ timeAgo(plant.last_watered_at) }}</span>
        </button>
        <button class="btn-care btn-fertilize" @click="oneClick('fertilize')" :disabled="caring.fertilize">
          <i class="icon" v-icon="'pills'"></i>
          <span>施肥</span>
          <span class="last-time" v-if="plant.last_fertilized_at">{{ timeAgo(plant.last_fertilized_at) }}</span>
        </button>
        <button class="btn-care btn-spray" @click="oneClick('spray')" :disabled="caring.spray">
          <i class="icon" v-icon="'medicine-bottle-one'"></i>
          <span>打药</span>
          <span class="last-time" v-if="plant.last_sprayed_at">{{ timeAgo(plant.last_sprayed_at) }}</span>
        </button>
      </div>
    </div>

    <!-- 统计 -->
    <div class="stats-row">
      <div class="stat-item"><div class="stat-num">{{ stats.days_kept }}</div><div class="stat-label">养护天数</div></div>
      <div class="stat-item"><div class="stat-num">{{ stats.water_count }}</div><div class="stat-label">浇水</div></div>
      <div class="stat-item"><div class="stat-num">{{ stats.fertilize_count }}</div><div class="stat-label">施肥</div></div>
      <div class="stat-item"><div class="stat-num">{{ stats.photo_count }}</div><div class="stat-label">照片</div></div>
    </div>

    <!-- Tab -->
    <div class="tabs">
      <button v-for="t in tabs" :key="t.key" class="tab" :class="{ active: activeTab === t.key }" @click="activeTab = t.key">
        <i v-icon="t.icon"></i> {{ t.label }}
      </button>
    </div>

    <!-- 详情 -->
    <div v-if="activeTab === 'info'" class="tab-content">
      <div class="card">
        <div class="info-row">
          <span class="info-label"><i v-icon="'tag'"></i> 分类</span>
          <span>{{ plant.category || '未分类' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'local'"></i> 位置</span>
          <span>{{ plant.location || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'sun'"></i> 光照</span>
          <span>{{ plant.light_requirement || '未设置' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'calendar'"></i> 获得日期</span>
          <span>{{ formatDate(plant.acquired_at) || '未记录' }}</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'currency'"></i> 购买价格</span>
          <span class="price-value" v-if="plant.price > 0">¥{{ Number(plant.price).toFixed(2) }}</span>
          <span v-else>未记录</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'kettle-one'"></i> 浇水周期</span>
          <span>{{ plant.water_cycle }} 天</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'pills'"></i> 施肥周期</span>
          <span>{{ plant.fertilize_cycle }} 天</span>
        </div>
        <div class="info-row">
          <span class="info-label"><i v-icon="'medicine-bottle-one'"></i> 打药周期</span>
          <span>{{ plant.spray_cycle }} 天</span>
        </div>
        <div class="info-row" v-if="plant.description">
          <span class="info-label"><i v-icon="'file-text'"></i> 描述</span>
          <span class="info-desc">{{ plant.description }}</span>
        </div>
      </div>
      <button class="btn btn-danger btn-block" style="margin-top:14px" @click="onDelete"><i v-icon="'delete'"></i> 删除植物</button>
    </div>

    <!-- 相册 -->
    <div v-if="activeTab === 'album'" class="tab-content">
      <div class="card" style="padding:14px">
        <div class="flex-between mb-16">
          <span class="text-bold"><i v-icon="'pic'"></i> 成长记录 ({{ photos.length }})</span>
          <button class="btn btn-primary btn-sm" @click="showUpload = true"><i v-icon="'plus'"></i> 上传</button>
        </div>
        <div v-if="photos.length" class="photo-grid">
          <div v-for="p in photos" :key="p.id" class="photo-item" @click="previewPhoto(p)">
            <img :src="p.path" :alt="p.remark" />
            <div class="photo-date">{{ formatDate(p.taken_at) }}</div>
          </div>
        </div>
        <div v-else class="empty-state" style="padding:30px">
          <i class="emoji" v-icon="'picture-album'"></i>
          <div class="text">还没有照片</div>
        </div>
      </div>
    </div>

    <!-- 养护记录 -->
    <div v-if="activeTab === 'care'" class="tab-content">
      <div class="card" style="padding:14px">
        <div class="flex-between mb-16">
          <span class="text-bold"><i v-icon="'timer'"></i> 养护记录 ({{ careRecords.length }})</span>
          <button class="btn btn-primary btn-sm" @click="showCareForm = true"><i v-icon="'plus'"></i> 记录</button>
        </div>
        <div v-if="careRecords.length" class="timeline">
          <div v-for="r in careRecords" :key="r.id" class="timeline-item">
            <div class="timeline-dot" :class="careTypeMap[r.type].color"><i v-icon="careTypeMap[r.type].icon"></i></div>
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="tag" :class="'tag-' + careTypeMap[r.type].color">{{ careTypeMap[r.type].label }}</span>
                <span class="timeline-time">{{ formatDate(r.record_time, true) }}</span>
              </div>
              <div class="timeline-remark" v-if="r.remark">{{ r.remark }}</div>
            </div>
            <button class="timeline-del" @click="deleteCare(r)"><i v-icon="'close'"></i></button>
          </div>
        </div>
        <div v-else class="empty-state" style="padding:30px">
          <i class="emoji" v-icon="'timer'"></i>
          <div class="text">暂无养护记录</div>
        </div>
      </div>
    </div>

    <!-- 笔记 -->
    <div v-if="activeTab === 'note'" class="tab-content">
      <div class="card" style="padding:14px">
        <div class="flex-between mb-16">
          <span class="text-bold"><i v-icon="'book'"></i> 植物日记 ({{ notes.length }})</span>
          <button class="btn btn-primary btn-sm" @click="showNoteForm = true"><i v-icon="'plus'"></i> 写日记</button>
        </div>
        <div v-if="notes.length" class="note-list">
          <div v-for="n in notes" :key="n.id" class="note-item">
            <div class="note-content">{{ n.content }}</div>
            <div class="note-footer">
              <span>{{ formatDate(n.created_at, true) }}</span>
              <button @click="deleteNote(n)" class="note-del">删除</button>
            </div>
          </div>
        </div>
        <div v-else class="empty-state" style="padding:30px">
          <i class="emoji" v-icon="'book'"></i>
          <div class="text">还没有日记，记录植物的成长点滴～</div>
        </div>
      </div>
    </div>

    <!-- 换盆记录 -->
    <div v-if="activeTab === 'repot'" class="tab-content">
      <div class="card" style="padding:14px">
        <div class="flex-between mb-16">
          <span class="text-bold"><i v-icon="'box'"></i> 换盆记录 ({{ repottings.length }})</span>
          <button class="btn btn-primary btn-sm" @click="showRepotForm = true"><i v-icon="'plus'"></i> 换盆</button>
        </div>
        <div v-if="repottings.length" class="timeline">
          <div v-for="r in repottings" :key="r.id" class="timeline-item">
            <div class="timeline-dot fertilize"><i v-icon="'refresh'"></i></div>
            <div class="timeline-content">
              <div class="timeline-header">
                <span class="tag tag-fertilize">换盆</span>
                <span class="timeline-time">{{ formatDate(r.repot_time, true) }}</span>
              </div>
              <div class="timeline-remark">
                {{ r.from_pot_name || '原花盆' }} → {{ r.to_pot_name || '新花盆' }}
              </div>
              <div class="timeline-remark" v-if="r.remark">{{ r.remark }}</div>
            </div>
            <button class="timeline-del" @click="deleteRepot(r)"><i v-icon="'close'"></i></button>
          </div>
        </div>
        <div v-else class="empty-state" style="padding:30px">
          <i class="emoji" v-icon="'box'"></i>
          <div class="text">暂无换盆记录</div>
        </div>
      </div>
    </div>
  </div>

  <div v-if="!plant" class="loading"><div class="spinner"></div>加载中...</div>

  <!-- 上传照片弹窗 -->
  <div v-if="showUpload" class="modal-mask" @click.self="showUpload = false">
    <div class="modal">
      <div class="modal-handle"></div>
      <div class="modal-title">上传照片</div>
      <div class="form-group">
        <label class="form-label">选择照片</label>
        <input type="file" accept="image/*" ref="photoInput" class="form-input" style="padding:8px" />
      </div>
      <div class="form-group">
        <label class="form-label">备注</label>
        <input v-model="uploadForm.remark" class="form-input" placeholder="记录此刻..." />
      </div>
      <div class="form-group">
        <label class="form-label">拍摄时间</label>
        <input v-model="uploadForm.taken_at" type="datetime-local" class="form-input" />
      </div>
      <div class="modal-actions">
        <button class="btn btn-ghost" @click="showUpload = false">取消</button>
        <button class="btn btn-primary" @click="uploadPhoto" :disabled="uploading">{{ uploading ? '上传中...' : '上传' }}</button>
      </div>
    </div>
  </div>

  <!-- 养护记录弹窗 -->
  <div v-if="showCareForm" class="modal-mask" @click.self="showCareForm = false">
    <div class="modal">
      <div class="modal-handle"></div>
      <div class="modal-title">添加养护记录</div>
      <div class="form-group">
        <label class="form-label">类型</label>
        <div class="care-type-select">
          <button v-for="t in careTypes" :key="t.value" class="care-type-btn"
                  :class="{ active: careForm.type === t.value, ['btn-'+t.color]: careForm.type === t.value }"
                  @click="careForm.type = t.value">
            <i v-icon="t.icon"></i> {{ t.label }}
          </button>
        </div>
      </div>
      <div class="form-group">
        <label class="form-label">时间</label>
        <input v-model="careForm.record_time" type="datetime-local" class="form-input" />
      </div>
      <div class="form-group">
        <label class="form-label">备注</label>
        <input v-model="careForm.remark" class="form-input" placeholder="如：肥料名称、用量..." />
      </div>
      <div class="modal-actions">
        <button class="btn btn-ghost" @click="showCareForm = false">取消</button>
        <button class="btn btn-primary" @click="addCare">保存</button>
      </div>
    </div>
  </div>

  <!-- 笔记弹窗 -->
  <div v-if="showNoteForm" class="modal-mask" @click.self="showNoteForm = false">
    <div class="modal">
      <div class="modal-handle"></div>
      <div class="modal-title">写日记</div>
      <div class="form-group">
        <label class="form-label">内容</label>
        <textarea v-model="noteContent" class="form-textarea" placeholder="记录今天的观察、心情..." style="min-height:120px"></textarea>
      </div>
      <div class="modal-actions">
        <button class="btn btn-ghost" @click="showNoteForm = false">取消</button>
        <button class="btn btn-primary" @click="addNote">保存</button>
      </div>
    </div>
  </div>

  <!-- 换盆弹窗 -->
  <div v-if="showRepotForm" class="modal-mask" @click.self="showRepotForm = false">
    <div class="modal">
      <div class="modal-handle"></div>
      <div class="modal-title">换盆记录</div>
      <div class="form-group">
        <label class="form-label">原花盆</label>
        <select v-model="repotForm.from_pot_id" class="form-select">
          <option :value="null">不指定</option>
          <option v-for="p in pots" :key="p.id" :value="p.id">{{ potLabel(p) }}</option>
        </select>
      </div>
      <div class="form-group">
        <label class="form-label">新花盆</label>
        <select v-model="repotForm.to_pot_id" class="form-select">
          <option :value="null">不指定</option>
          <option v-for="p in pots" :key="p.id" :value="p.id">{{ potLabel(p) }}</option>
        </select>
      </div>
      <div class="form-group">
        <label class="form-label">换盆时间</label>
        <input v-model="repotForm.repot_time" type="datetime-local" class="form-input" />
      </div>
      <div class="form-group">
        <label class="form-label">备注</label>
        <input v-model="repotForm.remark" class="form-input" placeholder="如换盆原因、土壤..." />
      </div>
      <div class="modal-actions">
        <button class="btn btn-ghost" @click="showRepotForm = false">取消</button>
        <button class="btn btn-primary" @click="addRepot">保存</button>
      </div>
    </div>
  </div>

  <!-- 照片预览 -->
  <div v-if="previewing" class="photo-preview" @click="previewing = null">
    <img :src="previewing.path" />
    <div class="preview-info" v-if="previewing.remark || previewing.taken_at">
      <div v-if="previewing.remark">{{ previewing.remark }}</div>
      <div class="text-light">{{ formatDate(previewing.taken_at) }}</div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { plantApi, careApi, photoApi, noteApi, repotApi, potApi } from '../api'
import { toast, formatDate, timeAgo, careTypeMap, categoryMap } from '../utils'

const route = useRoute()
const router = useRouter()
const plant = ref(null)
const stats = ref({})
const photos = ref([])
const careRecords = ref([])
const notes = ref([])
const repottings = ref([])
const pots = ref([])
const caring = reactive({ water: false, fertilize: false, spray: false })

const activeTab = ref('info')
const tabs = [
  { key: 'info', label: '详情', icon: 'file-text' },
  { key: 'album', label: '相册', icon: 'pic' },
  { key: 'care', label: '养护', icon: 'kettle-one' },
  { key: 'note', label: '日记', icon: 'book' },
  { key: 'repot', label: '换盆', icon: 'box' }
]

const coverStyle = computed(() => {
  if (plant.value?.avatar) return { backgroundImage: `url(${plant.value.avatar})` }
  return { background: 'linear-gradient(135deg, #2d8659, #4ba878)' }
})

const healthTagClass = (status) => ({
  '长势良好': 'tag-primary', '正在缓苗': 'tag-primary',
  '生长缓慢': 'tag-warning', '状态不佳': 'tag-warning',
  '生病枯萎': 'tag-danger', '含苞待放': 'tag-primary',
  '已经开花': 'tag-primary', '已经结果': 'tag-primary'
}[status] || 'tag-grey')

const careTypes = [
  { value: 'water', label: '浇水', icon: 'kettle-one', color: 'water' },
  { value: 'fertilize', label: '施肥', icon: 'pills', color: 'fertilize' },
  { value: 'spray', label: '打药', icon: 'medicine-bottle-one', color: 'spray' }
]

const showUpload = ref(false)
const showCareForm = ref(false)
const showNoteForm = ref(false)
const showRepotForm = ref(false)
const previewing = ref(null)
const photoInput = ref(null)
const uploading = ref(false)

const uploadForm = reactive({ remark: '', taken_at: '' })
const careForm = reactive({ type: 'water', record_time: '', remark: '' })
const noteContent = ref('')
const repotForm = reactive({ from_pot_id: null, to_pot_id: null, repot_time: '', remark: '' })

const potLabel = (p) => {
  let s = p.name || `花盆#${p.id}`
  if (p.diameter) s += ` (⌀${p.diameter}cm)`
  return s
}

const nowDatetime = () => {
  const d = new Date()
  const pad = n => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth()+1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

const loadAll = async () => {
  const id = route.params.id
  try {
    const [p, s, ph, cr, nt, rp, pt] = await Promise.all([
      plantApi.get(id), plantApi.stats(id), photoApi.list({ plant_id: id }),
      careApi.list({ plant_id: id }), noteApi.list({ plant_id: id }),
      repotApi.list({ plant_id: id }), potApi.list()
    ])
    plant.value = p.data
    stats.value = s.data
    photos.value = ph.data || []
    careRecords.value = cr.data || []
    notes.value = nt.data || []
    repottings.value = rp.data || []
    pots.value = pt.data || []
  } catch (e) {
    toast('加载失败')
  }
}

const oneClick = async (type) => {
  caring[type] = true
  try {
    await careApi.oneClick(plant.value.id, type)
    toast(careTypeMap[type].label + '完成 ✅')
    await loadAll()
  } catch (e) {
    toast(e.message || '操作失败')
  } finally {
    caring[type] = false
  }
}

const uploadPhoto = async () => {
  const file = photoInput.value?.files[0]
  if (!file) { toast('请选择照片'); return }
  uploading.value = true
  try {
    const fd = new FormData()
    fd.append('file', file)
    fd.append('plant_id', route.params.id)
    fd.append('remark', uploadForm.remark)
    if (uploadForm.taken_at) fd.append('taken_at', uploadForm.taken_at.replace('T', ' ') + ':00')
    await photoApi.upload(fd)
    toast('上传成功 📸')
    showUpload.value = false
    uploadForm.remark = ''
    uploadForm.taken_at = ''
    await loadAll()
  } catch (e) {
    toast('上传失败')
  } finally {
    uploading.value = false
  }
}

const previewPhoto = (p) => { previewing.value = p }

const addCare = async () => {
  try {
    await careApi.create({
      plant_id: Number(route.params.id),
      type: careForm.type,
      remark: careForm.remark,
      record_time: careForm.record_time ? new Date(careForm.record_time) : null
    })
    toast('记录已保存 ✅')
    showCareForm.value = false
    careForm.remark = ''
    careForm.record_time = ''
    await loadAll()
  } catch (e) {
    toast(e.message || '保存失败')
  }
}

const deleteCare = async (r) => {
  if (!confirm('删除这条记录？')) return
  try {
    await careApi.delete(r.id)
    toast('已删除')
    loadAll()
  } catch (e) { toast('删除失败') }
}

const addNote = async () => {
  if (!noteContent.value.trim()) { toast('请输入内容'); return }
  try {
    await noteApi.create({ plant_id: Number(route.params.id), content: noteContent.value })
    toast('日记已保存 📔')
    showNoteForm.value = false
    noteContent.value = ''
    loadAll()
  } catch (e) { toast('保存失败') }
}

const deleteNote = async (n) => {
  if (!confirm('删除这条日记？')) return
  try {
    await noteApi.delete(n.id)
    toast('已删除')
    loadAll()
  } catch (e) { toast('删除失败') }
}

const addRepot = async () => {
  try {
    await repotApi.create({
      plant_id: Number(route.params.id),
      from_pot_id: repotForm.from_pot_id ? Number(repotForm.from_pot_id) : null,
      to_pot_id: repotForm.to_pot_id ? Number(repotForm.to_pot_id) : null,
      remark: repotForm.remark,
      repot_time: repotForm.repot_time ? new Date(repotForm.repot_time) : null
    })
    toast('换盆记录已保存 🏺')
    showRepotForm.value = false
    repotForm.remark = ''
    repotForm.from_pot_id = null
    repotForm.to_pot_id = null
    repotForm.repot_time = ''
    await loadAll()
  } catch (e) { toast(e.message || '保存失败') }
}

const deleteRepot = async (r) => {
  if (!confirm('删除这条换盆记录？')) return
  try {
    await repotApi.delete(r.id)
    toast('已删除')
    loadAll()
  } catch (e) { toast('删除失败') }
}

const onDelete = async () => {
  if (!confirm(`确定删除「${plant.value.name}」？所有相关记录将一并删除。`)) return
  try {
    await plantApi.delete(plant.value.id)
    toast('已删除')
    router.replace('/plants')
  } catch (e) { toast('删除失败') }
}

onMounted(loadAll)
</script>

<style scoped>
.detail-page { padding: 0; padding-bottom: calc(var(--nav-height) + var(--safe-bottom) + 24px); }

.detail-cover {
  position: relative;
  height: 280px;
  overflow: hidden;
}
.cover-bg {
  position: absolute;
  inset: 0;
  background-size: cover;
  background-position: center;
  filter: blur(2px);
  transform: scale(1.1);
}
.cover-bg::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.5) 100%);
}
.cover-back, .cover-edit {
  position: absolute;
  top: 16px;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(255,255,255,0.25);
  backdrop-filter: blur(10px);
  color: #fff;
  font-size: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
}
.cover-back { left: 16px; }
.cover-back i { font-size: 24px; }
.cover-edit { right: 16px; }
.cover-edit i { font-size: 18px; }
.cover-content {
  position: absolute;
  bottom: 40px;
  left: 0;
  right: 0;
  text-align: center;
  color: #fff;
  z-index: 2;
}
.cover-avatar {
  width: 84px;
  height: 84px;
  border-radius: 26px;
  margin: 0 auto 10px;
  background: rgba(255,255,255,0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 3px solid rgba(255,255,255,0.6);
  font-size: 40px;
}
.cover-avatar i { color: #fff; font-size: 40px; }
.cover-avatar img { width: 100%; height: 100%; object-fit: cover; }
.cover-name { font-size: 22px; font-weight: 800; margin-bottom: 4px; }
.cover-meta { display: flex; gap: 8px; justify-content: center; align-items: center; font-size: 13px; opacity: 0.95; }
.cover-meta .tag { background: rgba(255,255,255,0.25); color: #fff; }

.quick-care-card {
  margin: -10px 16px 0;
  position: relative;
  z-index: 3;
  background: var(--card);
  border-radius: var(--radius-lg);
  padding: 16px;
  box-shadow: var(--shadow-lg);
}
.quick-care-title { font-size: 13px; font-weight: 700; color: var(--text-secondary); margin-bottom: 12px; text-align: center; }
.quick-care-actions { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.btn-care { padding: 14px 6px; }
.btn-care .last-time {
  font-size: 10px;
  font-weight: 500;
  opacity: 0.7;
  margin-top: 2px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  margin: 16px;
  background: var(--card);
  border-radius: var(--radius);
  padding: 14px 0;
  box-shadow: var(--shadow-sm);
}
.stat-item { text-align: center; border-right: 1px solid var(--border); }
.stat-item:last-child { border-right: none; }
.stat-num { font-size: 20px; font-weight: 800; color: var(--primary); }
.stat-label { font-size: 11px; color: var(--text-secondary); margin-top: 2px; }

.tabs {
  display: flex;
  margin: 0 16px 16px;
  background: var(--card);
  border-radius: var(--radius);
  padding: 4px;
  box-shadow: var(--shadow-sm);
  gap: 2px;
}
.tab {
  flex: 1;
  padding: 9px 4px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary);
  transition: all 0.2s;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  line-height: 1;
}
.tab i { font-size: 16px; }
.tab.active {
  background: var(--primary);
  color: #fff;
}

.tab-content { padding: 0 16px; }

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
  gap: 12px;
}
.info-row:last-child { border-bottom: none; }
.info-label { font-weight: 600; color: var(--text-secondary); white-space: nowrap; flex-shrink: 0; display: inline-flex; align-items: center; gap: 6px; }
.info-label i { font-size: 16px; color: var(--primary); }
.info-desc { text-align: right; color: var(--text); }
.price-value { color: var(--accent); font-weight: 700; }

/* 标题(图标+文字)垂直对齐 */
.text-bold { display: inline-flex; align-items: center; gap: 6px; line-height: 1.2; }
.text-bold i { font-size: 18px; color: var(--primary); display: inline-flex; }

.photo-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.photo-item {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius-sm);
  overflow: hidden;
  cursor: pointer;
}
.photo-item img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.photo-item:active img { transform: scale(1.1); }
.photo-date {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.7));
  color: #fff;
  font-size: 10px;
  padding: 12px 6px 4px;
  text-align: center;
}

.timeline { display: flex; flex-direction: column; gap: 4px; }
.timeline-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 0;
  border-bottom: 1px solid var(--border);
  position: relative;
}
.timeline-item:last-child { border-bottom: none; }
.timeline-dot {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  flex-shrink: 0;
}
.timeline-dot.water { background: var(--water-soft); }
.timeline-dot.fertilize { background: var(--fertilize-soft); }
.timeline-dot.spray { background: var(--spray-soft); }
.timeline-content { flex: 1; min-width: 0; }
.timeline-header { display: flex; align-items: center; gap: 8px; margin-bottom: 4px; }
.timeline-time { font-size: 12px; color: var(--text-light); }
.timeline-remark { font-size: 13px; color: var(--text-secondary); }
.timeline-del {
  color: var(--text-light);
  font-size: 14px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
}
.timeline-del:active { background: #fdeaea; color: var(--danger); }

.note-list { display: flex; flex-direction: column; gap: 12px; }
.note-item {
  background: var(--bg);
  border-radius: var(--radius);
  padding: 12px;
}
.note-content { color: var(--text); margin-bottom: 8px; white-space: pre-wrap; }
.note-footer { display: flex; justify-content: space-between; align-items: center; font-size: 12px; color: var(--text-light); }
.note-del { color: var(--text-light); font-size: 12px; }

.care-type-select { display: grid; grid-template-columns: repeat(3, 1fr); gap: 8px; }
.care-type-btn {
  padding: 12px;
  border-radius: var(--radius);
  background: var(--bg);
  font-weight: 600;
  border: 1.5px solid transparent;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}
.care-type-btn i { font-size: 18px; }
.care-type-btn.active { border-color: currentColor; }

.photo-preview {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.92);
  z-index: 2000;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  animation: fadeIn 0.2s;
}
.photo-preview img { max-width: 100%; max-height: 80vh; border-radius: var(--radius); }
.preview-info { color: #fff; text-align: center; margin-top: 16px; }
</style>
