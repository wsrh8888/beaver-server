# 好友模块数据模型

## 模型说明

### FriendModel (好友关系表)
- **UUID**: 全局唯一标识符，替代数据库自增ID用于前端同步
- **SendUserID/RevUserID**: 好友关系的双方用户ID
- **SendUserNotice/RevUserNotice**: 双方备注信息
- **Source**: 好友关系来源（qrcode/search/group/recommend）
- **IsDeleted**: 软删除标记
- **Version**: 版本号，用于数据同步和乐观锁

### FriendVerifyModel (好友验证表)
- **UUID**: 全局唯一标识符，替代数据库自增ID用于前端同步
- **SendUserID/RevUserID**: 验证请求的双方用户ID
- **SendStatus/RevStatus**: 发送方和接收方的验证状态
- **Message**: 验证附加消息
- **Source**: 添加好友来源（qrcode/search/group/recommend）
- **Version**: 版本号，用于数据同步和乐观锁

## Version 字段管理规则

### 版本递增时机
1. **创建记录**: 新记录的初始版本号为1（通过VersionGen.GetNextVersion获取）
2. **更新记录**: 每次更新时递增版本号
3. **删除记录**: 删除时递增版本号（软删除）

### 版本使用场景
- **数据同步**: 前端通过版本号判断数据是否需要更新
- **乐观锁**: 避免并发更新冲突
- **增量同步**: 只同步版本号大于指定值的记录

### Logic 层实现

#### 创建操作
```go
// 获取下一个版本号
nextVersion, err := l.svcCtx.VersionGen.GetNextVersion("friends")
// 创建记录时设置初始版本
friendModel := FriendModel{
    UUID:    uuid.New().String(),
    Version: nextVersion, // 初始版本为1
    // ... 其他字段
}
```

#### 更新操作
```go
// 获取下一个版本号
nextVersion, err := l.svcCtx.VersionGen.GetNextVersion("friends")
// 批量更新多个字段
err = l.svcCtx.DB.Model(&FriendModel{}).Where("uuid = ?", friendUUID).Updates(map[string]interface{}{
    "field1": newValue,
    "version": nextVersion, // 递增版本号
}).Error
```

#### 删除操作
```go
// 获取下一个版本号
nextVersion, err := l.svcCtx.VersionGen.GetNextVersion("friends")
// 软删除时更新版本
err = l.svcCtx.DB.Model(&FriendModel{}).Where("uuid = ?", friendUUID).Updates(map[string]interface{}{
    "is_deleted": 1,
    "version": nextVersion, // 递增版本号
}).Error
```

## API 响应格式

### 所有返回ID的字段已改为UUID
- `FriendValidInfo.Id`: 验证记录UUID
- `FriendSyncItem.UUID`: 好友记录UUID
- `FriendVerifySyncItem.UUID`: 验证记录UUID
- `SearchValidInfoRes.ValidID`: 验证记录UUID

### 版本字段在同步API中返回
- `FriendSyncItem.Version`: 好友记录版本号
- `FriendVerifySyncItem.Version`: 验证记录版本号

## 前端同步策略

前端通过版本号实现增量同步：
1. 记录本地最大版本号
2. 请求版本号大于本地版本的记录
3. 更新本地数据和版本号

这样可以避免全量同步，提高性能。
