package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type entityDao struct {
}

var EntityDao *entityDao

func newEntityDao() *entityDao {
	return &entityDao{}
}

func init() {
	EntityDao = newEntityDao()
}

func (entity *entityDao) Create(newEntity model.Entity) error {
	result := db.Model(&model.Entity{}).Create(&newEntity)
	return utils.DBError(result)
}

func (entity *entityDao) Delete(id []uint) error {
	result := db.Model(&model.Entity{}).Where("id in (?)", id).Delete(&model.Entity{})
	return utils.DBError(result)
}

func (entity *entityDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Entity{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (entity *entityDao) AllEntity() (list []model.Entity, err error) {
	result := db.Model(&model.Entity{}).Find(&list)
	err = utils.DBError(result)
	return
}

func (entity *entityDao) GetEntityByName(name string) (*model.Entity, error) {
	ret := &model.Entity{}
	result := db.Model(&model.Entity{}).Where("name = ?", name).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (entity *entityDao) GetEntityByID(id uint) (*model.Entity, error) {
	ret := &model.Entity{}
	result := db.Model(&model.Entity{}).Where("id = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (entity *entityDao) GetEntitysByNames(name []string) (list []model.Entity, err error) {
	result := db.Model(&model.Entity{}).Where("name IN (?)", name).Order("id").Find(&list)
	err = utils.DBError(result)
	return
}

func (entity *entityDao) EntityCount() (count int64, err error) {
	result := db.Model(&model.Entity{}).Count(&count)
	err = utils.DBError(result)
	return
}

// entity and user
func (entity *entityDao) GetEntityAllUser(id uint, offset int, limit int) (users []*model.User, count int64, err error) {
	// query_entity, err := entity.GetEntityByName(name)
	// if err != nil {
	// 	return
	// }
	err = utils.DBError(db.Model(&model.User{}).Preload("Department").Preload("Entity").Where("entity_id = ?", id).Count(&count).Offset(offset).Limit(limit).Find(&users))
	return
}

func (entity *entityDao) GetEntityManager(id uint) (managers []*model.User, err error) {
	err = utils.DBError(db.Model(&model.User{}).Where("entity_id = ? and entity_super = ?", id, true).Find(&managers))
	return
}

// entity and department
func (entity *entityDao) GetEntityAllDepartment(id uint) (departments []*model.Department, err error) {
	// query_entity, err := entity.GetEntityByName(name)
	// if err != nil {
	// 	return
	// }
	err = utils.DBError(db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("entity_id = ?", id).Find(&departments))
	return
}

func (entity *entityDao) GetEntitySubDepartment(name string) (departments []*model.Department, err error) {
	query_entity, err := entity.GetEntityByName(name)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("entity_id = ? and parent_id IS NULL", query_entity.ID, 0).Find(&departments))
	return
}

func (entity *entityDao) GetEntitySubDepartmentByID(id uint) (departments []*model.Department, err error) {
	result := db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("entity_id = ? and parent_id IS NULL", id).Find(&departments)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
		departments = nil
		return
	}
	err = utils.DBError(result)
	return
}
