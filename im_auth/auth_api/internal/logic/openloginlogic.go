package logic

import (
	"context"
	"errors"
	"fmt"
	"server/im_auth/auth_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils/jwt"
	"server/utils/open_login"

	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type Open_loginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOpen_loginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *Open_loginLogic {
	return &Open_loginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *Open_loginLogic) Open_login(req *types.OpenLoginRequest) (resp *types.LoginResponse, err error) {
	type OpenInfo struct {
		NickName string
		OpenID   string
		Avatar   string
	}
	var info OpenInfo
	switch req.Flag {
	case "qq":
		qqInfo, openErr := open_login.NewQQLogin(req.Code, open_login.QQConfig{
			AppID:    l.svcCtx.Config.QQ.AppID,
			AppKey:   l.svcCtx.Config.QQ.AppKey,
			Redirect: l.svcCtx.Config.QQ.Redirect,
		})
		info = OpenInfo{
			OpenID:   qqInfo.OpenID,
			NickName: qqInfo.Nickname,
			Avatar:   qqInfo.Avatar,
		}
		err = openErr
	default:
		err = errors.New("不支持第三方登录")
	}
	if err != nil {
		logx.Error(err)
		return nil, errors.New("登录失败")
	}
	var user auth_models.UserModel
	err = l.svcCtx.DB.Take(&user, "open_id = ?", info.OpenID).Error
	if err != nil {
		//注册逻辑
		fmt.Println("注册服务")
		res, err := l.svcCtx.UserRpc.UserCreate(context.Background(), &user_rpc.UserCreateRequest{
			NickName:       info.NickName,
			Password:       "",
			Role:           2,
			Avatar:         info.Avatar,
			OpenId:         info.OpenID,
			RegisterSource: "qq",
		})
		if err != nil {
			logx.Error(err)
			return nil, errors.New("登录失败")
		}
		user.Model.ID = uint(res.UserId)
		user.Role = 2
		user.Nickname = info.NickName
	}
	//登录逻辑
	//生成jtw令牌，token
	token, err1 := jwt.GenToken(jwt.JwtPayload{
		UserID:   user.ID,
		NickName: user.Nickname,
		Role:     user.Role,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err1 != nil {
		logx.Error(err1)
		err = errors.New("服务内部错误")
		return nil, err1
	}
	//返回响应
	return &types.LoginResponse{
		Token: token,
	}, nil
	return
}
