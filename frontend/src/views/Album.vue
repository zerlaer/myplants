<template>
  <div class="page">
    <div class="page-header">
      <button class="back-btn" @click="$router.back()"><i v-icon="'left'"></i></button>
      <div class="page-title"><i v-icon="'pic'"></i> 成长相册</div>
      <button class="btn btn-primary btn-sm" @click="showUpload = true"><i v-icon="'plus'"></i> 上传</button>
    </div>

    <div v-if="loading" class="loading"><div class="spinner"></div>加载中...</div>

    <template v-else>
      <!-- 时间线视图 -->
      <div v-if="photos.length" class="album-timeline">
        <div v-for="(group, month) in groupedPhotos" :key="month" class="album-group">
          <div class="group-header">
            <span class="group-dot"></span>
            <span class="group-label">{{ month }}</span>
            <span class="group-count">{{ group.length }} 张</span>
          </div>
          <div class="group-photos">
            <div v-for="p in group" :key="p.id" class="album-photo" @click="previewPhoto(p)">
              <img :src="p.path" :alt="p.remark" />
              <div class="photo-overlay">
                <div class="photo-date">{{ formatDate(p.taken_at) }}</div>
                <div class="photo-remark" v-if="p.remark">{{ p.remark }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="empty-state">
        <i class="emoji" v-icon="'picture-album'"></i>
        <div class="text">还没有照片，记录植物的成长瞬间吧～</div>
        <button class="btn btn-primary" @click="showUpload = true"><i v-icon="'plus'"></i> 上传照片</button>
      </div>
    </template>

    <!-- 上传弹窗 -->
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

    <!-- 预览 -->
    <div v-if="previewing" class="photo-preview" @click="previewing = null">
      <img :src="previewing.path" />
      <div class="preview-info" v-if="previewing.remark || previewing.taken_at">
        <div v-if="previewing.remark">{{ previewing.remark }}</div>
        <div class="text-light">{{ formatDate(previewing.taken_at) }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { photoApi } from '../api'
import { toast, formatDate } from '../utils'

const route = useRoute()
const loading = ref(true)
const photos = ref([])
const showUpload = ref(false)
const uploading = ref(false)
const previewing = ref(null)
const photoInput = ref(null)
const uploadForm = reactive({ remark: '', taken_at: '' })

const groupedPhotos = computed(() => {
  const groups = {}
  photos.value.forEach(p => {
    const d = new Date(p.taken_at)
    const key = `${d.getFullYear()}年${d.getMonth() + 1}月`
    if (!groups[key]) groups[key] = []
    groups[key].push(p)
  })
  // 按时间倒序
  return Object.keys(groups).sort((a, b) => {
    const ya = parseInt(a.match(/(\d+)年/)[1])
    const yb = parseInt(b.match(/(\d+)年/)[1])
    const ma = parseInt(a.match(/(\d+)月/)[1])
    const mb = parseInt(b.match(/(\d+)月/)[1])
    return yb - ya || mb - ma
  }).reduce((obj, key) => { obj[key] = groups[key]; return obj }, {})
})

const loadPhotos = async () => {
  loading.value = true
  try {
    const res = await photoApi.list({ plant_id: route.params.id })
    photos.value = res.data || []
  } catch (e) {
    toast('加载失败')
  } finally {
    loading.value = false
  }
}

const previewPhoto = (p) => { previewing.value = p }

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
    loadPhotos()
  } catch (e) {
    toast('上传失败')
  } finally {
    uploading.value = false
  }
}

onMounted(loadPhotos)
</script>

<style scoped>
.back-btn {
  width: 32px; height: 32px;
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

.album-timeline { position: relative; }
.album-group { margin-bottom: 20px; position: relative; }
.album-group::before {
  content: '';
  position: absolute;
  left: 7px;
  top: 24px;
  bottom: -20px;
  width: 2px;
  background: var(--border);
}
.album-group:last-child::before { display: none; }

.group-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  position: relative;
}
.group-dot {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--primary);
  border: 3px solid var(--card);
  box-shadow: 0 0 0 2px var(--primary);
  flex-shrink: 0;
}
.group-label { font-size: 16px; font-weight: 700; }
.group-count { font-size: 12px; color: var(--text-light); }

.group-photos {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding-left: 26px;
}
.album-photo {
  position: relative;
  aspect-ratio: 1;
  border-radius: var(--radius);
  overflow: hidden;
  cursor: pointer;
  box-shadow: var(--shadow-sm);
}
.album-photo img { width: 100%; height: 100%; object-fit: cover; transition: transform 0.3s; }
.album-photo:active img { transform: scale(1.08); }
.photo-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(transparent, rgba(0,0,0,0.7));
  color: #fff;
  padding: 20px 10px 8px;
}
.photo-date { font-size: 11px; font-weight: 600; }
.photo-remark { font-size: 12px; opacity: 0.9; margin-top: 2px; }

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
