package logic

import (
	"context"
	"errors"
	"server/im_user/user_models"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserValidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserValidLogic {
	return &UserValidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserValidLogic) UserValid(req *types.UserValidRequest) (resp *types.UserValidResponse, err error) {
	//如果互为好友，就不需要添加
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("已经互为好友")
	}

	var userConf user_models.UserConfModel
	err = l.svcCtx.DB.Take(&userConf, "user_id=?", req.FriendID).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	resp = new(types.UserValidResponse)
	resp.Verification = userConf.Verification
	switch userConf.Verification {
	case 0: //不允许任何人添加
	case 1: //允许任何人添加，直接成为好友
	//首先往验证表添加一条记录，然后通过
	case 2: //需要验证问题
	case 3, 4: //需要正确回答问题，返回问题
		if userConf.VerificationQuestion != nil {
			resp.VerificationQuestion = types.VerificationQuestion{
				Problem1: userConf.VerificationQuestion.Problem1,
				Problem2: userConf.VerificationQuestion.Problem2,
				Problem3: userConf.VerificationQuestion.Problem3,
			}
		}
	default:

	}
	return
}
