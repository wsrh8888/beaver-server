package robot

import (
	"context"
	"encoding/json"
	"errors"

	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_api/internal/utils"
	"beaver/app/open/open_rpc/types/open_rpc"
	"beaver/app/open/openevent"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveRobotFromGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRemoveRobotFromGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveRobotFromGroupLogic {
	return &RemoveRobotFromGroupLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RemoveRobotFromGroupLogic) RemoveRobotFromGroup(req *types.RemoveRobotFromGroupReq, authorization string) (resp *types.RemoveRobotFromGroupRes, err error) {
	if req.GroupID == "" {
		return nil, errors.New("groupId 不能为空")
	}

	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppCapability(app, true, false); err != nil {
		return nil, err
	}

	robot, err := utils.EnsureAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	_, err = l.svcCtx.GroupRpc.RemoveGroupMember(l.ctx, &group_rpc.RemoveGroupMemberReq{
		GroupId:    req.GroupID,
		UserId:     robot.RobotID,
		OperatedBy: token.AppID,
	})
	if err != nil {
		return nil, err
	}

	go func() {
		body, _ := json.Marshal(map[string]interface{}{
			"group_id":    req.GroupID,
			"robot_id":    robot.RobotID,
			"operator_id": token.AppID,
		})
		_, _ = l.svcCtx.OpenRpc.DispatchPlatformEvent(context.Background(), &open_rpc.DispatchPlatformEventReq{
			AppId:     token.AppID,
			EventType: openevent.EventIMChatMemberBotRemoved,
			EventJson: string(body),
		})
	}()

	return &types.RemoveRobotFromGroupRes{
		Success: true,
	}, nil
}
