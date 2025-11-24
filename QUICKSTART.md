# 🚀 快速启动指南

## ✅ 已完成的集成

### 后端
- ✅ Redis 配置和初始化
- ✅ 在线用户追踪服务（基于 Redis）
- ✅ 访问时长记录服务（基于 MySQL）
- ✅ 扩展统计 API（包含所有统计项）
- ✅ 数据库迁移（page_visits 表）
- ✅ 网站启动日期配置

### 前端
- ✅ 心跳服务（useHeartbeat composable）
- ✅ 访问追踪服务（useVisitTracking composable）
- ✅ App.vue 集成心跳和访问追踪
- ✅ 统计页面组件（Stats.vue）
- ✅ 路由配置
- ✅ 导航栏添加统计入口

## 📊 统计功能列表

| 统计项 | 数据来源 | 更新频率 |
|--------|---------|---------|
| 网站运行时长 | settings 表 | 实时计算 |
| 24小时访问量（PV） | views 表 | 实时 |
| 当前在线人数 | Redis | 实时（30秒） |
| 平均访问时长 | page_visits 表 | 实时 |
| 文章篇数 | articles 表 | 实时 |
| 笔记篇数 | chapters 表 | 实时 |
| 文章分类数 | categories 表 | 实时 |
| 文章标签数 | tags 表 | 实时 |
| 总浏览量 | articles.view_count | 实时 |
| 评论总数 | comments 表 | 实时 |
| 注册用户数 | users 表 | 实时 |

## 🎯 启动步骤

### 1. 确认 Redis 运行

```bash
redis-cli ping
# 应该返回: PONG
```

### 2. 后端已启动

后端服务已经在后台运行：
- 端口: 8888
- 进程日志: `logs/app.log`
- 查看日志: `tail -f logs/app.log`

### 3. 测试后端接口

```bash
# 测试统计接口
curl http://localhost:8888/blog/stats | python3 -m json.tool

# 测试心跳接口
curl -X POST http://localhost:8888/blog/heartbeat

# 测试访问记录
curl -X POST http://localhost:8888/blog/visit \
  -H "Content-Type: application/json" \
  -d '{"path":"/test","duration":60}'
```

### 4. 启动前端（如果还没启动）

```bash
cd blog-frontend
npm install  # 首次运行需要安装依赖
npm run dev
```

前端会在 http://localhost:5173 启动

## 🔍 验证功能

### 1. 测试心跳功能

1. 打开浏览器访问: http://localhost:5173
2. 打开开发者工具的 Network 面板
3. 等待 30 秒，应该看到 `/blog/heartbeat` 请求
4. 访问统计页面: http://localhost:5173/stats
5. 查看"当前在线人数"应该显示 1

### 2. 测试访问时长记录

1. 访问任意文章页面
2. 停留 10 秒以上
3. 切换到其他页面或关闭标签页
4. 查看数据库记录:
   ```bash
   mysql -h 127.0.0.1 -u root -p123456 leaf_admin \
     -e "SELECT * FROM page_visits ORDER BY id DESC LIMIT 5;"
   ```
5. 刷新统计页面，查看"平均访问时长"

### 3. 测试统计页面

访问: http://localhost:5173/stats

应该看到所有统计数据，包括：
- 网站运行天数（基于 2023-01-01 计算）
- 当前在线人数
- 平均访问时长
- 文章、笔记、分类、标签等数据

## 🐛 故障排查

### 问题 1: 在线人数始终为 0

**检查步骤:**

1. 确认 Redis 运行:
   ```bash
   redis-cli ping
   ```

2. 查看后端日志:
   ```bash
   tail -f logs/app.log | grep Redis
   # 应该看到: Redis connected successfully
   ```

3. 检查前端是否发送心跳:
   - 打开浏览器开发者工具
   - Network 面板搜索 "heartbeat"
   - 应该每 30 秒有一个请求

4. 手动测试心跳:
   ```bash
   curl -X POST http://localhost:8888/blog/heartbeat
   curl http://localhost:8888/blog/stats | grep online_count
   ```

### 问题 2: 平均访问时长为 0

**检查步骤:**

1. 确认数据库有记录:
   ```bash
   mysql -h 127.0.0.1 -u root -p123456 leaf_admin \
     -e "SELECT COUNT(*) FROM page_visits;"
   ```

2. 检查前端是否发送访问记录:
   - 打开浏览器开发者工具
   - Network 面板搜索 "visit"
   - 切换页面时应该有请求

3. 手动添加测试数据:
   ```bash
   curl -X POST http://localhost:8888/blog/visit \
     -H "Content-Type: application/json" \
     -d '{"path":"/test","duration":120}'
   ```

### 问题 3: 后端未响应

**检查步骤:**

1. 确认后端进程运行:
   ```bash
   lsof -i:8888
   ```

2. 查看日志:
   ```bash
   tail -30 logs/app.log
   ```

3. 重启后端:
   ```bash
   # 停止
   lsof -ti:8888 | xargs kill

   # 启动
   ./bin/leaf-api > logs/app.log 2>&1 &
   ```

## 📈 性能优化建议

### 1. Redis 缓存统计结果

当前每次请求都查询数据库，可以添加 30 秒缓存：

```go
// 在 stats.go 中添加
const STATS_CACHE_KEY = "stats:site"
const STATS_CACHE_TTL = 30 * time.Second

// GetStats 中先检查缓存
cached, err := redis.Get(STATS_CACHE_KEY)
if err == nil && cached != "" {
    // 返回缓存的数据
}
```

### 2. 定期清理历史数据

创建定时任务清理 30 天前的访问记录：

```bash
# 添加到 crontab
0 3 * * * mysql -u root -pPASSWORD leaf_admin -e "DELETE FROM page_visits WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);"
```

### 3. 数据库索引优化

确认以下索引存在：

```sql
-- 已自动创建
SHOW INDEX FROM page_visits;
SHOW INDEX FROM views;
```

## 📝 API 文档

### 统计接口

```
GET /blog/stats

响应示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "article_count": 4,
    "chapter_count": 2,
    "category_count": 2,
    "tag_count": 2,
    "user_count": 2,
    "comment_count": 4,
    "total_views": 45,
    "today_views": 0,
    "online_count": 1,
    "avg_visit_duration": 74.25,
    "site_runtime": 1058
  }
}
```

### 心跳接口

```
POST /blog/heartbeat

响应示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  }
}
```

### 访问记录接口

```
POST /blog/visit

请求体:
{
  "path": "/blog/articles/1",
  "duration": 120
}

响应示例:
{
  "code": 0,
  "message": "success",
  "data": {
    "status": "ok"
  }
}
```

## 🎉 完成！

恭喜！所有统计功能已成功集成。

### 快速访问

- 博客首页: http://localhost:5173
- 统计页面: http://localhost:5173/stats
- 后端 API: http://localhost:8888

### 下一步

1. 访问博客并浏览几个页面
2. 等待 1 分钟后查看统计页面
3. 应该能看到在线人数和访问时长的变化

### 需要帮助？

查看详细文档：
- `docs/FRONTEND_INTEGRATION.md` - 前端集成指南
- `docs/IMPLEMENTATION_SUMMARY.md` - 实现总结
- `logs/app.log` - 后端运行日志
