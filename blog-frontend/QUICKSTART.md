# 快速启动指南

## 前置要求

确保你的系统已安装：
- Node.js (版本 >= 16)
- npm (版本 >= 8)

## 启动步骤

### 1. 进入项目目录

```bash
cd /Users/dujie/go/src/github.com/ydcloud-dy/leaf-api/blog-frontend
```

### 2. 安装依赖

```bash
npm install
```

这将安装所有必需的依赖包，包括：
- Vue 3
- Vue Router
- Pinia
- Element Plus
- Axios
- Markdown-it
- Dayjs

### 3. 启动开发服务器

```bash
npm run dev
```

启动成功后，访问：http://localhost:3000

### 4. 确保后端服务运行

项目需要后端服务支持，确保后端服务运行在：http://localhost:8888

如果后端地址不同，请修改 `vite.config.js` 中的 proxy 配置。

## 常见问题

### 1. 端口被占用

如果 3000 端口已被占用，可以修改 `vite.config.js` 中的端口：

```javascript
server: {
  port: 3001, // 修改为其他端口
  // ...
}
```

### 2. 后端连接失败

检查：
1. 后端服务是否已启动
2. 后端地址是否正确（默认 http://localhost:8888）
3. 防火墙是否阻止了连接

### 3. 依赖安装失败

尝试：
```bash
# 清除缓存
npm cache clean --force

# 删除 node_modules 和 package-lock.json
rm -rf node_modules package-lock.json

# 重新安装
npm install
```

## 构建生产版本

```bash
npm run build
```

构建完成后，产物在 `dist` 目录。

## 预览生产构建

```bash
npm run preview
```

## 主要功能路由

- `/` - 首页（文章列表）
- `/articles` - 文章列表（带搜索和筛选）
- `/articles/:id` - 文章详情
- `/archive` - 归档页面
- `/guestbook` - 留言板
- `/about` - 关于页面
- `/login` - 登录/注册
- `/profile` - 个人中心（需登录）

## 测试账号

根据后端配置创建测试账号，或使用注册功能创建新账号。

## 技术支持

如有问题，请查看：
1. README.md - 详细文档
2. 浏览器控制台 - 查看错误信息
3. 网络面板 - 检查 API 请求
