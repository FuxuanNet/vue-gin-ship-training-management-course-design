# 企业培训管理系统 - 前端

基于 Vue 3 + Vite + Element Plus 的企业培训管理系统前端应用，支持员工、讲师、课程大纲制定者三种角色的培训管理功能。

## 技术栈

- **框架**: Vue 3.5.25 (Composition API)
- **构建工具**: Vite 6.0.5
- **UI 组件库**: Element Plus 2.13.0
- **HTTP 客户端**: Axios 1.7.2
- **路由管理**: Vue Router 4.5.0
- **状态管理**: Pinia 2.3.0
- **图表可视化**: ECharts 5.6.0
- **开发工具**: Vue DevTools

## 项目结构

```
front/
├── public/                    # 静态资源目录
├── src/
│   ├── api/                   # API 接口模块
│   │   ├── auth.js           # 登录认证接口
│   │   ├── employee.js       # 员工端接口
│   │   ├── home.js           # 主页接口
│   │   ├── planner.js        # 课程大纲制定者接口
│   │   ├── request.js        # Axios 实例配置
│   │   └── teacher.js        # 讲师端接口
│   │
│   ├── components/            # 公共组件
│   │   └── CourseCard.vue    # 课程卡片组件
│   │
│   ├── router/                # 路由配置
│   │   └── index.js          # 路由定义和权限控制
│   │
│   ├── stores/                # Pinia 状态管理
│   │   ├── counter.js        # 计数器示例
│   │   ├── mockData.js       # 模拟数据
│   │   └── user.js           # 用户状态管理
│   │
│   ├── views/                 # 页面组件
│   │   ├── employee/         # 员工端页面
│   │   │   ├── Evaluation.vue    # 课程自评
│   │   │   ├── Schedule.vue      # 课程表
│   │   │   └── Scores.vue        # 成绩查看（含能力分析）
│   │   │
│   │   ├── planner/          # 课程大纲制定者页面
│   │   │   ├── Analytics.vue     # 数据分析（平台统计）
│   │   │   ├── Courses.vue       # 课程管理
│   │   │   └── Plans.vue         # 培训计划管理
│   │   │
│   │   ├── teacher/          # 讲师端页面
│   │   │   ├── Grading.vue       # 待评分学员
│   │   │   ├── Schedule.vue      # 授课表
│   │   │   └── Scores.vue        # 成绩统计（含可视化）
│   │   │
│   │   ├── Home.vue          # 系统主页
│   │   └── Login.vue         # 登录页面
│   │
│   ├── App.vue                # 根组件（含导航菜单）
│   └── main.js                # 应用入口
│
├── index.html                 # HTML 模板
├── vite.config.js            # Vite 配置
├── jsconfig.json             # JavaScript 配置
└── package.json              # 项目依赖

```

## 功能模块

### 1. 认证模块
- 用户登录（支持员工/讲师/课程大纲制定者）
- Session 会话管理
- 自动登录状态检查
- 登出功能

### 2. 员工端功能
- **课程表查看**: 查看个人课程安排，按周/日期筛选
- **课程自评**: 提交课程学习自评（AI 辅助打分）
- **成绩查看**: 
  - 成绩明细（自评分、讲师评分、综合得分）
  - 能力分析雷达图（按课程类型维度）
  - 各类型平均分统计

### 3. 讲师端功能
- **授课表查看**: 查看授课安排，显示学员数和评分进度
- **待评分学员**: 查看学员自评内容并提交讲师评分
- **成绩统计**:
  - 课程成绩统计（按课次、按学员）
  - 授课统计可视化（课程类型分布、月度授课时长）

### 4. 课程大纲制定者功能
- **培训计划管理**: 创建、编辑、删除培训计划，管理参与员工
- **课程管理**: 创建、编辑、删除课程，管理课程安排
- **数据分析**:
  - 课程评分排名 Top10
  - 培训计划排名 Top10
  - 员工成绩排名 Top10
  - 讲师授课统计
  - 课程类型分布
  - 培训计划状态统计

## 路由设计

```javascript
/login              # 登录页面（公开）
/home               # 系统主页（需认证）

# 员工端路由
/employee/schedule      # 课程表
/employee/evaluation    # 课程自评
/employee/scores        # 成绩查看

# 讲师端路由
/teacher/schedule       # 授课表
/teacher/grading        # 待评分学员
/teacher/scores         # 成绩统计

# 课程大纲制定者路由
/planner/plans          # 培训计划管理
/planner/courses        # 课程管理
/planner/analytics      # 数据分析
```

