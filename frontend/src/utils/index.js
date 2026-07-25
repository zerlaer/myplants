// 全局 Toast 提示
let wrap = null

function getWrap() {
  if (!wrap) {
    wrap = document.createElement('div')
    wrap.className = 'toast-wrap'
    document.body.appendChild(wrap)
  }
  return wrap
}

export function toast(message, duration = 2200) {
  const el = document.createElement('div')
  el.className = 'toast'
  el.textContent = message
  getWrap().appendChild(el)
  setTimeout(() => {
    el.style.opacity = '0'
    el.style.transform = 'translateY(-12px)'
    el.style.transition = 'all 0.3s'
    setTimeout(() => el.remove(), 300)
  }, duration)
}

// 格式化日期 (移动端显示格式: M月DD日)
export function formatDate(date, withTime = false) {
  if (!date) return ''
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  let s = `${d.getMonth() + 1}月${pad(d.getDate())}日`
  if (withTime) {
    s += ` ${pad(d.getHours())}:${pad(d.getMinutes())}`
  }
  return s
}

// 格式化日期为 ISO 格式 (用于 <input type="date">)
export function formatDateISO(date) {
  if (!date) return ''
  const d = new Date(date)
  if (isNaN(d.getTime())) return ''
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

// 相对时间
export function timeAgo(date) {
  if (!date) return '从未'
  const d = new Date(date)
  const diff = Date.now() - d.getTime()
  if (diff < 0) return '刚刚'
  const min = Math.floor(diff / 60000)
  if (min < 1) return '刚刚'
  if (min < 60) return `${min}分钟前`
  const hour = Math.floor(min / 60)
  if (hour < 24) return `${hour}小时前`
  const day = Math.floor(hour / 24)
  if (day < 30) return `${day}天前`
  const month = Math.floor(day / 30)
  if (month < 12) return `${month}个月前`
  return `${Math.floor(month / 12)}年前`
}

// 距今天数(可负)
export function daysFromNow(date) {
  if (!date) return null
  const d = new Date(date)
  const diff = d.getTime() - Date.now()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
}

// 养护类型显示 (IconPark 图标名称)
export const careTypeMap = {
  water: { label: '浇水', icon: 'kettle-one', color: 'water' },
  fertilize: { label: '施肥', icon: 'pills', color: 'fertilize' },
  spray: { label: '打药', icon: 'medicine-bottle-one', color: 'spray' }
}

// 植物分类: 名称 → IconPark 图标名 (value 为图标名,模板中用 v-icon 渲染)
export const categoryMap = {
  '绿植': 'natural-mode',
  '多肉': 'cactus',
  '花卉': 'bloom',
  '草本': 'seedling',
  '木本': 'tree-two',
  '果树': 'fruiter'
}

export const healthStatusMap = {
  '优秀': { color: 'status-excellent', icon: '🌟' },
  '良好': { color: 'status-good', icon: '✨' },
  '一般': { color: 'status-normal', icon: '⚡' },
  '需关注': { color: 'status-warning', icon: '⚠️' }
}

export const materialMap = {
  '塑料': '🥤', '陶土': '🏺', '陶瓷': '🍵', '水泥': '🧱', '木质': '🪵', '其他': '📦'
}
