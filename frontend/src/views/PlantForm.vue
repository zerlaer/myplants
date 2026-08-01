<template>
  <div class="page">
    <div class="page-header">
      <button class="back-btn" @click="$router.back()"><i v-icon="'left'"></i></button>
      <div class="page-title">{{ isEdit ? '编辑植物' : '添加植物' }}</div>
      <div style="width:32px"></div>
    </div>

    <div class="card">
      <!-- 头像 -->
      <div class="avatar-upload" @click="triggerAvatarUpload">
        <div class="avatar-preview">
          <img v-if="form.avatar" :src="form.avatar" />
          <i v-else class="avatar-placeholder" v-icon="'natural-mode'"></i>
        </div>
        <div class="avatar-tip">点击上传头像</div>
        <input ref="avatarInput" type="file" accept="image/*" hidden @change="onAvatarChange" />
      </div>

      <div class="form-group">
        <label class="form-label">名称 <span class="req">*</span></label>
        <input v-model="form.name" class="form-input" placeholder="给植物起个名字" />
      </div>

      <div class="form-group">
        <label class="form-label">品种</label>
        <input v-model="form.species" class="form-input" placeholder="如：龟背竹、绿萝" />
      </div>

      <div class="form-group">
        <label class="form-label">分类</label>
        <div class="category-grid">
          <button v-for="c in categories" :key="c" type="button" class="category-btn"
                  :class="{ active: form.category === c }"
                  @click="form.category = c">
            <i class="cat-icon" v-icon="categoryMap[c]"></i>
            <span>{{ c }}</span>
          </button>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">摆放位置</label>
        <input v-model="form.location" class="form-input" list="location-options" placeholder="可输入或选择，如：东南阳台" />
        <datalist id="location-options">
          <option value="东南阳台"></option>
          <option value="西南阳台"></option>
          <option value="东北阳台"></option>
          <option value="西北阳台"></option>
          <option value="客厅"></option>
          <option value="卧室"></option>
          <option value="书房"></option>
          <option value="厨房"></option>
          <option value="窗台"></option>
        </datalist>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">到家时间</label>
          <div class="date-input-wrap">
            <input ref="dateInput" v-model="acquiredAt" type="date" class="date-input" @change="onDateChange" />
            <div class="date-display" @click="triggerDatePicker">
              <span class="date-text">{{ formattedDate || '选择日期' }}</span>
              <i class="date-icon" v-icon="'calendar'"></i>
            </div>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">植物状态</label>
          <select v-model="form.health_status" class="form-select">
            <option value="长势良好">长势良好</option>
            <option value="正在缓苗">正在缓苗</option>
            <option value="生长缓慢">生长缓慢</option>
            <option value="状态不佳">状态不佳</option>
            <option value="生病枯萎">生病枯萎</option>
            <option value="含苞待放">含苞待放</option>
            <option value="已经开花">已经开花</option>
            <option value="已经结果">已经结果</option>
          </select>
        </div>
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">购买价格</label>
          <div class="price-input-wrap">
            <input v-model.number="form.price" type="number" min="0" step="0.01" class="form-input" placeholder="0.00" />
            <span class="price-unit">元</span>
          </div>
        </div>
        <div class="form-group">
          <label class="form-label">选择花盆</label>
          <select v-model="form.pot_id" class="form-select">
            <option :value="null">不关联花盆</option>
            <option v-for="pot in pots" :key="pot.id" :value="pot.id">
              {{ pot.name || `花盆#${pot.id}` }}{{ pot.type ? ` (${pot.type})` : '' }}
            </option>
          </select>
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">光照需求</label>
        <div class="radio-group">
          <button v-for="l in lightOptions" :key="l" class="radio-btn"
                  :class="{ active: form.light_requirement === l }"
                  @click="form.light_requirement = l">{{ l }}</button>
        </div>
      </div>

      <div class="section-title" style="margin-top:8px"><i v-icon="'refresh'"></i> 养护周期(天)</div>
      <div class="cycle-grid">
        <div class="cycle-item">
          <i class="cycle-icon" v-icon="'kettle-one'"></i>
          <div class="cycle-label">浇水</div>
          <input v-model.number="form.water_cycle" type="number" min="1" class="cycle-input" />
          <span class="cycle-unit">天</span>
        </div>
        <div class="cycle-item">
          <i class="cycle-icon" v-icon="'pills'"></i>
          <div class="cycle-label">施肥</div>
          <input v-model.number="form.fertilize_cycle" type="number" min="1" class="cycle-input" />
          <span class="cycle-unit">天</span>
        </div>
        <div class="cycle-item">
          <i class="cycle-icon" v-icon="'medicine-bottle-one'"></i>
          <div class="cycle-label">打药</div>
          <input v-model.number="form.spray_cycle" type="number" min="1" class="cycle-input" />
          <span class="cycle-unit">天</span>
        </div>
      </div>

      <div class="form-group" style="margin-top:16px">
        <label class="form-label">描述</label>
        <textarea v-model="form.description" class="form-textarea" placeholder="记录植物的特点、习性等..."></textarea>
      </div>
    </div>

    <div class="form-actions">
      <button class="btn btn-ghost" @click="$router.back()">取消</button>
      <button class="btn btn-primary" @click="onSave" :disabled="saving">{{ saving ? '保存中...' : '保存' }}</button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { plantApi, photoApi, potApi, configApi } from '../api'
