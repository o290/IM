package list_query

import (
	"fmt"
	"gorm.io/gorm"
	"server/common/models"
)

// Option 定义一组用于数据库查询的选项
type Option struct {
	PageInfo models.PageInfo      //分页查询
	Where    *gorm.DB             //高级查询
	Joins    string               //联合查询
	Likes    []string             //模糊匹配的字段
	Preload  []string             //预加载字段
	Table    func() (string, any) //用于指定表明,子查询
	Groups   []string             //分组
}

// ListQuery 泛型类函数
// model 要查询的数据库对象 option 封装的各种可选配置信息,相当于select后面的内容
func ListQuery[T any](db *gorm.DB, model T, option Option) (list []T, count int64, err error) {
	//条件查询
	query := db.Where(model) //把结构体自己的查询条件查了

	//模糊匹配
	//fmt.Println("lisr", option.Where)
	if option.PageInfo.Key != "" && len(option.Likes) > 0 {
		likeQuery := db.Where("")
		//遍历 option.Likes 切片，构建 LIKE 查询条件，多个条件之间使用 OR 连接
		for index, like := range option.Likes {
			if index == 0 {
				//where name like '%xxxx%'
				likeQuery.Where(fmt.Sprintf("%s like '%%?%%'", like), option.PageInfo.Key)
			} else {
				likeQuery.Or(fmt.Sprintf("%s like '%%?%%'", like), option.PageInfo.Key)
			}
		}
		query.Where(likeQuery)
	}

	if option.Table != nil {
		table, data := option.Table()
		query = query.Table(table, data)
	}

	if option.Joins != "" {
		query = query.Joins(option.Joins)
	}

	if option.Where != nil {
		query = query.Where(option.Where)
	}
	if len(option.Groups) > 0 {
		for _, group := range option.Groups {
			query = query.Group(group)
		}
	}
	//求总数
	query.Model(model).Count(&count)
	//预加载
	for _, s := range option.Preload {
		query = query.Preload(s)
	}

	//分页查询
	if option.PageInfo.Page <= 0 {
		option.PageInfo.Page = 1
	}
	if option.PageInfo.Limit != -1 { //-1查全部
		if option.PageInfo.Limit <= 0 {
			option.PageInfo.Limit = 10
		}
	}

	//计算偏移量
	offset := (option.PageInfo.Page - 1) * option.PageInfo.Limit

	if option.PageInfo.Sort != "" {
		query.Order(option.PageInfo.Sort)
	}
	err = query.Limit(option.PageInfo.Limit).Offset(offset).Find(&list).Error
	return
}
