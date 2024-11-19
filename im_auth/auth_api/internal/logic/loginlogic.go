package logic

import (
	"context"
	"errors"
	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"
	"server/im_auth/auth_models"
	"server/utils/jwt"
	"server/utils/pwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	//提供日志功能
	logx.Logger
	//用于传递请求的上下文信息
	ctx context.Context
	//包含了服务相关的配置信息和数据库连接
	svcCtx *svc.ServiceContext
}

// 创建一个新的 LoginLogic 实例，并初始化其中的字段
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 执行登录逻辑的
func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	//查找用户
	//user用于存放从数据库中查询到的用户数据
	var user auth_models.UserModel
	//l.svcCtx.DB.Take(&user, req.UserName) 通过数据库查询 req.UserName 对应的用户
	//并将结果填充到 user 变量中。
	//这里的id指的就是username，省去了用户名
	err = l.svcCtx.DB.Take(&user, "id=?", req.UserName).Error
	if err != nil {
		err = errors.New("用户名或密码错误")
		return
	}
	//验证密码
	if !pwd.CheckPwd(user.Pwd, req.Password) {
		err = errors.New("用户名或密码错误")
		return
	}
	//判断用户的注册来源，第三方登录来的不能通过用户密码登录

	//生成jtw令牌，token
	token, err := jwt.GenToken(jwt.JwtPayload{
		UserID:   user.ID,
		NickName: user.Nickname,
		Role:     user.Role,
	}, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Error(err)
		err = errors.New("服务内部错误")
		return
	}
	//返回响应
	return &types.LoginResponse{
		Token: token,
	}, nil
}
