package logic

import (
	"context"
	"encoding/json"
	"errors"
	"server/im_user/user_models"
	"server/im_user/user_rpc/types/user_rpc"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo(req *types.UserInfoRequest) (resp *types.UserInfoResponse, err error) {
	//1.获取用户信息
	res, err := l.svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
		UserId: uint32(req.UserID),
	})
	if err != nil {
		return nil, err
	}
	//2.解析json，反序列化
	var user user_models.UserModel
	err = json.Unmarshal(res.Data, &user)
	if err != nil {
		logx.Error(err)
		return nil, errors.New("数据错误")
	}
	//3.构造响应
	p := &types.UserInfoResponse{
		UserID:        user.ID,
		NickName:      user.Nickname,
		Abstract:      user.Abstract,
		Avatar:        user.Avatar,
		RecallMessage: nil,
		FriendOnline:  user.UserConfModel.FriendOnline,
		Sound:         user.UserConfModel.Sound,
		SecureLink:    user.UserConfModel.SecureLink,
		SavePwd:       user.UserConfModel.SavePwd,
		SearchUser:    user.UserConfModel.SearchUser,
		Verification:  user.UserConfModel.Verification,
	}
	resp = p
	if user.UserConfModel.VerificationQuestion != nil {
		resp.VerificationQuestion = &types.VerificationQuestion{
			Problem1: user.UserConfModel.VerificationQuestion.Problem1,
			Problem2: user.UserConfModel.VerificationQuestion.Problem2,
			Problem3: user.UserConfModel.VerificationQuestion.Problem3,
			Answer1:  user.UserConfModel.VerificationQuestion.Answer1,
			Answer2:  user.UserConfModel.VerificationQuestion.Answer2,
			Answer3:  user.UserConfModel.VerificationQuestion.Answer3,
		}
	}
	return resp, nil
}
