package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"beaver/app/user/user_rpc/types/user_rpc"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidListLogic {
	return &ValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidListLogic) ValidList(req *types.ValidListReq) (resp *types.ValidListRes, err error) {
	// 参数验证
	if req.UserID == "" {
		return nil, errors.New("用户ID不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 查询好友验证列表
	fvs, count, _ := list_query.ListQuery(l.svcCtx.DB, friend_models.FriendVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where: l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
		// 移除Preload，微服务架构中通过RPC获取用户信息
	})

	// 收集需要获取用户信息的UserID列表
	var userIds []string
	userIdSet := make(map[string]bool)
	for _, fv := range fvs {
		if fv.SendUserID != "" && !userIdSet[fv.SendUserID] {
			userIds = append(userIds, fv.SendUserID)
			userIdSet[fv.SendUserID] = true
		}
		if fv.RevUserID != "" && !userIdSet[fv.RevUserID] {
			userIds = append(userIds, fv.RevUserID)
			userIdSet[fv.RevUserID] = true
		}
	}

	// 批量获取用户信息
	userInfoMap := make(map[string]*user_rpc.UserInfo)
	if len(userIds) > 0 {
		userListResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIds,
		})
		if err != nil {
			l.Logger.Errorf("批量获取用户信息失败: %v", err)
			// 不返回错误，继续处理，为没有用户信息的设置默认值
		} else {
			userInfoMap = userListResp.UserInfo
		}
	}

	var list []types.FriendValidInfo
	for _, fv := range fvs {
		info := types.FriendValidInfo{
			Message:   fv.Message,
			Id:        fv.UUID,
			Source:    fv.Source,
			CreatedAt: fv.CreatedAt.String(),
		}

		if fv.SendUserID == req.UserID {
			// 我是发起方
			info.UserID = fv.RevUserID
			if userInfo, exists := userInfoMap[fv.RevUserID]; exists && userInfo != nil {
				info.NickName = userInfo.NickName
				info.Avatar = userInfo.Avatar
			} else {
				info.NickName = "未知用户"
				info.Avatar = ""
			}
			info.Flag = "send"
			info.Status = fv.RevStatus
		} else if fv.RevUserID == req.UserID {
			// 我是接收方
			info.UserID = fv.SendUserID
			if userInfo, exists := userInfoMap[fv.SendUserID]; exists && userInfo != nil {
				info.NickName = userInfo.NickName
				info.Avatar = userInfo.Avatar
			} else {
				info.NickName = "未知用户"
				info.Avatar = ""
			}
			info.Flag = "receive"
			info.Status = fv.RevStatus
		} else {
			// 这种情况理论上不应该发生，跳过
			continue
		}

		list = append(list, info)
	}

	l.Logger.Infof("获取好友验证列表成功: userID=%s, count=%d", req.UserID, len(list))
	return &types.ValidListRes{
		Count: count,
		List:  list,
	}, nil
}
