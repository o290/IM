package logic

import (
	"context"
	"errors"
	"fmt"
	"server/utils"
	"server/utils/jwt"

	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
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

func (l *AuthenticationLogic) Authentication(req *types.AuthenticationRequest) (resp *types.AuthenticationResponse, err error) {
	if utils.InitListByRegex(l.svcCtx.Config.WriteList, req.ValidPath) {
		logx.Infof("%s在白名单中", req.ValidPath)
		return
	}
	if req.Token == "" {
		err = errors.New("认证失败")
		return
	}
	claims, err := jwt.ParseToken(req.Token, l.svcCtx.Config.Auth.AccessSecret)
	if err != nil {
		logx.Error(err.Error())
		err = errors.New("认证失败")
		return
	}
	_, err = l.svcCtx.Redis.Get(l.ctx, fmt.Sprintf("logout_%s", req.Token)).Result()
	//Get err==nil时表示存在
	if err == nil {
		logx.Error("已注销")
		err = errors.New("认证失败")
		return
	}
	return &types.AuthenticationResponse{
		UserID: claims.UserID,
		Role:   int(claims.Role),
	}, nil
}
