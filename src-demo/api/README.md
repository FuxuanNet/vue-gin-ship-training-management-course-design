# 前端网络请求接口封装说明

## 概述

本文档介绍了前端网络请求接口的封装方式，包括请求拦截器、响应拦截器以及API的使用方法。

## 目录结构

```txt
api/
├── request.js           # 请求封装，包含拦截器配置
├── auth.js              # 认证相关API（登录、注册等）
└── README.md            # 说明文档
```

## 请求封装（request.js）

`request.js` 使用 Axios 库封装了HTTP请求，实现了以下功能：

1. **基础配置**
   - 设置API基础URL
   - 设置请求超时时间

2. **请求拦截器**
   - 自动添加认证令牌到请求头
   - 处理请求错误

3. **响应拦截器**
   - 统一处理响应数据格式
   - 处理特定错误码（如401未授权）
   - 自动重定向到登录页面（当认证失败时）

## 使用方法

### 1. 导入API函数

```javascript
import { login, register, logout } from '@/api/auth'
```

### 2. 调用API函数

**示例：用户登录**:

```javascript
// 登录请求
login({
  username: 'user123',
  password: 'password123'
})
  .then(response => {
    // 登录成功处理
    const { token, user } = response.data
    
    // 存储令牌和用户信息
    localStorage.setItem('token', token)
    localStorage.setItem('user', JSON.stringify(user))
    
    // 跳转到首页
    this.$router.push('/')
  })
  .catch(error => {
    // 错误处理
    console.error('登录失败:', error)
    this.$message.error(error.message || '登录失败')
  })
```

**示例：用户注册**:

```javascript
// 注册请求
register({
  username: 'newuser',
  password: 'password123',
  confirmPassword: 'password123',
  email: 'user@example.com',
  company: '示例公司',
  role: 'data_provider'
})
  .then(response => {
    // 注册成功处理
    this.$message.success('注册成功，请登录')
    this.$router.push('/login')
  })
  .catch(error => {
    // 错误处理
    console.error('注册失败:', error)
    this.$message.error(error.message || '注册失败')
  })
```

**示例：获取当前用户信息**：

```javascript
// 获取当前用户信息
getCurrentUser()
  .then(response => {
    // 成功获取用户信息
    const userData = response.data
    this.userInfo = userData
  })
  .catch(error => {
    // 错误处理
    console.error('获取用户信息失败:', error)
  })
```

**示例：修改密码**：

```javascript
// 修改密码
updatePassword(data)
  .then(response => {
    this.$message.success('密码修改成功')
  })
  .catch(error => {
    this.$message.error(error.message || '修改密码失败')
  })
```

**示例：上传头像**：

```javascript
import { uploadAvatar } from '@/api/auth'

// 选择文件后上传
const file = event.target.files[0]
uploadAvatar(file)
  .then(response => {
    this.userInfo.avatar = response.data.avatar
    this.$message.success('头像上传成功')
  })
  .catch(error => {
    this.$message.error(error.message || '头像上传失败')
  })
```

**示例：获取头像URL**：

```javascript
import { getAvatarUrl } from '@/api/auth'

// 获取用户头像URL
const avatarUrl = getAvatarUrl(userId)
// 在模板中使用
<img :src="avatarUrl" />
```

## 后端响应格式

所有后端API响应均遵循以下格式：

```json
{
  "code": 200,           // 状态码，200表示成功
  "message": "成功",     // 响应消息
  "data": {              // 响应数据（可选）
    // 具体数据内容
  }
}
```

## 错误处理

前端API封装会自动处理以下错误情况：

1. **未授权访问（401）**
   - 清除本地存储的token和用户信息
   - 自动重定向到登录页面

2. **请求超时**
   - 返回超时错误信息

3. **服务器错误**
   - 返回服务器错误信息

## 认证机制

1. **Token存储**
   - 登录成功后，token存储在localStorage中
   - 保存路径: `localStorage.getItem('token')`

2. **Token使用**
   - 每次请求自动将token添加到请求头
   - 格式: `Authorization: Bearer {token}`

3. **Token失效**
   - 当接收到401响应时，自动清除token并重定向到登录页

## 数据资源API

### 上传数据资源

用于上传数据资源文件和相关元数据。

```javascript
import { uploadDataResource } from '@/api/resources'

// 表单数据
const dataForm = {
  name: '金融交易数据集',
  category: 'finance',
  dataStructure: 'CSV格式，包含交易ID、金额、时间等字段',
  isConfidential: false,
  needHosting: false,
  description: '# 金融交易数据集\n\n## 数据概述\n这是一个金融交易数据集...',
  tags: ['金融', '交易', 'CSV'],
  fileExtensions: ['csv', 'xlsx']
}

// 文件对象
const dataFile = file // 从文件选择器获取的文件对象
const sampleFile = sampleFile // 可选的样例文件

// 调用上传API
try {
  const result = await uploadDataResource(dataForm, dataFile, sampleFile)
  console.log('上传成功:', result)
} catch (error) {
  console.error('上传失败:', error)
}
```

### 获取用户资源列表

获取当前用户上传的所有数据资源。

```javascript
import { getUserResources } from '@/api/resources'

try {
  const result = await getUserResources()
  console.log('用户资源列表:', result.data)
} catch (error) {
  console.error('获取失败:', error)
}
```

### 获取资源详情

获取指定ID的数据资源详情。

```javascript
import { getResourceDetail } from '@/api/resources'

const resourceId = 1 // 资源ID

try {
  const result = await getResourceDetail(resourceId)
  console.log('资源详情:', result.data)
} catch (error) {
  console.error('获取失败:', error)
}
```

### 获取资源分类列表

获取系统支持的资源分类列表。

```javascript
import { getResourceCategories, getMockCategories } from '@/api/resources'

// 正式环境
try {
  const result = await getResourceCategories()
  console.log('资源分类:', result.data)
} catch (error) {
  console.error('获取失败:', error)
}

// 开发环境（使用模拟数据）
try {
  const result = await getMockCategories()
  console.log('资源分类:', result.data)
} catch (error) {
  console.error('获取失败:', error)
}
```
