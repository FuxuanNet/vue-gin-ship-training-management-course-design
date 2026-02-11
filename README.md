# 计算机软件综合课程设计——船舶企业培训管理系统

## 前言

本项目是计算机软件综合设计/数据库课程设计/信息管理系统实现课程的课程设计大作业，可能存在很多设计与代码不合理不规范之处，也会有很多 bug，请不要用于生产环境。

## 1. 课题题目

船舶企业培训管理系统（船舶人力资源管理系统——培训与发展）。

## 2. 开发目的及意义

当前，船舶人力资源管理领域存在着以下问题：

- 培训计划的制定和管理混乱；
- 员工难以掌握课程时间和安排；
- 效果评估方式不规范，无法掌握员工技能掌握情况。

针对以上现状，本系统围绕着船舶人力资源管理系统——培训与发展的选题进行拓展完善，实现了船舶企业培训管理系统。

该系统可以保证计划的指定与执行，培训课程与时间的管理以及培训的效果评估，致力于解决上述现状。

## 3．系统需求分析

船舶企业培训管理系统支持员工培训计划的制定和执行，包括培训课程、培训时间、培训效果评估等，帮助员工提升技能水平。

### 3.1 需求分析

#### 目标用户

本系统面向参训员工，课程大纲指定负责人与培训讲师，制定三种用户角色，分别是：员工、课程大纲制定者（管理员）、培训讲师。

#### 功能概述

- 针对课程大纲制定者，本系统可以进行培训计划管理与课程管理，以及培训大纲效果评估。
- 针对讲师，本系统可以提供授课安排提醒，打分与教学效果评估。
- 针对员工，本系统可以提供上课安排提醒，打分与上课表现评估。

详细的功能性需求的描述如图 3.1 所示。

![图 3.1 功能性需求](docs/image/README/1766817219584.png)

#### 业务流程

三大用户使用本系统的大体流程如下图 3.2 所示。

![图 3.2 系统业务流程泳道图](docs/image/README/1766817284013.png)

### 3.2 开发语言及工具

- 前端

  - Vue3 框架 + Element Plus 组件库 + Axios 库 + ECharts 库

- 后端

  - Gin 框架 +  Gorm 库

- 数据库

  - GreatSQL（国产开源数据库）+ Docker

- 工具

  - 开发工具：VSCode

  - 开发环境：Windows 11

- 模块架构设计

![图 3.3 系统模块架构图](docs/image/README/1766817402261.png)

### 5.1 系统功能划分

本节将简要介绍系统的大致功能划分与代码实现，具体功能实现内容和效果将在5.2 系统实现界面结合展示。

本系统所提供的主要功能可分为以下四大类：培训计划管理、课程管理、时间安排提醒、效果评估与可视化。

培训计划管理包括：

- 创建、删除、修改、查询培训计划以及其详细信息。
- 在培训计划中加入课程信息，并指定上课的时间地点等信息。
- 在培训计划中添加员工名单。

课程管理包括：

- 创建、删除、修改、查询培训课程，并选定授课讲师、上课信息等详细信息。

时间安排提醒：

- 针对员工和讲师，展示上课课程表和授课课程表。

效果评估：

- 一门课程的成绩组成包括员工自评和讲师评价，共同组成该员工在该课程的成绩。讲师可以设置分数权重。员工自评和讲师评价内容都可以进行 AI 评价生成评价结果，同时讲师也可以直接打分，最终形成单个员工的加权评价体系，进而生成课程加权评价体系和培训计划的加权评价体系。

### 5.2 系统实现界面

#### 登录端

功能一：员工、讲师、课程大纲制定者可以通过登录界面进入系统。

功能二：提供了用于快速测试的模拟账号，涵盖上述角色。

![1766818388935](docs/image/README/1766818388935.png)

#### 员工端

功能三：登录后展示员工界面首页。教师端与课程大纲制定者端也具有同样的功能，仅仅是导航栏内容不同，故后续不再重复展示，后续其余类似功能也是如此，仅在员工端详细展示。

