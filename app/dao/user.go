package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"errors"

	"gorm.io/gorm"
)

type userDao struct {
}

var UserDao *userDao

func newUserDao() *userDao {
	return &userDao{}
}

func init() {
	UserDao = newUserDao()
}

func (user *userDao) Create(newUser model.User) error {
	result := db.Model(&model.User{}).Create(&newUser)
	return utils.DBError(result)
}

func (user *userDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.User{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (user *userDao) UpdateByName(username string, data map[string]interface{}) error {
	result := db.Model(&model.User{}).Where("username = ?", username).Updates(data)
	return utils.DBError(result)
}

func (user *userDao) Delete(id []uint) error {
	result := db.Model(&model.User{}).Where("id in (?)", id).Delete(&model.User{})
	return utils.DBError(result)
}

func (user *userDao) AllUser(Offset int, Limit int) (list []*model.User, count int64, err error) {
	result := db.Model(&model.User{}).Preload("Department").Preload("Entity").Count(&count).Offset(Offset).Limit(Limit).Find(&list)
	err = utils.DBError(result)
	return
}

// 预计给分页器使用
func (user *userDao) GetLimitUser(begin int, length int) (list []model.User, err error) {
	if begin <= 0 || length <= 0 {
		err = errors.New("invalid number")
		return
	}
	result := db.Model(&model.User{}).Find(&list).Offset(begin - 1).Limit(length)
	err = utils.DBError(result)
	return
}

func (user *userDao) GetUserByName(username string) (*model.User, error) {
	ret := &model.User{}
	result := db.Model(&model.User{}).Preload("Department").Preload("Entity").Where("username = ?", username).First(ret)
	// department := model.Department{}
	// err := db.Model(&ret).Association("Department").Find(&department)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.Department = &department
	// entity := &model.Entity{}
	// err = db.Model(&ret).Association("Entity").Find(&entity)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.Entity = entity
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (user *userDao) GetUserByID(id uint) (*model.User, error) {
	ret := &model.User{}
	result := db.Model(&model.User{}).Preload("Department").Preload("Entity").Where("id = ?", id).First(ret)
	// department := model.Department{}
	// err := db.Model(&ret).Association("Department").Find(&department)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.Department = &department
	// entity := &model.Entity{}
	// err = db.Model(&ret).Association("Entity").Find(&entity)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.Entity = entity
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (user *userDao) GetUsersByNames(username []string) (list []model.User, err error) {
	result := db.Model(&model.User{}).Where("username IN ?", username).Order("id").Find(&list)
	err = utils.DBError(result)
	return
}

func (user *userDao) UserCount() (count int64, err error) {
	result := db.Model(&model.User{}).Count(&count)
	err = utils.DBError(result)
	return
}

func (user *userDao) ModifyUserIdentity(username string, identity int) error {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	if identity == 0 {
		err = user.Update(thisUser.ID, map[string]interface{}{
			"system_super":     false,
			"entity_super":     false,
			"department_super": false,
		})
	} else if identity == 1 {
		err = user.Update(thisUser.ID, map[string]interface{}{
			"department_super": true,
		})
	} else if identity == 2 {
		err = user.Update(thisUser.ID, map[string]interface{}{
			"entity_super": true,
		})
	} else if identity == 3 {
		err = user.Update(thisUser.ID, map[string]interface{}{
			"system_super": true,
		})
	} else {
		err = errors.New("invalid identity number")
	}
	return err
}

func (user *userDao) ModifyUserPassword(username string, password string) error {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	err = user.Update(thisUser.ID, map[string]interface{}{
		"password": password,
	})
	return err
}

func (user *userDao) ModifyUserBanstate(username string, ban bool) error {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	err = user.Update(thisUser.ID, map[string]interface{}{
		"ban": ban,
	})
	return err
}

// User Entity Part
func (user *userDao) GetUserEntity(username string) (entity model.Entity, err error) {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return
	}
	if thisUser == nil {
		err = errors.New("user doesn't exist")
		return
	}
	entity = *thisUser.Entity
	//db.Model(&thisUser).Where("id = ?", thisUser.ID).Preload("entity").Find(&entity)
	return
}

func (user *userDao) ModifyUserEntity(username string, entity model.Entity) error {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	thisUser.EntityID = entity.ID
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisUser))
	return err
}

func (user *userDao) ModifyUserEntityByID(id uint, entityID uint) error {
	thisUser, err := user.GetUserByID(id)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	err = utils.DBError(db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"entity_id": entityID,
	}))
	return err
}

// User Department Part
func (user *userDao) GetUserDepartment(username string) (department model.Department, err error) {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return
	}
	if thisUser == nil {
		err = errors.New("user doesn't exist")
		return
	}
	department = *thisUser.Department
	//err = db.Model(&thisUser).Association("Department").Find(&department)
	return
}

func (user *userDao) ModifyUserDepartment(username string, department model.Department) error {
	thisUser, err := user.GetUserByName(username)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	thisUser.DepartmentID = department.ID
	if department.EntityID != 0 {
		thisUser.EntityID = department.Entity.ID
	}
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisUser))
	//err = utils.DBError(db.Save(&thisUser))
	return err
}

func (user *userDao) ModifyUserDepartmentByID(id uint, departmentID uint) error {
	thisUser, err := user.GetUserByID(id)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	err = utils.DBError(db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"department_id": departmentID,
	}))
	return err
}

// feishu
func (user *userDao) BindFeishu(UserID uint, FeishuID string) error {
	thisUser, err := user.GetUserByID(UserID)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	err = utils.DBError(db.Model(&model.User{}).Where("id = ?", UserID).Updates(map[string]interface{}{
		"feishu_id": FeishuID,
	}))
	return err
}

func (user *userDao) UpdateFeishuToken(UserID uint, FeishuToken string, RefreshToken string) error {
	thisUser, err := user.GetUserByID(UserID)
	if err != nil {
		return err
	}
	if thisUser == nil {
		return errors.New("user doesn't exist")
	}
	thisUser.FeishuToken = FeishuToken
	thisUser.RefreshToken = RefreshToken
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisUser))
	return err
}

func (user *userDao) GetUserByFeishuID(FeishuID string) (*model.User, error) {
	ret := &model.User{}
	result := db.Model(&model.User{}).Preload("Department").Preload("Entity").Where("feishu_id = ?", FeishuID).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

/*func (user *userDao) GetFeishuTokenByID(UserID uint) (token string, err error) {
	thisUser, err := user.GetUserByID(UserID)
	if err != nil {
		return
	}
	if thisUser == nil {
		err = errors.New("user doesn't exist")
		return
	}
	token = thisUser.FeishuToken
	return
}

func (user *userDao) GetRefreshTokenByID(UserID uint) (token string, err error) {
	thisUser, err := user.GetUserByID(UserID)
	if err != nil {
		return
	}
	if thisUser == nil {
		err = errors.New("user doesn't exist")
		return
	}
	token = thisUser.RefreshToken
	return
}*/
