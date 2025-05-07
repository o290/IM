package group_models

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"server/common/models"
	"time"
)

// 群成员表
type GroupMemberModel struct {
	models.Model
	GroupID         uint       `json:"groupID"`
	GroupModel      GroupModel `gorm:"foreignKey:GroupID" json:"-"`
	UserID          uint       `json:"userID"`
	MemberNickname  string     `gorm:"size:32" json:"memberNickname"` //群昵称
	Role            int8       `json:"role"`                          //1：群主2：管理员3：普通成员
	ProhibitionTime *int       `json:"prohibitionTime"`               //禁言时长，单位分钟
}

// GetProhibitionTime 获取群成员禁言时间
func (gm GroupMemberModel) GetProhibitionTime(client *redis.Client, db *gorm.DB) *int {
	//没有被禁言
	if gm.ProhibitionTime == nil {
		return nil
	}
	//获取剩余时间
	t, err := client.TTL(context.Background(), fmt.Sprintf("prohibition__%d", gm.ID)).Result()
	if err != nil {
		//查不到就说明过期了，就把这个值改为nil，这个是查询到redis过期后再修改
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	//检查是否已经过期
	if t == -2*time.Nanosecond {
		db.Model(&gm).Update("prohibition_time", nil)
		return nil
	}
	//计算禁言的分钟数
	res := int(t / time.Minute)
	return &res
}
