# 船舶企业员工培训管理系统 - 后端

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM ORM
- GreatSQL (MySQL兼容)
- Docker

## 项目结构

```
backend/
├── main.go              # 程序入口，路由定义
├── database/            # 数据库相关
│   ├── db.go           # 数据库连接
│   └── models.go       # 数据模型定义
├── handlers/            # 请求处理器
│   ├── auth/           # 认证相关接口
│   ├── home/           # 主页接口
│   ├── teacher/        # 讲师端接口
│   ├── employee/       # 员工端接口
│   └── planner/        # 课程大纲制定者接口
├── middleware/          # 中间件
│   ├── auth.go         # 简单鉴权中间件
│   └── cors.go         # CORS中间件
├── config/              # 配置
│   └── config.go       # 配置加载
├── docker-compose.yml   # Docker配置
├── .env                 # 环境变量
└── go.mod              # Go模块依赖
```

## 快速开始

### 1. 启动数据库

```bash
docker-compose up -d
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## 接口文档

详见 `/docs/后端接口设计文档.md`

## 接口概览

- 认证相关：4个接口
- 主页：1个接口
- 讲师端：8个接口
- 员工端：7个接口
- 课程大纲制定者：18个接口

共计 **37个接口**