import { toast, categoryMap, formatDateISO } from '../utils'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const saving = ref(false)
const avatarInput = ref(null)
const dateInput = ref(null)
const acquiredAt = ref('')
const pots = ref([])

const formattedDate = computed(() => {
  if (!acquiredAt.value) return ''
  const parts = acquiredAt.value.split('-')
  if (parts.length === 3) {
    const month = parseInt(parts[1])
    const day = parseInt(parts[2])
    return `${month}月${day.toString().padStart(2, '0')}日`
  }
  return ''
})

const triggerDatePicker = () => dateInput.value?.click()

const onDateChange = () => {}

const categories = ['绿植', '多肉', '花卉', '草本', '木本', '果树']
const lightOptions = ['喜阳', '半阴', '喜阴']

const form = reactive({
  name: '', species: '', category: '', location: '', avatar: '',
  acquired_at: null, health_status: '长势良好', light_requirement: '半阴',
  water_cycle: 0, fertilize_cycle: 0, spray_cycle: 0, price: 0,
  pot_id: null, description: ''
})

watch(acquiredAt, (v) => {
  form.acquired_at = v ? new Date(v + 'T00:00:00') : null
})

const triggerAvatarUpload = () => avatarInput.value?.click()

const onAvatarChange = async (e) => {
  const file = e.target.files[0]
  if (!file) return
  const formData = new FormData()
  formData.append('file', file)
  try {
    const res = await photoApi.uploadAvatar(formData)
    form.avatar = res.data.path
    toast('头像上传成功')
  } catch (err) {
    toast('上传失败')
  }
  e.target.value = ''
}

const loadPlant = async () => {
  if (!isEdit.value) return
  try {
    const res = await plantApi.get(route.params.id)
    if (!res.data || !res.data.id) {
      toast('植物不存在')
      router.replace('/plants')
      return
    }
    Object.assign(form, res.data)
    if (form.acquired_at) {
      acquiredAt.value = formatDateISO(form.acquired_at)
    }
  } catch (e) {
    toast('加载失败: ' + (e.message || '植物不存在'))
    router.replace('/plants')
  }
}

const loadPots = async () => {
  try {
    const res = await potApi.list()
    pots.value = res.data || []
  } catch (e) {
    console.error('加载花盆失败', e)
  }
}

