package logic

import (
	"context"
	"errors"
	"fmt"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidStatusLogic {
	return &GroupValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidStatusLogic) GroupValidStatus(req *types.GroupValidStatusRequest) (resp *types.GroupValidStatusResponse, err error) {
	//1.获取群验证记录信息
	var groupValidModel group_models.GroupVerifyModel
	err = l.svcCtx.DB.Take(&groupValidModel, req.ValidID).Error
	if err != nil {
		return nil, errors.New("不存在的验证记录")
	}
	//2.处理群验证信息，判断是否已经被处理过，判断调用者是否有权限处理
	if groupValidModel.Status != 0 {
		return nil, errors.New("已经处理过该验证请求了")
	}
	// 判断我有没有权限处理这个验证请求
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "user_id = ? and group_id = ?", req.UserID, groupValidModel.GroupID).Error
	if err != nil {
		return nil, errors.New("没有处理该操作的权限")
	}
	fmt.Println(member.Role)
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("没有处理该操作的权限")
	}

	//3.处理验证信息
	switch req.Status {
	case 0: // 未操作
		return
	case 1: // 同意
		//将用户加入群里
		var member1 = group_models.GroupMemberModel{
			GroupID: groupValidModel.GroupID,
			UserID:  groupValidModel.UserID,
			Role:    3,
		}
		l.svcCtx.DB.Create(&member1)
	case 2: // 拒绝
	case 3: // 忽略
	}
	//4.更新验证记录
	l.svcCtx.DB.Model(&groupValidModel).UpdateColumn("status", req.Status)
	return
}