功能四：导航栏展示，展示课程表、课程自评和我的成绩页面的导航，用户可以点击跳转对应链接。

功能五：显示个人信息与退出账号。

![1766818398642](docs/image/README/1766818398642.png)

功能六：首页快速导航与系统特殊介绍。

![1766818407156](docs/image/README/1766818407156.png)

功能七：显示员工上课课程表，展示对应上课时间、地点与周次。

![1766818414179](docs/image/README/1766818414179.png)

功能八：员工可以查看待评价的课程自评。

![1766818418830](docs/image/README/1766818418830.png)

功能九：员工可以进行员工评价并提交，生成对应评价分数。

![1766818425075](docs/image/README/1766818425075.png)

功能九：展示员工的自评成绩、讲师评分和综合成绩。

![1766818438037](docs/image/README/1766818438037.png)

功能十：展示课程类型掌握度雷达图和各类型平均分。

![1766818443892](docs/image/README/1766818443892.png)

#### 讲师端

功能十一：显示讲师的授课表，包括上课时间地点与周次等。

![1766818453990](docs/image/README/1766818453990.png)

功能十二：查看对学员的评分信息和自评内容。

![1766818459595](docs/image/README/1766818459595.png)

功能十三：对学员进行打分与修改。可以直接评分，也可以进行评语，让 AI 生成分数。

![1766818488940](docs/image/README/1766818488940.png)

功能十四：授课课程的学员成绩统计分析可视化与课次统计。

![1766818501864](docs/image/README/1766818501864.png)

#### 课程大纲制定者端

功能十五：查看当前已创建的培训计划。

![1766818511983](docs/image/README/1766818511983.png)

功能十六：创建新的培训计划。

![1766818517519](docs/image/README/1766818517519.png)

功能十七：展示培训计划的详细信息。

![1766818526276](docs/image/README/1766818526276.png)

功能十八：在培训计划里添加员工。

![1766818534158](docs/image/README/1766818534158.png)

功能十九：在培训计划里添加课程安排信息，包括课程、课程日期、开始时间、结束时间和上课地点。

![1766818541773](docs/image/README/1766818541773.png)

功能二十：编辑培训计划。

![1766818549586](docs/image/README/1766818549586.png)

功能二十一：删除培训计划。

![1766818555206](docs/image/README/1766818555206.png)

如果计划下存在课程安排或关联员工，则会触发错误提示：

![1766818565157](docs/image/README/1766818565157.png)

功能二十二：查看当前开设的全部课程。

![1766818570975](docs/image/README/1766818570975.png)

功能二十三：新建课程，填写课程名称、描述、要求、类型和主讲讲师等信息。

![1766818577025](docs/image/README/1766818577025.png)

功能二十四：编辑修改课程相关信息。

![1766818584131](docs/image/README/1766818584131.png)

功能二十五：删除课程相关信息。

![1766818594284](docs/image/README/1766818594284.png)

功能二十六：全局数据分析，包括课程评分排名，培训计划评分排名，课程类型分布和培训计划状态统计等可视化。

![1766818607482](docs/image/README/1766818607482.png)


## 4. 数据库设计

### 4.1 概念结构设计

- E-R 图

本系统的 E-R 图如图 4.1 所示。

![图 4.1 数据库设计 E-R 图](docs/image/README/1766817479739.png)

包含的实体有：培训计划、培训课程安排、培训课程、人员、账号。

包含的关系有：

- **拥有**

人员 拥有 账号：1 对 1

- **包含**

培训计划 包含 培训课程安排：1 对 n

- **安排**

培训课程安排 安排 培训课程：n 对 1

- **参与和评价**

人员（角色="员工"） 参与和评价 培训课程安排：m 对 n

- 属性：
  - 所属安排 ID
  - 所属人员 ID（角色="员工"）
  - 自评分数
  - 自评内容
  - 讲师评分
  - 讲师评语
  - 讲师评分构成占比

- **制定**

人员（角色="课程大纲指定者"） 指定 培训计划：1 对 n

- **规划**

