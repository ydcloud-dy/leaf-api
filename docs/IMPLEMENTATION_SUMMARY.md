# 网站统计功能实现总结

## 📋 功能概览

本次实现了以下统计功能：

### ✅ 已实现的统计项

1. **文章篇数** - 已发布文章总数
2. **笔记篇数** - 章节总数（Chapter）
3. **文章分类数** - Category 总数
4. **文章标签数** - Tag 总数
5. **评论数** - 已审核评论总数
6. **总浏览量** - 所有文章浏览量之和
7. **今日浏览量** - 今天的浏览记录数
8. **当前在线人数** ⭐ 新功能（Redis 实时统计）
9. **平均访问时长** ⭐ 新功能（数据库记录，按页面统计）
10. **网站运行天数** - 从 settings 表读取启动日期

## 🎯 核心技术方案

### 1. 在线人数统计

**技术栈**: Redis + 心跳机制

**工作原理**:
- 前端每 30 秒发送心跳请求
- 后端将用户标识存入 Redis（登录用户用 UserID，游客用 IP）
- Redis Key 设置 60 秒过期时间
- 统计时查询所有未过期的 Key 数量

**优点**:
- ✅ 实时性高
- ✅ 自动清理过期数据
- ✅ 不增加数据库压力
- ✅ 支持未登录用户

### 2. 平均访问时长统计

**技术栈**: MySQL + 前端上报

**工作原理**:
- 前端记录页面打开时间
- 页面关闭/切换时使用 `navigator.sendBeacon` 上报停留时长
- 后端存入 `page_visits` 表
- 统计时计算最近 24 小时的平均时长

**优点**:
- ✅ 数据准确
- ✅ 可追溯历史
- ✅ 支持多维度分析

## 📁 文件变更清单

### 新增文件

```
pkg/redis/redis.go              # Redis 客户端封装
internal/service/online.go       # 在线用户追踪服务
internal/model/po/models.go      # 添加 PageVisit 模型
docs/FRONTEND_INTEGRATION.md     # 前端集成文档
docs/IMPLEMENTATION_SUMMARY.md   # 本文档
```

### 修改文件

```
config.yaml                      # 添加 Redis 配置
config/config.go                 # 添加 RedisConfig 结构体
internal/service/stats.go        # 扩展统计接口
internal/server/router.go        # 添加新路由
internal/server/http.go          # 初始化新服务
cmd/app.go                       # 添加 Redis 初始化和清理
go.mod                           # 添加 Redis 依赖
```

## 🚀 部署步骤

### 1. 安装 Redis

**macOS**:
```bash
brew install redis
brew services start redis
```

**Ubuntu/Debian**:
```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis
```

**Docker**:
```bash
docker run -d --name redis -p 6379:6379 redis:alpine
```

### 2. 更新配置

编辑 `config.yaml`：

```yaml
redis:
  host: 127.0.0.1
  port: 6379
  password:           # 如果有密码请填写
  db: 0
  pool_size: 10
```

### 3. 初始化网站启动日期

在数据库中执行：

```sql
INSERT INTO settings (`key`, `value`, updated_at)
VALUES ('site_start_date', '2020-01-01', NOW())
ON DUPLICATE KEY UPDATE value = value;
```

将 `2020-01-01` 替换为你的网站实际启动日期。

### 4. 构建并运行

```bash
# 安装依赖
go mod download

# 构建
go build -o bin/leaf-api

# 运行
./bin/leaf-api
```

### 5. 验证 Redis 连接

查看启动日志，应该看到：

```
INFO Redis connected successfully
```

如果看到警告：

```
WARN Failed to initialize Redis: ...
WARN Online user tracking and visit duration recording will be disabled
```

说明 Redis 连接失败，请检查配置和 Redis 服务状态。

### 6. 数据库迁移

应用启动时会自动创建 `page_visits` 表，无需手动操作。

## 🔌 前端集成

详细的前端集成文档请查看 `docs/FRONTEND_INTEGRATION.md`。

### 快速集成步骤

1. **创建 API 服务** (`src/api/stats.js`)
2. **实现心跳 Composable** (`src/composables/useHeartbeat.js`)
3. **实现访问追踪 Composable** (`src/composables/useVisitTracking.js`)
4. **在根组件使用** (`App.vue`)
5. **展示统计数据** (统计页面)

### 关键 API 接口

```
POST /blog/heartbeat          # 发送心跳
POST /blog/visit              # 记录访问时长
GET  /blog/stats              # 获取统计数据
```

## 📊 API 响应示例

### GET /blog/stats

```json
{
  "code": 200,
  "data": {
    "article_count": 92,
    "chapter_count": 1124,
    "category_count": 11,
    "tag_count": 9,
    "user_count": 5,
    "comment_count": 12,
    "total_views": 123456,
    "today_views": 2913,
    "online_count": 9,
    "avg_visit_duration": 77.5,
    "site_runtime": 1849
  }
}
```

## 🧪 功能测试

### 1. 测试在线人数

```bash
# 1. 发送心跳
curl -X POST http://localhost:8888/blog/heartbeat

# 2. 查看统计（应该看到 online_count 增加）
curl http://localhost:8888/blog/stats

# 3. 等待 60 秒后再查看（online_count 应该减少）
sleep 60
curl http://localhost:8888/blog/stats
```

### 2. 测试访问时长记录

```bash
curl -X POST http://localhost:8888/blog/visit \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/blog/articles/1",
    "duration": 120
  }'

# 查看数据库
mysql -u root -p leaf_admin -e "SELECT * FROM page_visits ORDER BY id DESC LIMIT 5;"
```

### 3. 测试统计接口

```bash
curl http://localhost:8888/blog/stats | jq '.'
```

