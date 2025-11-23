package po

import "time"

type Chapter struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TagID     uint      `json:"tag_id" gorm:"index"`
	Tag       *Tag      `json:"tag,omitempty" gorm:"foreignKey:TagID"`
	Name      string    `json:"name" gorm:"type:varchar(200);not null"`
	Sort      int       `json:"sort" gorm:"default:0;comment:排序,数字越小越靠前"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Chapter) TableName() string {
	return "chapters"
}
