package logic

import (
	"context"
	"errors"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
	"server/im_group/group_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidLogic {
	return &GroupValidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GroupValid 显示群验证信息
func (l *GroupValidLogic) GroupValid(req *types.GroupValidRequest) (resp *types.GroupValidResponse, err error) {
	//自己已经在群里面
	//1.判断自己是否在群中
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err == nil {
		return nil, errors.New("请勿重复加群")
	}
	//2.判断群是否存在
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, req.GroupID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}
	//3.构造响应
	resp = new(types.GroupValidResponse)
	resp.Verification = group.Verification
	switch group.Verification {
	case 0: //不允许任何人添加
	case 1: //允许任何人添加，直接成为好友
	//首先往验证表添加一条记录，然后通过
	case 2: //需要验证问题
	case 3, 4: //需要正确回答问题，返回问题
		if group.VerificationQuestion != nil {
			resp.VerificationQuestion = types.VerificationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
			}
		}
	default:

	}

	return
}
