package po

import "time"

type Chapter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TagID     uint      `json:"tag_id" gorm:"index"`
	Tag       *Tag      `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	ParentID  *uint     `json:"parent_id" gorm:"index;comment:父章节ID,为空表示一级章节"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null"`
	Sort      int       `json:"sort" gorm:"default:0;comment:排序,数字越小越靠前"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Chapter) TableName() string {
	return "chapters"
}
