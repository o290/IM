package logic

import (
	"context"
	"errors"
	"fmt"
	"server/im_group/group_models"
	"time"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupProhibitionUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupProhibitionUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupProhibitionUpdateLogic {
	return &GroupProhibitionUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupProhibitionUpdateLogic) GroupProhibitionUpdate(req *types.GroupProhibitionUpdateRequest) (resp *types.GroupProhibitionUpdateResponse, err error) {
	// 只能是群主或者是管理员才能调用
	//1.查询调用者的群成员信息
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("当前用户错误")
	}
	//2.判断调用者的角色是否是群主或管理员
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("当前用户角色错误")
	}

	//3.查询禁言对象的群成员信息
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.GroupID, req.MemberID).Error
	if err != nil {
		return nil, errors.New("当前用户错误")
	}
	//4.判断角色
	if !((member.Role == 1 && member1.Role == 2 || member1.Role == 3) || (member.Role == 2 && member1.Role == 3)) {
		return nil, errors.New("角色错误")
	}

	//5.更新禁言时间
	l.svcCtx.DB.Model(&member1).Update("prohibition_time", req.ProhibitionTime)

	//6.利用redis的过期时间去做这个禁言时间
	key := fmt.Sprintf("prohibition__%d", member1.ID)
	if req.ProhibitionTime != nil {
		//给redis设置一个key,过期时间是xxxx
		l.svcCtx.Redis.Set(context.Background(), key, "1", time.Duration(*req.ProhibitionTime)*time.Minute)
	} else {
		//取消禁言
		l.svcCtx.Redis.Del(context.Background(), key)
	}
	return
}
