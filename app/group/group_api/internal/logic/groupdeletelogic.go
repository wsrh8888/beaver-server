package logic

import (
	"context"
	"errors"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"
	"beaver/app/group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupDeleteLogic {
	return &GroupDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupDeleteLogic) GroupDelete(req *types.GroupDeleteReq) (resp *types.GroupDeleteRes, err error) {
	// todo: add your logic here and delete this line
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&groupMember, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("用户不是群组成员")
	}
	if groupMember.Role != 1 {
		return nil, errors.New("只有群主才可以删除群组")
	}

	// 获取该群的版本号（独立递增）
	groupVersion := l.svcCtx.VersionGen.GetNextVersion("groups", "group_id", req.GroupID)
	if groupVersion == -1 {
		l.Logger.Errorf("获取群组版本号失败")
		return nil, errors.New("获取版本号失败")
	}

	//  群成员要删掉
	var memberList []group_models.GroupMemberModel
	err = l.svcCtx.DB.Find(&memberList, "group_id = ?", req.GroupID).Delete(&memberList).Error
	if err != nil {
		return nil, errors.New("删除群成员失败")
	}

	// 群组要删掉
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, req.GroupID).Delete(&group).Error
	if err != nil {
		return nil, errors.New("删除群组失败")
	}

	return &types.GroupDeleteRes{
		Version: groupVersion,
	}, nil
}
