# 我的花园 (MyPlants)

一个基于 Go + Vue 3 的绿植管理系统，帮助您记录和管理家中的植物。

## 功能特性

- **植物管理**：添加、编辑、删除植物，记录名称、品种、分类、价格等信息
- **养护提醒**：浇水、施肥、打药的周期管理和智能提醒
- **一键养护**：快速记录养护操作，自动更新下次提醒时间
- **成长相册**：上传植物照片，记录成长过程
- **花盆管理**：管理花盆信息，关联植物与花盆
- **健康状态**：记录植物健康状况，支持多状态标记
- **数据统计**：植物总价值、分类统计、养护趋势图表

## 技术栈

### 后端
- Go 1.21+
- Gin 框架
- GORM ORM
- SQLite 数据库（纯 Go 驱动）
- Viper 配置管理
- Zap 日志

### 前端
- Vue 3 + Vite
- Vue Router
- Pinia 状态管理
- Axios HTTP 客户端
- IconPark 图标库

## 快速开始

### 环境要求
- Go 1.21+
- Node.js 18+

### 安装依赖

```bash
# 后端依赖
go mod download

# 前端依赖
cd frontend
npm install
```

### 开发模式

```bash
# 启动后端服务（默认端口 8080）
go run main.go

# 启动前端开发服务器（默认端口 5173）
cd frontend
npm run dev
```

### 生产构建

```bash
# 构建前端
cd frontend
npm run build

# 构建后端（Windows）
go build -ldflags "-s -w" -o myplants.exe

# 构建后端（Linux）
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o myplants

# 构建后端（macOS）
GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o myplants
```

### 运行

```bash
# Windows
./myplants.exe

# Linux/macOS
./myplants
```

访问 http://localhost:8080 即可使用。

## 配置说明

配置文件位于 `config.yaml`：

```yaml
server:
  port: 8080
  mode: release

database:
  path: ./data/myplants.db

logging:
  level: info
  path: ./logs/app.log
```

## 项目结构

```
myplants/
├── frontend/                 # 前端代码
│   ├── src/
│   │   ├── api/              # API 请求封装
│   │   ├── assets/           # 静态资源
│   │   ├── components/       # 公共组件
│   │   ├── router/           # 路由配置
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   └── dist/                 # 构建产物
├── internal/                 # 后端代码
│   ├── config/               # 配置管理
│   ├── controller/           # 控制器
│   ├── database/             # 数据库连接
│   ├── logger/               # 日志管理
│   ├── model/                # 数据模型
│   ├── response/             # 响应封装
│   └── router/               # 路由注册
├── config.yaml               # 配置文件
├── go.mod                    # Go 模块依赖
└── main.go                   # 入口文件
```

## API 接口

| 接口 | 方法 | 描述 |
|------|------|------|
| `/api/plants` | GET | 获取植物列表 |
| `/api/plants/:id` | GET | 获取植物详情 |
| `/api/plants` | POST | 创建植物 |
| `/api/plants/:id` | PUT | 更新植物 |
| `/api/plants/:id` | DELETE | 删除植物 |
| `/api/care/one-click` | POST | 一键养护 |
| `/api/care` | GET | 获取养护记录 |
| `/api/photos` | POST | 上传照片 |
| `/api/pots` | GET | 获取花盆列表 |
| `/api/dashboard` | GET | 获取首页数据 |

## 植物分类

- 绿植
- 多肉
- 花卉
- 草本
- 木本
- 果树

## 健康状态

- 优秀
- 良好
- 一般
- 需关注

## 许可证

MIT License