// 新建植物时从 config.yaml 读取默认养护周期
const loadConfig = async () => {
  if (isEdit.value) return
  try {
    const res = await configApi.get()
    const cfg = res.data || {}
    if (form.water_cycle === 0) form.water_cycle = cfg.default_water_days || 7
    if (form.fertilize_cycle === 0) form.fertilize_cycle = cfg.default_fertilize_days || 30
    if (form.spray_cycle === 0) form.spray_cycle = cfg.default_spray_days || 45
  } catch (e) {
    // 接口失败时使用兜底默认值
    form.water_cycle = form.water_cycle || 7
    form.fertilize_cycle = form.fertilize_cycle || 30
    form.spray_cycle = form.spray_cycle || 45
  }
}

const onSave = async () => {
  if (!form.name.trim()) {
    toast('请输入植物名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await plantApi.update(route.params.id, form)
      toast('更新成功 ✅')
    } else {
      const res = await plantApi.create(form)
      toast('添加成功 🎉')
      router.replace(`/plants/${res.data.id}`)
      return
    }
    router.back()
  } catch (e) {
    toast(e.message || '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadPots()
  loadPlant()
  loadConfig()
})
</script>

<style scoped>
.back-btn {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--card);
  font-size: 24px;
  color: var(--text);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-sm);
  line-height: 1;
}
.back-btn i { font-size: 22px; }

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
}
.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 30px;
  background: var(--primary-soft);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  border: 3px solid var(--card);
  box-shadow: var(--shadow);
  margin-bottom: 8px;
}
.avatar-preview img { width: 100%; height: 100%; object-fit: cover; }
.avatar-placeholder { font-size: 48px; color: var(--primary); }
.avatar-tip { font-size: 12px; color: var(--text-light); }

.radio-group { display: flex; gap: 8px; }
.radio-btn {
  flex: 1;
  padding: 10px;
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--text-secondary);
  font-weight: 600;
  border: 1.5px solid transparent;
  transition: all 0.2s;
}
.radio-btn.active {
  background: var(--primary-soft);
  color: var(--primary);
  border-color: var(--primary);
}

/* 分类图标按钮组 */
.category-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}
.category-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 6px;
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--text-secondary);
  border: 1.5px solid transparent;
  font-weight: 600;
  font-size: 13px;
  transition: all 0.2s;
}
.category-btn .cat-icon {
  font-size: 28px;
  color: var(--text-secondary);
  transition: color 0.2s;
}
.category-btn.active {
  background: var(--primary-soft);
  color: var(--primary);
  border-color: var(--primary);
}
.category-btn.active .cat-icon { color: var(--primary); }

.cycle-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}
.cycle-item {
  background: var(--bg);
  border-radius: var(--radius);
  padding: 14px 8px;
  text-align: center;
}
.cycle-icon { font-size: 26px; margin-bottom: 4px; color: var(--primary); display: inline-flex; }
.cycle-label { font-size: 12px; color: var(--text-secondary); margin-bottom: 8px; }
.cycle-input {
  width: 100%;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 6px;
  text-align: center;
  font-weight: 700;
  font-size: 16px;
  background: var(--card);
}
.cycle-unit { font-size: 11px; color: var(--text-light); }

.form-actions {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}
.form-actions .btn { flex: 1; }

.price-input-wrap {
  position: relative;
}
.price-input-wrap .form-input { padding-right: 36px; }
.price-unit {
  position: absolute;
  right: 14px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-light);
  font-size: 14px;
  font-weight: 600;
}

.date-input-wrap {
  position: relative;
}
.date-input {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
  z-index: 2;
}
.date-display {
  width: 100%;
  padding: 12px 14px;
  border: 1.5px solid var(--border);
  border-radius: var(--radius);
  background: var(--card);
  font-size: 15px;
  color: var(--text);
  box-sizing: border-box;
  height: 46px;
  line-height: 1.4;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.date-text {
  flex: 1;
  text-align: left;
}
.date-display .date-text:empty::before {
  content: '选择日期';
  color: var(--text-light);
}
.date-icon {
  font-size: 16px;
  color: var(--text-light);
  flex-shrink: 0;
  margin-left: 8px;
}
</style>
