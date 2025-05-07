package logic

import (
	"context"
	"errors"
	"server/common/models/ctype"
	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"
	"server/im_user/user_models"
	"server/utils/maps"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoUpdateLogic {
	return &UserInfoUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoUpdateLogic) UserInfoUpdate(req *types.UserInfoUpdateRequest) (resp *types.UserInfoUpdateResponse, err error) {
	//1.利用反射将req指针指向的结构体根据user标签转换成map
	userMaps := maps.RefToMap(*req, "user")
	if len(userMaps) != 0 {
		var user user_models.UserModel
		//2.根据id获取用户信息
		err = l.svcCtx.DB.Take(&user, req.UserID).Error
		if err != nil {
			return nil, errors.New("用户不存在")
		}
		//3.更新
		//Updates用于批量更新，会把 userMaps 中的键值对更新到数据库中 user 记录对应的字段上
		//model是用于指定操作的数据表
		err = l.svcCtx.DB.Model(&user).Updates(userMaps).Error
		if err != nil {
			logx.Error(userMaps)
			logx.Error(err)
			return nil, errors.New("用户信息更新失败")
		}
	}
	//4.利用反射将req指针指向的结构体根据user_conf标签转换成map
	userConfMaps := maps.RefToMap(*req, "user_conf")
	//fmt.Println(userConfMaps)
	if len(userConfMaps) != 0 {
		var userConf user_models.UserConfModel
		//5.根据id获取用户信息
		err = l.svcCtx.DB.Take(&userConf, req.UserID).Error
		if err != nil {
			return nil, errors.New("用户不存在")
		}
		//6.检查是否存在verification_question字段
		//如果存在就将他删除，并转换成ctype.VerificationQuestion类型的结构体
		//map[save_pwd:true verification_question:map[answer1:school answer2:schhol2 problem1:question2222 problem2:Q2]]
		//所以要将map转换为struct
		VerificationQuestion, ok := userConfMaps["verification_question"]
		if ok {
			delete(userConfMaps, "verification_question")
			data := ctype.VerificationQuestion{}
			//转化，类型断言
			maps.MapToStruct(VerificationQuestion.(map[string]any), &data)
			l.svcCtx.DB.Model(&userConf).Updates(&user_models.UserConfModel{
				VerificationQuestion: &data,
			})
		}
		//7.更新
		err = l.svcCtx.DB.Model(&userConf).Updates(userConfMaps).Error
		if err != nil {
			logx.Error(userConfMaps)
			logx.Error(err)
			return nil, errors.New("用户配置信息更新失败")
		}
	}
	return
}
