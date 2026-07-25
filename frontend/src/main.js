import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/styles.css'

// IconPark 图标 (按需引入，tree-shaking 友好)
// 每个图标是一个函数，调用后返回 SVG 字符串，默认 size=1em, 颜色 currentColor
import Home from '@icon-park/svg/es/icons/Home.js'
import Leaf from '@icon-park/svg/es/icons/Leaf.js'
import Plus from '@icon-park/svg/es/icons/Plus.js'
import Remind from '@icon-park/svg/es/icons/Remind.js'
import Box from '@icon-park/svg/es/icons/Box.js'
import Water from '@icon-park/svg/es/icons/Water.js'
import Shield from '@icon-park/svg/es/icons/Shield.js'
import Search from '@icon-park/svg/es/icons/Search.js'
import Local from '@icon-park/svg/es/icons/Local.js'
import Caution from '@icon-park/svg/es/icons/Caution.js'
import Right from '@icon-park/svg/es/icons/Right.js'
import Left from '@icon-park/svg/es/icons/Left.js'
import ChartLine from '@icon-park/svg/es/icons/ChartLine.js'
import GridFour from '@icon-park/svg/es/icons/GridFour.js'
import Like from '@icon-park/svg/es/icons/Like.js'
import Lightning from '@icon-park/svg/es/icons/Lightning.js'
import Calendar from '@icon-park/svg/es/icons/Calendar.js'
import Camera from '@icon-park/svg/es/icons/Camera.js'
import Pic from '@icon-park/svg/es/icons/Pic.js'
import Refresh from '@icon-park/svg/es/icons/Refresh.js'
import Timer from '@icon-park/svg/es/icons/Timer.js'
import Book from '@icon-park/svg/es/icons/Book.js'
import Delete from '@icon-park/svg/es/icons/Delete.js'
import Close from '@icon-park/svg/es/icons/Close.js'
import Good from '@icon-park/svg/es/icons/Good.js'
import Sun from '@icon-park/svg/es/icons/Sun.js'
import Cloudy from '@icon-park/svg/es/icons/Cloudy.js'
import Round from '@icon-park/svg/es/icons/Round.js'
import Cup from '@icon-park/svg/es/icons/Cup.js'
import Paint from '@icon-park/svg/es/icons/Paint.js'
import Tag from '@icon-park/svg/es/icons/Tag.js'
import FileTxt from '@icon-park/svg/es/icons/FileTxt.js'
import Move from '@icon-park/svg/es/icons/Move.js'
import Cube from '@icon-park/svg/es/icons/Cube.js'
import Editor from '@icon-park/svg/es/icons/Editor.js'
// 新增图标: 植物/养护/相册/分类
import NaturalMode from '@icon-park/svg/es/icons/NaturalMode.js'
import Kettle from '@icon-park/svg/es/icons/Kettle.js'
import Pills from '@icon-park/svg/es/icons/Pills.js'
import MedicineBottleOne from '@icon-park/svg/es/icons/MedicineBottleOne.js'
import PictureAlbum from '@icon-park/svg/es/icons/PictureAlbum.js'
import Bloom from '@icon-park/svg/es/icons/Bloom.js'
import TreeTwo from '@icon-park/svg/es/icons/TreeTwo.js'
import Seedling from '@icon-park/svg/es/icons/Seedling.js'
import KettleOne from '@icon-park/svg/es/icons/KettleOne.js'
import Cactus from '@icon-park/svg/es/icons/Cactus.js'
import Fruiter from '@icon-park/svg/es/icons/Fruiter.js'
// 花盆类型图标
import Pot from '@icon-park/svg/es/icons/Pot.js'
import BottleOne from '@icon-park/svg/es/icons/BottleOne.js'
import Cylinder from '@icon-park/svg/es/icons/Cylinder.js'
import CupOne from '@icon-park/svg/es/icons/CupOne.js'

// 图标注册表: 模板中使用的名称 → IconPark 图标函数
// 包含别名 (模板名可能与 IconPark 原名不同)
const iconRegistry = {
  home: Home, leaf: Leaf, plus: Plus, remind: Remind, box: Box,
  water: Water, shield: Shield, search: Search, local: Local,
  caution: Caution, right: Right, left: Left,
  'chart-line': ChartLine, 'grid-four': GridFour, like: Like,
  flash: Lightning, calendar: Calendar, camera: Camera, pic: Pic,
  refresh: Refresh, timer: Timer, book: Book, delete: Delete,
  close: Close, smile: Good, sun: Sun, cloud: Cloudy, circle: Round,
  coffee: Cup, palette: Paint, tag: Tag, 'file-text': FileTxt,
  'vertical-align': Move, cube: Cube, editor: Editor, package: Box,
  // 植物与养护图标
  'natural-mode': NaturalMode, kettle: Kettle, pills: Pills,
  'medicine-bottle-one': MedicineBottleOne, 'picture-album': PictureAlbum,
  'kettle-one': KettleOne,
  // 植物分类图标
  cactus: Cactus, bloom: Bloom, seedling: Seedling,
  'tree-two': TreeTwo, fruiter: Fruiter,
  // 花盆类型图标
  pot: Pot, 'bottle-one': BottleOne, cylinder: Cylinder, 'cup-one': CupOne,
  // 别名 (兼容 IconPark 原名)
  lightning: Lightning, good: Good, cloudy: Cloudy, round: Round,
  cup: Cup, paint: Paint, 'file-txt': FileTxt, move: Move
}

const app = createApp(App)

// IconPark 指令: <i v-icon="'home'"></i> 或 <i v-icon="name"></i>
app.directive('icon', {
  mounted(el, binding) {
    renderIcon(el, binding.value)
  },
  updated(el, binding) {
    if (binding.value !== binding.oldValue) {
      renderIcon(el, binding.value)
    }
  }
})

function renderIcon(el, name) {
  el.classList.add('iconpark-icon')
  const fn = iconRegistry[name]
  if (fn) {
    el.innerHTML = fn({ size: '1em' })
  } else {
    el.innerHTML = ''
  }
}

app.use(createPinia())
app.use(router)
app.mount('#app')