培训计划 规划 人员（角色="员工"）：m 对 n

- **讲授**

人员（角色="讲师"） 讲授 培训课程：1 对 n

### 4.2 逻辑结构设计

#### 培训计划表：training_plan

| 字段名              | 数据类型    | 主/外键 | 备注                               |
| ------------------- | ----------- | ------- | ---------------------------------- |
| plan_id             | BIGINT      | 主键    | 计划 ID                            |
| plan_name           | VARCHAR(50) | -       | 计划名称                           |
| plan_status         | CHAR(3)     | -       | 计划状态：规划中/进行中/已完成     |
| plan_start_datetime | varchar(8)  | -       | 计划开始时间                       |
| plan_end_datetime   | varchar(8)  | -       | 计划结束时间                       |
| creator_id          | BIGINT      | 外键    | 制定人（人员ID，表示“制定”关系） |

#### 规划员工表：plan_employee

| 字段名                                 | 数据类型 | 主/外键  | 备注    |
| -------------------------------------- | -------- | -------- | ------- |
| plan_id                                | BIGINT   | 主键外键 | 计划 ID |
| person_id                              | BIGINT   | 主键外键 | 人员 ID |

#### 培训课程表：course

| 字段名         | 数据类型     | 主/外键 | 备注                             |
| -------------- | ------------ | ------- | -------------------------------- |
| course_id      | BIGINT       | 主键    | 课程 ID                          |
| course_name    | VARCHAR(50)  | -       | 课程名称                         |
| course_desc    | VARCHAR(100) | -       | 课程描述                         |
| course_require | VARCHAR(500) | -       | 课程要求                         |
| course_class   | VARCHAR(20)  | -       | 课程类型                         |
| teacher_id     | BIGINT       | 外键    | 主讲讲师（人员ID，“讲授”关系） |

#### 人员表：person

| 字段名    | 数据类型    | 主/外键 | 备注                           |
| --------- | ----------- | ------- | ------------------------------ |
| person_id | BIGINT      | 主键    | 人员 ID                        |
| name      | VARCHAR(20) | -       | 姓名                           |
| role      | VARCHAR(7)  | -       | 角色：员工/讲师/课程大纲制定者 |

#### 账号表：account

| 字段名        | 数据类型     | 主/外键 | 备注     |
| ------------- | ------------ | ------- | -------- |
| account_id    | BIGINT       | 主键，外键 | 账号 ID  |
| person_id     | BIGINT       | 主键    | 人员 ID  |
| login_name    | VARCHAR(20)  | -       | 用户名   |
| password_hash | VARCHAR(255) | -       | 密码哈希 |

#### 培训课程安排表：plan_course_item

| 字段名           | 数据类型     | 主/外键 | 备注                      |
| ---------------- | ------------ | ------- | ------------------------- |
| item_id          | BIGINT       | 主键    | 安排 ID                   |
| plan_id          | BIGINT       | 外键    | 所属计划 ID（"包含"关系） |
| course_id        | BIGINT       | 外键    | 对应课程 ID（"安排"关系） |
| class_date       | DATE         | -       | 上课日期                  |
| class_begin_time | TIME         | -       | 课程开始时间              |
| class_end_time   | TIME         | -       | 课程结束时间              |
| location         | VARCHAR(100) | -       | 上课地点                  |

#### 参与和评价表：`attendance_evaluation`

| 字段名          | 数据类型 | 主/外键    | 备注                                     |
| --------------- | -------- | ---------- | ---------------------------------------- |
| person_id       | BIGINT   | 主键，外键 | 人员ID（员工，连到【人员】实体）         |
| item_id         | BIGINT   | 主键，外键 | 培训课程安排ID(连到【培训课程安排】实体) |
| self_score      | FLOAT(2) | -          | 员工自评分（算法生成）                   |
| self_comment    | TEXT     | -          | 自评文本/问卷内容                        |
| teacher_score   | FLOAT(2) | -          | 讲师给该员工的评分                       |
| teacher_comment | TEXT     | -          | 讲师评语                                 |
| score_ratio     | FLOAT(2) | -          | 讲师评分构成占比，由此可计算出员工自评   |