## 🐛 故障排查

### 问题 1: 在线人数始终为 0

**检查清单**:
1. Redis 是否启动？
   ```bash
   redis-cli ping  # 应该返回 PONG
   ```

2. 后端日志是否显示 Redis 连接成功？

3. 前端是否发送心跳请求？
   - 打开浏览器 Network 面板
   - 查找 `/blog/heartbeat` 请求

4. Redis 中是否有数据？
   ```bash
   redis-cli
   KEYS online:*
   ```

### 问题 2: 平均访问时长为 0

**检查清单**:
1. 数据库中是否有访问记录？
   ```sql
   SELECT COUNT(*) FROM page_visits;
   ```

2. 前端是否正确发送访问时长数据？

3. 时间范围是否正确？（默认统计最近 24 小时）

### 问题 3: 网站运行天数为 0

**原因**: `settings` 表中没有 `site_start_date` 配置。

**解决**:
```sql
INSERT INTO settings (`key`, `value`, updated_at)
VALUES ('site_start_date', '2020-01-01', NOW());
```

## 📈 性能优化建议

### 1. 统计接口缓存

当前每次请求都实时查询数据库和 Redis，可以添加缓存：

```go
// 使用 Redis 缓存统计结果，有效期 30 秒
func (s *StatsService) GetStats(c *gin.Context) {
    cacheKey := "stats:site"

    // 尝试从缓存获取
    cached, err := redis.Get(cacheKey)
    if err == nil && cached != "" {
        var stats StatsData
        json.Unmarshal([]byte(cached), &stats)
        response.Success(c, stats)
        return
    }

    // 查询数据库...
    // ...

    // 写入缓存
    data, _ := json.Marshal(stats)
    redis.SetWithExpire(cacheKey, string(data), 30*time.Second)

    response.Success(c, stats)
}
```

### 2. 心跳请求优化

如果用户打开多个标签页，可以使用 `localStorage` 避免重复发送心跳。

详见 `docs/FRONTEND_INTEGRATION.md` 的"可选优化"章节。

### 3. 数据库索引

`page_visits` 表已包含必要的索引：

```sql
-- 已自动创建
INDEX idx_user_id (user_id)
INDEX idx_ip (ip)
INDEX idx_created_at (created_at)
```

### 4. 定期清理历史数据

建议定期清理超过 30 天的 `page_visits` 记录：

```sql
DELETE FROM page_visits WHERE created_at < DATE_SUB(NOW(), INTERVAL 30 DAY);
```

可以使用 cron 任务或后台定时任务执行。

## 🔒 安全考虑

### 1. 防止恶意刷心跳

当前实现按 IP 或 UserID 去重，但仍可能被恶意刷新。

**建议**:
- 添加速率限制中间件
- 记录异常 IP 并加入黑名单
- 监控单个 IP 的心跳频率

### 2. 访问时长数据验证

当前已做基本验证（`duration >= 0`），但可以添加更严格的限制：

```go
// 拒绝超过 24 小时的异常时长
if req.Duration > 86400 {
    response.Error(c, 400, "无效的访问时长")
    return
}
```

### 3. Redis 密码保护

生产环境强烈建议为 Redis 设置密码：

```yaml
redis:
  password: "your-strong-password"
```

## 📚 扩展功能建议

### 1. 实时在线用户列表

可以扩展显示哪些用户在线：

```go
// 获取在线用户详情
func (s *OnlineService) GetOnlineUsers() ([]OnlineUser, error) {
    keys, _ := redis.Keys(onlineUserPrefix + "*")
    users := make([]OnlineUser, 0)

    for _, key := range keys {
        userID := strings.TrimPrefix(key, onlineUserPrefix)
        // 查询用户信息...
        users = append(users, user)
    }

    return users, nil
}
```

### 2. 按页面统计访问量

当前只统计平均时长，可以扩展统计每个页面的访问量：

```go
func (s *VisitService) GetPageStats() ([]PageStat, error) {
    var stats []PageStat
    s.data.GetDB().Model(&po.PageVisit{}).
        Select("path, COUNT(*) as visit_count, AVG(duration) as avg_duration").
        Where("created_at >= ?", time.Now().Add(-24*time.Hour)).
        Group("path").
        Order("visit_count DESC").
        Limit(10).
        Scan(&stats)
    return stats, nil
}
```

### 3. 访客地域分析

基于 IP 地址进行地域分析（需要集成 IP 地理位置库）。

### 4. 访问趋势图

记录每小时或每天的访问统计，生成趋势图表。

## ✅ 完成检查清单

- [x] Redis 配置和初始化
- [x] 在线用户追踪服务
- [x] 访问时长记录服务
- [x] 统计 API 扩展
- [x] 路由和服务注册
- [x] 数据模型迁移
- [x] 前端集成文档
- [x] 部署文档
- [ ] Redis 启动并运行
- [ ] 数据库 `settings` 表配置网站启动日期
- [ ] 前端集成心跳和访问追踪
- [ ] 功能测试验证

## 🎉 总结

本次实现为博客系统添加了完整的在线统计功能，包括：

1. **实时在线人数追踪** - 基于 Redis 的高性能实时统计
2. **访问时长分析** - 基于数据库的持久化记录
3. **全面的网站统计** - 文章、分类、标签、评论等多维度统计
4. **优雅降级** - Redis 失败不影响主要功能
5. **详细文档** - 完整的前后端集成指南

核心代码位置：
- 后端服务：`internal/service/online.go`
- 统计接口：`internal/service/stats.go`
- Redis 封装：`pkg/redis/redis.go`
- 前端文档：`docs/FRONTEND_INTEGRATION.md`

如有问题，请查阅文档或检查日志。
