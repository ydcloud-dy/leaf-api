# 博客前端项目

基于 Vue 3 + Element Plus 构建的现代化博客前台系统。

## 技术栈

- **框架**: Vue 3 (Composition API)
- **路由**: Vue Router 4
- **状态管理**: Pinia
- **UI 组件库**: Element Plus
- **HTTP 客户端**: Axios
- **Markdown 渲染**: markdown-it
- **日期处理**: dayjs
- **构建工具**: Vite

## 项目结构

```
blog-frontend/
├── package.json          # 项目依赖配置
├── vite.config.js        # Vite 配置文件
├── index.html            # HTML 入口文件
└── src/
    ├── main.js           # 应用入口
    ├── App.vue           # 根组件
    ├── router/           # 路由配置
    │   └── index.js
    ├── stores/           # Pinia 状态管理
    │   └── user.js       # 用户状态
    ├── api/              # API 接口封装
    │   ├── request.js    # Axios 封装
    │   ├── auth.js       # 认证相关接口
    │   ├── article.js    # 文章相关接口
    │   └── comment.js    # 评论相关接口
    ├── layouts/          # 布局组件
    │   └── MainLayout.vue
    ├── views/            # 页面组件
    │   ├── Home.vue          # 首页
    │   ├── Articles.vue      # 文章列表
    │   ├── ArticleDetail.vue # 文章详情
    │   ├── Archive.vue       # 归档页面
    │   ├── About.vue         # 关于页面
    │   ├── Guestbook.vue     # 留言板
    │   ├── Login.vue         # 登录/注册
    │   └── Profile.vue       # 个人中心
    └── components/       # 公共组件
        ├── ArticleCard.vue   # 文章卡片
        ├── Comment.vue       # 评论组件
        └── Header.vue        # 导航栏
```

## 功能特性

### 已实现功能

1. **用户系统**
   - 用户注册/登录
   - 个人中心
   - 账号设置
   - 修改密码
   - 权限控制

2. **文章系统**
   - 文章列表（分页）
   - 文章详情（Markdown 渲染）
   - 文章搜索
   - 分类筛选
   - 排序功能
   - 文章归档

3. **互动功能**
   - 文章点赞/取消点赞
   - 文章收藏/取消收藏
   - 发表评论
   - 回复评论
   - 评论点赞
   - 留言板

4. **其他功能**
   - 响应式布局
   - 路由守卫
   - 用户状态持久化
   - 统一错误处理
   - 热门文章
   - 标签云
   - 站点统计

## 安装依赖

```bash
cd blog-frontend
npm install
```

## 开发模式

```bash
npm run dev
```

项目将在 http://localhost:3000 启动

## 生产构建

```bash
npm run build
```

构建产物将生成在 `dist` 目录

## 预览生产构建

```bash
npm run preview
```

## API 配置

项目默认后端地址为 `http://localhost:8888`，如需修改，请编辑 `vite.config.js` 中的 proxy 配置：

```javascript
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8888',  // 修改为你的后端地址
      changeOrigin: true,
      rewrite: (path) => path.replace(/^\/api/, '')
    }
  }
}
```

## 环境要求

- Node.js >= 16
- npm >= 8

## 主要依赖版本

- vue: ^3.4.0
- vue-router: ^4.2.5
- pinia: ^2.1.7
- element-plus: ^2.5.0
- axios: ^1.6.2
- markdown-it: ^14.0.0
- dayjs: ^1.11.10

## 开发说明

### 路由配置

路由配置在 `src/router/index.js` 中，使用嵌套路由结构，主要路由包括：

- `/` - 首页
- `/articles` - 文章列表
- `/articles/:id` - 文章详情
- `/archive` - 归档
- `/guestbook` - 留言板
- `/about` - 关于
- `/login` - 登录/注册
- `/profile` - 个人中心（需登录）

### 状态管理

使用 Pinia 管理全局状态，当前只有用户状态（user store）：

- 用户信息
- 登录状态
- Token 管理
- 登录/注册/登出方法

### API 调用

所有 API 调用都通过 `src/api` 目录中的模块进行：

```javascript
import { getArticles } from '@/api/article'

const { data } = await getArticles({ page: 1, page_size: 10 })
```

### 样式规范

- 使用 Element Plus 主题色
- 响应式断点：768px
- 采用 scoped 样式避免污染

## 后端接口要求

项目需要以下后端接口支持：

### 认证接口
- POST /auth/login - 用户登录
- POST /auth/register - 用户注册
- GET /auth/me - 获取当前用户信息
- PUT /auth/profile - 更新用户信息
- PUT /auth/password - 修改密码

### 文章接口
- GET /articles - 获取文章列表
- GET /articles/:id - 获取文章详情
- GET /articles/search - 搜索文章
- GET /articles/archive - 获取归档文章
- POST /articles/:id/like - 点赞文章
- DELETE /articles/:id/like - 取消点赞
- POST /articles/:id/favorite - 收藏文章
- DELETE /articles/:id/favorite - 取消收藏

### 评论接口
- GET /comments - 获取评论列表
- GET /articles/:id/comments - 获取文章评论
- POST /comments - 发表评论
- POST /comments/:id/reply - 回复评论
- DELETE /comments/:id - 删除评论
- POST /comments/:id/like - 点赞评论

### 用户数据接口
- GET /user/likes - 获取我的点赞
- GET /user/favorites - 获取我的收藏

## 浏览器支持

- Chrome >= 90
- Firefox >= 88
- Safari >= 14
- Edge >= 90

## 注意事项

1. 确保后端服务已启动在 http://localhost:8888
2. 首次运行需要安装依赖：`npm install`
3. 如需修改端口，请编辑 `vite.config.js` 中的 `server.port`
4. Token 存储在 localStorage 中
5. 图片资源需要配置正确的 CDN 地址或本地路径

## 许可证

MIT