视图设计：

本系统主要设计了员工每节课得分视图、员工×课程类型综合得分视图、课程整体评分视图、培训计划整体评分视图、员工课表视图、讲师课表视图六类视图。由于内容过长，将移至**附录**展示。

### 4.3 物理结构设计

本系统使用 GreatSQL 数据库，使用 Docker 进行安装。GreatSQL是由开放原子开源基金会孵化的金融级开源数据库项目，全面支持ARM、x86、loongArch、SW-64等多样性计算架构[1]，是国产开源的数据库管理系统。

Docker 是一个开源的应用容器引擎，基于 Go 语言并遵从 Apache2.0 协议开源。Docker 可以让开发者打包他们的应用以及依赖包到一个轻量级、可移植的容器中，然后发布到任何流行的 Linux 机器上，也可以实现虚拟化[1]。因此，本系统使用 Docker 来安装 GreatSQL。

GORM 通过将 Go 结构体（Go structs） 映射到数据库表来简化数据库交互。避免直接操控 SQL 语句，提高代码的可维护性，并保证跨数据库兼容性[2]。

#### 索引/触发器设计

本系统所有表均采用主键索引，采用 BIGINT 类型自增主键，并在外键创建普通索引，以提高查询性能。本系统暂时没有使用触发器的需求。

### 5. 系统详细设计和实现

#### 代码文件结构

项目主要结构分为前端文件夹，后端文件夹，文档文件夹，Git 文件夹。

前端的项目结构如下所示：

```txt
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
│   │   │   └── Scores.vue        # 成绩查看/能力分析
│   │   │
│   │   ├── planner/          # 课程大纲制定者页面
│   │   │   ├── Analytics.vue     # 数据分析/平台统计
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

后端的项目结构如下所示：

```txt
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

文档文件夹主要包含需求分析文档、界面设计的一些想法的草稿、数据库设计文档还有接口设计文档等。

#### 内部接口设计

##### 一、认证相关接口

| 接口名称             | 接口路径                   | 请求方式 | 功能描述                                                         |
| -------------------- | -------------------------- | -------- | ---------------------------------------------------------------- |
| 用户登录接口         | `/api/auth/login`        | POST     | 前端提交用户名和密码，后端验证后返回JWT Token和用户基本信息      |
| 用户注册接口         | `/api/auth/register`     | POST     | 前端提交注册信息，后端验证数据后创建用户账号和个人信息           |
| 退出登录接口         | `/api/auth/logout`       | POST     | 前端提交Token，后端验证后将Token加入黑名单（如果实现黑名单机制） |
| 获取当前用户信息接口 | `/api/auth/current-user` | GET      | 前端提交Token，后端解析并验证Token后返回用户详细信息             |

##### 二、主页相关接口

| 接口名称             | 接口路径                 | 请求方式 | 功能描述                                                           |
| -------------------- | ------------------------ | -------- | ------------------------------------------------------------------ |
| 获取平台统计数据接口 | `/api/home/statistics` | GET      | 前端请求平台整体统计数据，后端查询并统计后返回全局或个性化统计数据 |

##### 三、讲师端接口

| 接口名称               | 接口路径                             | 请求方式 | 功能描述                                                                           |
| ---------------------- | ------------------------------------ | -------- | ---------------------------------------------------------------------------------- |
| 获取讲师授课表接口     | `/api/teacher/schedule`            | GET      | 前端请求指定时间范围内的授课表，后端验证权限后返回讲师的授课安排                   |
| 获取待评分学员列表接口 | `/api/teacher/pending-evaluations` | GET      | 前端请求待评分的学员列表，后端验证权限后返回讲师未评分的学员记录                   |
| 提交学员评分接口       | `/api/teacher/submit-evaluation`   | POST     | 前端提交评分信息，后端验证权限后更新学员的讲师评分和评语                           |
| 批量提交评分接口       | `/api/teacher/batch-evaluation`    | POST     | 前端提交多个学员的评分信息数组，后端验证权限后批量更新学员评分                     |
| 设置评分占比接口       | `/api/teacher/score-ratio`         | PUT      | 前端提交课程安排ID和新的评分占比，后端验证权限后更新该课程安排下所有学员的评分占比 |
| 获取课程成绩统计接口   | `/api/teacher/course-statistics`   | GET      | 前端请求指定课程的成绩统计信息，后端验证权限后返回课程的成绩统计数据               |
| 获取讲师授课统计接口   | `/api/teacher/teaching-statistics` | GET      | 前端请求讲师的整体授课统计信息，后端验证权限后返回讲师的授课统计数据               |

