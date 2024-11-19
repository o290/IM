package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"server/utils/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"server/im_auth/auth_api/internal/svc"
)

type AuthenticationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuthenticationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthenticationLogic {
	return &AuthenticationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuthenticationLogic) Authentication(token string) (resp string, err error) {
	if token == "" {
		err = errors.New("认证失败")
		return
	}
	payload, err := jwt.ParseToken(token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		err = errors.New("认证失败")
		return
	}
	_, err = l.svcCtx.Redis.Get(l.ctx, fmt.Sprintf("logout_%d", payload.UserID)).Result()
	//Get err==nil时表示不存在，没有注销
	if err == redis.Nil {
		resp = "认证成功"
	} else {
		err = errors.New("认证失败")
		return
	}
	return resp, nil
}
