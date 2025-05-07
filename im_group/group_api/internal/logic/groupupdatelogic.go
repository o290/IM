package logic

import (
	"context"
	"errors"
	"server/common/models/ctype"
	"server/im_group/group_models"
	"server/utils/maps"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUpdateLogic {
	return &GroupUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUpdateLogic) GroupUpdate(req *types.GroupUpdateRequest) (resp *types.GroupUpdateResponse, err error) {
	// 只能是群主或者是管理员才能调用
	//1.判断调用者是否是群成员
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Preload("GroupModel").Take(&groupMember, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("群不存在或用户不是群成员")
	}
	//2.判断调用者是否是群主或管理员
	if !(groupMember.Role == 1 || groupMember.Role == 2) {
		return nil, errors.New("群信息只能是群主或管理员更新")
	}

	//3.将请求req结构体根据标签conf反射成map
	groupMaps := maps.RefToMap(*req, "conf")
	if len(groupMaps) != 0 {
		//对于verificationQuestion发射成结构体
		verificationQuestion, ok := groupMaps["verification_question"]
		if ok {
			delete(groupMaps, "verification_question")
			data := ctype.VerificationQuestion{}
			maps.MapToStruct(verificationQuestion.(map[string]any), &data)
			//更新
			l.svcCtx.DB.Model(&groupMember.GroupModel).Updates(&group_models.GroupModel{
				VerificationQuestion: &data,
			})
		}
		//群信息更新
		err = l.svcCtx.DB.Model(&groupMember.GroupModel).Updates(groupMaps).Error
		if err != nil {
			logx.Error(err)
			return nil, errors.New("群更新失败")
		}
	}

	return
}