##### 四、员工端接口

| 接口名称                 | 接口路径                              | 请求方式 | 功能描述                                                               |
| ------------------------ | ------------------------------------- | -------- | ---------------------------------------------------------------------- |
| 获取员工课程表接口       | `/api/employee/schedule`            | GET      | 前端请求指定时间范围内的课程表，后端验证权限后返回员工的课程安排       |
| 获取待自评课程列表接口   | `/api/employee/pending-evaluations` | GET      | 前端请求待自评的课程列表，后端验证权限后返回员工未自评的课程记录       |
| 提交课程自评接口         | `/api/employee/submit-evaluation`   | POST     | 前端提交自评信息，后端验证权限后将自评内容发送给AI接口生成分数并存储   |
| 获取员工成绩列表接口     | `/api/employee/scores`              | GET      | 前端请求成绩列表，后端验证权限后返回员工的所有课程得分                 |
| 获取课程类型成绩分析接口 | `/api/employee/course-type-scores`  | GET      | 前端请求课程类型维度的成绩分析，后端验证权限后返回每个课程类型的平均分 |
| 获取员工学习进度接口     | `/api/employee/learning-progress`   | GET      | 前端请求学习进度信息，后端验证权限后返回员工的学习进度和统计数据       |

##### 五、课程大纲制定者端接口

| 接口名称               | 接口路径                                             | 请求方式 | 功能描述                                                         |
| ---------------------- | ---------------------------------------------------- | -------- | ---------------------------------------------------------------- |
| 获取培训计划列表接口   | `/api/planner/plans`                               | GET      | 前端请求培训计划列表，后端验证权限后返回所有培训计划及统计信息   |
| 创建培训计划接口       | `/api/planner/plans`                               | POST     | 前端提交培训计划信息，后端验证权限后创建新的培训计划             |
| 修改培训计划接口       | `/api/planner/plans/:planId`                       | PUT      | 前端提交修改的培训计划信息，后端验证权限后更新培训计划信息       |
| 删除培训计划接口       | `/api/planner/plans/:planId`                       | DELETE   | 前端提交要删除的计划ID，后端验证权限后删除培训计划               |
| 获取培训计划详情接口   | `/api/planner/plans/:planId`                       | GET      | 前端请求指定计划的详细信息，后端验证权限后返回培训计划的完整信息 |
| 为培训计划添加员工接口 | `/api/planner/plans/:planId/employees`             | POST     | 前端提交计划ID和员工ID列表，后端验证权限后为培训计划添加员工     |
| 从培训计划移除员工接口 | `/api/planner/plans/:planId/employees/:employeeId` | DELETE   | 前端提交计划ID和员工ID，后端验证权限后从培训计划移除员工         |
| 获取课程列表接口       | `/api/planner/courses`                             | GET      | 前端请求课程列表，后端验证权限后返回所有课程及统计信息           |
| 创建课程接口           | `/api/planner/courses`                             | POST     | 前端提交课程信息，后端验证权限后创建新的课程                     |
| 修改课程接口           | `/api/planner/courses/:courseId`                   | PUT      | 前端提交修改的课程信息，后端验证权限后更新课程信息               |
| 删除课程接口           | `/api/planner/courses/:courseId`                   | DELETE   | 前端提交要删除的课程ID，后端验证权限后删除课程                   |
| 获取课程安排列表接口   | `/api/planner/course-items`                        | GET      | 前端请求课程安排列表，后端验证权限后返回所有课程安排及统计信息   |
| 创建课程安排接口       | `/api/planner/course-items`                        | POST     | 前端提交课程安排信息，后端验证权限后创建新的课程安排             |
| 修改课程安排接口       | `/api/planner/course-items/:itemId`                | PUT      | 前端提交修改的课程安排信息，后端验证权限后更新课程安排信息       |
| 删除课程安排接口       | `/api/planner/course-items/:itemId`                | DELETE   | 前端提交要删除的课程安排ID，后端验证权限后删除课程安排           |
| 获取平台数据分析接口   | `/api/planner/analytics`                           | GET      | 前端请求平台整体数据分析，后端验证权限后返回综合数据分析结果     |
| 获取员工成绩详情接口   | `/api/planner/employees/:employeeId/scores`        | GET      | 前端请求指定员工的成绩详情，后端验证权限后返回员工的成绩完整信息 |
| 获取课程评价详情接口   | `/api/planner/courses/:courseId/evaluations`       | GET      | 前端请求指定课程的评价详情，后端验证权限后返回课程的评价完整信息 |

