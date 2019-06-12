package mysql

import (
	"github.com/jinzhu/gorm"
)

// 获取offset
func GetOffset(page int, perPage int) int {
	if page < 1 {
		page = 1
	}
	return (page - 1) * perPage
}

// 获取分页页数总数及数据列表。m是查询的表的model，result是列表结果传指针进来！
func GetDataByPageAndPerPage(db *gorm.DB, page int, perPage int, m interface{}, result interface{}) (totalPages int, totalCount int) {
	offset := GetOffset(page, perPage)
	if err := db.Offset(offset).Limit(perPage).Find(result).Error; err != nil {
		return
	}

	db = db.Offset(-1).Limit(-1)
	if err := db.Model(m).Count(&totalCount).Error; err != nil {
		return
	}
	if perPage <= 0 {
		return
	}
	totalPages = (totalCount + perPage -1)/perPage
	return
}