## API 接口设计

### 通用接口
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/check` - 检查登录状态
- `GET /api/home/courses` - 获取主页课程列表

### 员工端接口
- `GET /api/employee/schedule` - 获取课程表
- `GET /api/employee/pending-evaluations` - 获取待自评课程
- `POST /api/employee/submit-evaluation` - 提交自评
- `GET /api/employee/scores` - 获取成绩明细
- `GET /api/employee/course-type-scores` - 获取课程类型成绩分析
- `GET /api/employee/learning-progress` - 获取学习进度

### 讲师端接口
- `GET /api/teacher/schedule` - 获取授课表
- `GET /api/teacher/pending-evaluations` - 获取待评分学员
- `POST /api/teacher/submit-grading` - 提交讲师评分
- `GET /api/teacher/course-statistics` - 获取课程成绩统计
- `GET /api/teacher/teaching-statistics` - 获取授课统计

### 课程大纲制定者接口
- `GET /api/planner/plans` - 获取培训计划列表
- `POST /api/planner/plans` - 创建培训计划
- `PUT /api/planner/plans/:id` - 更新培训计划
- `DELETE /api/planner/plans/:id` - 删除培训计划
- `GET /api/planner/plans/:id` - 获取培训计划详情
- `POST /api/planner/plans/:id/employees` - 添加员工到计划
- `DELETE /api/planner/plans/:id/employees/:employeeId` - 从计划移除员工
- `GET /api/planner/courses` - 获取课程列表
- `POST /api/planner/courses` - 创建课程
- `PUT /api/planner/courses/:id` - 更新课程
- `DELETE /api/planner/courses/:id` - 删除课程
- `GET /api/planner/course-items` - 获取课程安排列表
- `POST /api/planner/course-items` - 创建课程安排
- `PUT /api/planner/course-items/:id` - 更新课程安排
- `DELETE /api/planner/course-items/:id` - 删除课程安排
- `GET /api/planner/employees` - 获取员工列表
- `GET /api/planner/teachers` - 获取讲师列表
- `GET /api/planner/analytics` - 获取平台数据分析
- `GET /api/planner/courses/:id/evaluations` - 获取课程评价统计
- `GET /api/planner/employees/:id/scores` - 获取员工成绩

## 开发环境配置

### 推荐 IDE

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) 

### 推荐浏览器插件

**Chrome/Edge:**
- [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)

**Firefox:**
- [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)

## 快速开始

### 安装依赖

```sh
npm install
```

### 开发模式运行

```sh
npm run dev
```

访问 http://localhost:5173

### 生产构建

```sh
npm run build
```

构建产物在 `dist/` 目录

## 配置说明

### API 基础路径

在 `src/api/request.js` 中配置后端 API 地址：

```javascript
const request = axios.create({
  baseURL: 'http://localhost:8080',  // 后端服务地址
  timeout: 10000,
  withCredentials: true
})
```

### Vite 代理配置

参考 `vite.config.js` 配置开发环境代理：

```javascript
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
```

## 核心技术实现

### 1. 状态管理 (Pinia)

使用 Pinia 管理全局用户状态：

```javascript
// stores/user.js
export const useUserStore = defineStore('user', {
  state: () => ({
    userInfo: null,
    isLoggedIn: false
  }),
  actions: {
    setUser(user) {
      this.userInfo = user
      this.isLoggedIn = true
    }
  }
})
```

### 2. 路由守卫

实现基于角色的路由权限控制：

```javascript
router.beforeEach((to, from, next) => {
  // 检查登录状态
  // 验证角色权限
  // 重定向到登录页或目标页
})
```

### 3. HTTP 拦截器

统一处理请求和响应：

```javascript
// 请求拦截器：添加认证信息
// 响应拦截器：统一错误处理
```

### 4. 图表可视化

使用 ECharts 实现数据可视化：
- 雷达图（员工能力分析）
- 柱状图（成绩统计）
- 饼图（课程类型分布）
- 折线图（月度授课趋势）

## 注意事项

1. **时间格式**: 时间使用 "HH:mm:ss" 字符串格式，避免时区转换问题
2. **评分逻辑**: 综合得分 = 自评分 × (1 - 评分占比) + 讲师评分 × 评分占比
3. **权限控制**: 不同角色只能访问对应的功能模块
4. **会话管理**: 使用 withCredentials 保持 Cookie 会话

## 相关文档

- [Vite 配置文档](https://vite.dev/config/)
- [Vue 3 官方文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [ECharts 文档](https://echarts.apache.org/)
