package logic

import (
	"context"
	"fmt"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取好友验证列表
func NewGetFriendVerifyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyListLogic {
	return &GetFriendVerifyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendVerifyListLogic) GetFriendVerifyList(req *types.GetFriendVerifyListReq) (resp *types.GetFriendVerifyListRes, err error) {
	// 设置默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	// 构建查询条件
	query := l.svcCtx.DB.Model(&friend_models.FriendVerifyModel{}).
		Preload("SendUserModel").
		Preload("RevUserModel")

	// 按发送用户ID筛选
	if req.SendUserId != "" {
		query = query.Where("send_user_id = ?", req.SendUserId)
	}

	// 按接收用户ID筛选
	if req.RevUserId != "" {
		query = query.Where("rev_user_id = ?", req.RevUserId)
	}

	// 按发送方状态筛选
	if req.SendStatus > 0 {
		query = query.Where("send_status = ?", req.SendStatus)
	}

	// 按接收方状态筛选
	if req.RevStatus > 0 {
		query = query.Where("rev_status = ?", req.RevStatus)
	}

	// 时间范围筛选
	if req.StartTime != "" {
		query = query.Where("created_at >= ?", req.StartTime)
	}
	if req.EndTime != "" {
		query = query.Where("created_at <= ?", req.EndTime)
	}

	// 查询总数
	var total int64
	err = query.Count(&total).Error
	if err != nil {
		logx.Errorf("查询好友验证总数失败: %v", err)
		return nil, err
	}

	// 查询列表
	var verifies []friend_models.FriendVerifyModel
	err = query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&verifies).Error
	if err != nil {
		logx.Errorf("查询好友验证列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.FriendVerifyInfo, len(verifies))
	for i, verify := range verifies {
		var sendUserName, revUserName string
		if verify.SendUserModel.UserID != "" {
			sendUserName = verify.SendUserModel.NickName
		}
		if verify.RevUserModel.UserID != "" {
			revUserName = verify.RevUserModel.NickName
		}

		list[i] = types.FriendVerifyInfo{
			Id:           fmt.Sprintf("%d", verify.Id),
			SendUserId:   verify.SendUserID,
			SendUserName: sendUserName,
			RevUserId:    verify.RevUserID,
			RevUserName:  revUserName,
			SendStatus:   int(verify.SendStatus),
			RevStatus:    int(verify.RevStatus),
			Message:      verify.Message,
			CreateTime:   time.Time(verify.CreatedAt).Format(time.RFC3339),
			UpdateTime:   time.Time(verify.UpdatedAt).Format(time.RFC3339),
		}
	}

	return &types.GetFriendVerifyListRes{
		List:  list,
		Total: total,
	}, nil
}
