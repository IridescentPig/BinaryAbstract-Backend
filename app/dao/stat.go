package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
)

type statDao struct {
}

var StatDao *statDao

func newStatDao() *statDao {
	return &statDao{}
}

func init() {
	StatDao = newStatDao()
}

func (stat *statDao) GetAllAssetStat() ([]*model.Stat, error) {
	var stats []*model.Stat

	result := db.Model(model.Asset{}).Select("department_id, SUM(net_worth) as total").Group("department_id").Scan(&stats)
	if result.Error != nil {
		return nil, utils.DBError(result)
	}

	return stats, nil
}

func (stat *statDao) CreateAssetStats(stats []*model.Stat) error {
	result := db.Model(&model.Stat{}).Create(stats)

	return utils.DBError(result)
}

func (stat *statDao) GetDepartmentStat(departmentID uint) (stats []*model.Stat, err error) {
	result := db.Model(&model.Stat{}).Where("department_id = ?", departmentID).Find(&stats)

	err = utils.DBError(result)
	return stats, err
}
