package logic

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"server/common/models/ctype"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
	"server/im_group/group_models"
)

type GroupValidAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidAddLogic {
	return &GroupValidAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidAddLogic) GroupValidAdd(req *types.AddGroupRequest) (resp *types.AddGroupResponse, err error) {
	//自己已经在群里面
	//1.判断自己是否在群中
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err == nil {
		return nil, errors.New("请勿重复加群")
	}
	//2.判断群是否存在,存在则获取信息
	var group group_models.GroupModel
	err = l.svcCtx.DB.Take(&group, req.GroupID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}

	//3.构造响应,构造验证记录
	resp = new(types.AddGroupResponse)
	var verifyModel = group_models.GroupVerifyModel{
		GroupID:           req.GroupID,
		UserID:            req.UserID,
		Status:            0,
		AdditionalMessage: req.Verify,
		Type:              1, //加群
	}
	//4.根据群的设置判断
	switch group.Verification {
	case 0: //不允许任何人添加
		return nil, errors.New("不允许任何人加群")
	case 1: //允许任何人添加，直接加群
		//首先往验证表添加一条记录，然后通过
		verifyModel.Status = 1
		var groupMember = group_models.GroupMemberModel{
			GroupID: req.GroupID,
			UserID:  req.UserID,
			Role:    3,
		}
		l.svcCtx.DB.Create(&groupMember)
	case 2: //需要验证问题

	case 3:
		if req.VerificationQuestion != nil {
			verifyModel.VerificationQuestion = &ctype.VerificationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
				Answer1:  req.VerificationQuestion.Answer1,
				Answer2:  req.VerificationQuestion.Answer2,
				Answer3:  req.VerificationQuestion.Answer3,
			}
		}
	case 4: //需要正确回答问题，返回问题
		//判断问题是否正确回答
		var count int
		if req.VerificationQuestion != nil && group.VerificationQuestion != nil {
			//要考虑一个问题两个问题三个问题的情况
			if req.VerificationQuestion.Answer1 != nil && group.VerificationQuestion.Answer1 != nil {
				if *req.VerificationQuestion.Answer1 == *group.VerificationQuestion.Answer1 {
					count++
				}
			}
			if req.VerificationQuestion.Answer2 != nil && group.VerificationQuestion.Answer2 != nil {
				if *req.VerificationQuestion.Answer2 == *group.VerificationQuestion.Answer2 {
					count++
				}
			}
			if req.VerificationQuestion.Answer3 != nil && group.VerificationQuestion.Answer3 != nil {
				if *req.VerificationQuestion.Answer3 == *group.VerificationQuestion.Answer3 {
					count++
				}
			}
			if count != group.ProblemCount() {
				return nil, errors.New("答案错误")
			}
			//直接加群
			verifyModel.Status = 1
			verifyModel.VerificationQuestion = group.VerificationQuestion
			// 把用户加到群里面
			var groupMember = group_models.GroupMemberModel{
				GroupID: req.GroupID,
				UserID:  req.UserID,
				Role:    3,
			}
			l.svcCtx.DB.Create(&groupMember)
		} else {
			return nil, errors.New("答案错误")
		}

	default:

	}
	//创建验证记录
	err = l.svcCtx.DB.Create(&verifyModel).Error
	if err != nil {
		return
	}
	return
}