#### 外部接口设计

本系统的外部接口是系统与 DeepSeek API 的接口，用于给员工自评和教师评价进行打分，接收 AI 的自动打分结果。

### 6. 系统测试

由于工程量较大，本项目仅针对部分核心主要功能进行功能验证，如表 6.1 所示。

![表 6.1 测试项与测试用例表格](docs/image/README/1766818205199.png)

## 结论

本项目成功设计实现了船舶企业培训管理系统，该系统围绕船舶人力资源管理领域的痛点问题，构建了一个功能完整、架构清晰的培训管理解决方案，围绕需求文档、数据库设计文档和系统接口设计文档，实现了四大核心功能，解决了船舶企业培训管理中存在的计划制定混乱、课程安排不清晰和效果评估不规范等问题，提高了培训管理效率。

## 参考文献

[1]GreatSQL 用户手册. 容器化安装（Docker）[EB/OL]. (2024-10-9)（引用日期: 2025-12-26）. [https://greatsql.cn/docs/8.0.32-26/3-quick-start/3-3-quick-start-with-docker.html](https://greatsql.cn/docs/8.0.32-26/3-quick-start/3-3-quick-start-with-docker.html)

[2]GORM.GORM 中文指南[EB/OL]. (引用日期: 2025-12-26). [https://gorm.io/zh_CN/docs/index.html](https://gorm.io/zh_CN/docs/index.html)

## 附录

### 员工每节课得分视图：`v_employee_item_score`

**功能说明：**

按“员工 × 课程安排”展示每一节课的详细信息和得分情况，用于：

- 员工个人查看每节课得分情况
- 后续做课程、计划、员工维度的各类统计的基础视图

class_time 来源已发生更改，后续将更改。

course_class 来源有可能也发生更改。

| 字段名           | 来源                                | 备注            |
| ---------------- | ----------------------------------- | --------------- |
| person_id        | attendance_evaluation.person_id     | 员工 ID         |
| person_name      | person.name                         | 员工姓名        |
| item_id          | plan_course_item.item_id            | 课程安排 ID     |
| class_date       | plan_course_item.class_date         | 上课日期        |
| class_begin_time | plan_course_item.class_begin_time   | 上课开始时间    |
| class_end_time   | plan_course_item.class_end_time     | 上课结束时间    |
| location         | plan_course_item.location           | 上课地点        |
| plan_id          | plan_course_item.plan_id            | 所属培训计划 ID |
| plan_name        | training_plan.plan_name             | 培训计划名称    |
| course_id        | plan_course_item.course_id          | 课程 ID         |
| course_name      | course.course_name                  | 课程名称        |
| course_class     | course.course_class                 | 课程类型        |
| self_score       | attendance_evaluation.self_score    | 员工自评分      |
| teacher_score    | attendance_evaluation.teacher_score | 讲师评分        |
| weighted_score   | 由 self_score 和 teacher_score 计算 | 加权得分        |

### 员工 × 课程类型综合得分视图：`v_employee_course_score`

**功能说明：**

按“员工 + 课程类型”维度，汇总该员工在该系列类型的课程下所有上课记录的平均加权得分，用于：

- 分析某员工在哪些课程上薄弱
- 后续做员工×课程类型的雷达图、热力图等可视化

| 字段名             | 来源                               | 备注                           |
| ------------------ | ---------------------------------- | ------------------------------ |
| person_id          | v_employee_item_score.person_id    | 员工 ID                        |
| course_id          | v_employee_item_score.course_class | 课程 ID                        |
| avg_weighted_score | 聚合计算                           | 员工在该类课程上的平均加权得分 |

### 课程整体评分视图：`v_course_score`

**功能说明：**

从“课程”维度汇总所有员工的表现，用于：

- 分析课程整体效果（平均得分）
- 看这门课参与人数，辅助判断课程难度、受欢迎程度等

| 字段名           | 来源                              | 备注                         |
| ---------------- | --------------------------------- | ---------------------------- |
| course_id        | v_employee_item_score.course_id   | 课程 ID                      |
| course_name      | v_employee_item_score.course_name | 课程名称                     |
| course_avg_score | 聚合计算                          | 课程整体平均加权得分         |
| student_count    | 聚合计算                          | 参与该课程的员工人数（去重） |

### 培训计划整体评分视图：`v_plan_score`

**功能说明：**

从“培训计划”维度，综合该计划下所有课程、所有员工的表现，用于：

- 对每个培训计划做整体效果评估
- 做不同培训计划之间的对比分析

| 字段名         | 来源                            | 备注                 |
| -------------- | ------------------------------- | -------------------- |
| plan_id        | v_employee_item_score.plan_id   | 培训计划 ID          |
| plan_name      | v_employee_item_score.plan_name | 培训计划名称         |
| plan_avg_score | 聚合计算                        | 培训计划整体平均得分 |

### 员工课表视图：`v_employee_timetable`

**功能说明：**

面向员工端，展示“员工个人课程表”，包括：

- 每节课的时间、地点、课程名称、课程要求等
- 用于首页“今日课程”、“本周课程”的课表展示

| 字段名           | 来源                              | 备注         |
| ---------------- | --------------------------------- | ------------ |
| person_id        | attendance_evaluation.person_id   | 员工 ID      |
| person_name      | person.name                       | 员工姓名     |
| item_id          | plan_course_item.item_id          | 课程安排 ID  |
| class_date       | plan_course_item.class_date       | 上课日期     |
| class_begin_time | plan_course_item.class_begin_time | 上课开始时间 |
| class_end_time   | plan_course_item.class_end_time   | 上课结束时间 |
| location         | plan_course_item.location         | 上课地点     |
| course_name      | course.course_name                | 课程名称     |
| course_require   | course.course_require             | 课程要求     |

### 讲师课表视图：`v_teacher_timetable`

**功能说明：**

面向讲师端，展示“讲师个人课程表”，包括：

- 每节课的时间、地点、课程名称、课程要求等
- 用于讲师查看自己每日/每周的授课安排

| 字段名           | 来源                                 | 备注         |
| ---------------- | ------------------------------------ | ------------ |
| teacher_id       | course.teacher_id / person.person_id | 讲师 ID      |
| teacher_name     | person.name                          | 讲师姓名     |
| item_id          | plan_course_item.item_id             | 课程安排 ID  |
| class_date       | plan_course_item.class_date          | 上课日期     |
| class_begin_time | plan_course_item.class_begin_time    | 上课开始时间 |
| class_end_time   | plan_course_item.class_end_time      | 上课结束时间 |
| location         | plan_course_item.location            | 上课地点     |
| course_id        | course.course_id                     | 课程 ID      |
| course_name      | course.course_name                   | 课程名称     |
| course_require   | course.course_require                | 课程要求     |
