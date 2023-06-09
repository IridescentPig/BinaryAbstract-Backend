package model

import (
	"gorm.io/datatypes"
)

type Department struct {
	ID        uint                        `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"department_id"`
	Name      string                      `gorm:"column:name;not null" json:"department_name"`
	ParentID  uint                        `gorm:"default:null;column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent    *Department                 `gorm:"foreignKey:ParentID;references:ID;default:null" json:"-"`
	EntityID  uint                        `gorm:"default:null;column:entity_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"entity_id"`
	Entity    Entity                      `gorm:"foreignKey:EntityID;references:ID;default:null" json:"-"`
	KeyList   datatypes.JSONSlice[string] `gorm:"column:key_list" json:"key_list"`
	Threshold uint                        `gorm:"column:threshold" json:"threshold"`
}

// func (Department) TableName() string {
// 	return "department"
// }
