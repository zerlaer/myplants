import axios from 'axios'

const request = axios.create({
  baseURL: '/api',
  timeout: 30000
})

// 响应拦截
request.interceptors.response.use(
  res => {
    const data = res.data
    if (data.code === 0) {
      return data
    }
    return Promise.reject(new Error(data.message || '请求失败'))
  },
  err => {
    return Promise.reject(err)
  }
)

export default request

// ===== 植物 =====
export const plantApi = {
  list: (params) => request.get('/plants', { params }),
  get: (id) => request.get(`/plants/${id}`),
  create: (data) => request.post('/plants', data),
  update: (id, data) => request.put(`/plants/${id}`, data),
  delete: (id) => request.delete(`/plants/${id}`),
  stats: (id) => request.get(`/plants/${id}/stats`)
}

// ===== 养护 =====
export const careApi = {
  list: (params) => request.get('/care', { params }),
  create: (data) => request.post('/care', data),
  delete: (id) => request.delete(`/care/${id}`),
  oneClick: (id, type) => request.post(`/care/one-click/${id}/${type}`),
  batch: (data) => request.post('/care/batch', data)
}

// ===== 提醒 =====
export const reminderApi = {
  list: () => request.get('/reminders')
}

// ===== 相册 =====
export const photoApi = {
  list: (params) => request.get('/photos', { params }),
  upload: (formData, onProgress) => request.post('/photos', formData, {
    timeout: 0,
    onUploadProgress: onProgress
  }),
  uploadAvatar: (formData) => request.post('/photos/avatar', formData, {
    timeout: 0
  }),
  update: (id, data) => request.put(`/photos/${id}`, data),
  delete: (id) => request.delete(`/photos/${id}`)
}

// ===== 花盆 =====
export const potApi = {
  list: (params) => request.get('/pots', { params }),
  create: (data) => request.post('/pots', data),
  update: (id, data) => request.put(`/pots/${id}`, data),
  delete: (id) => request.delete(`/pots/${id}`)
}

// ===== 换盆 =====
export const repotApi = {
  list: (params) => request.get('/repottings', { params }),
  create: (data) => request.post('/repottings', data),
  delete: (id) => request.delete(`/repottings/${id}`)
}

// ===== 笔记 =====
export const noteApi = {
  list: (params) => request.get('/notes', { params }),
  create: (data) => request.post('/notes', data),
  delete: (id) => request.delete(`/notes/${id}`)
}

// ===== 仪表盘 =====
export const dashboardApi = {
  get: () => request.get('/dashboard')
}
