package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type assetService struct{}

var AssetService *assetService

func newAssetService() *assetService {
	return &assetService{}
}

func init() {
	AssetService = newAssetService()
}

func (asset *assetService) GetSubAsset(parentID uint, departmentID uint) ([]*define.AssetInfo, error) {
	var subAssetList []*model.Asset
	var err error

	if parentID == 0 {
		subAssetList, err = dao.AssetDao.GetAssetDirectDepartment(departmentID)
	} else {
		subAssetList, err = dao.AssetDao.GetSubAsset(parentID)
	}

	if err != nil {
		return nil, err
	}

	subAssetTreeNodeList := []*define.AssetInfo{}
	err = copier.Copy(&subAssetTreeNodeList, subAssetList)
	if err != nil {
		return nil, err
	}

	for _, subNode := range subAssetTreeNodeList {
		subNode.Children, err = asset.GetSubAsset(subNode.AssetID, departmentID)
		if err != nil {
			return nil, err
		}
	}

	return subAssetTreeNodeList, err
}

func (asset *assetService) GetAssetByID(assetID uint) (*model.Asset, error) {
	return dao.AssetDao.GetAssetByID(assetID)
}

func (asset *assetService) ExistAsset(assetID uint) (bool, error) {
	thisAsset, err := dao.AssetDao.GetAssetByID(assetID)
	if err != nil {
		return false, err
	}
	return thisAsset != nil, nil
}

func (asset *assetService) CheckAssetInDepartment(assetID uint, departmentID uint) (bool, error) {
	thisAsset, err := dao.AssetDao.GetAssetByID(assetID)
	if err != nil || thisAsset == nil {
		return false, nil
	}
	return thisAsset.DepartmentID == departmentID, nil
}

func (asset *assetService) CheckIsAncestor(srcID uint, targetID uint) (bool, error) {
	targetAsset, err := dao.AssetDao.GetAssetByID(targetID)
	if err != nil {
		return true, err
	}

	flag := false
	for {
		if targetAsset == nil {
			break
		}
		if targetAsset.ID == srcID {
			flag = true
			break
		}
		targetAsset, err = dao.AssetDao.GetAssetByID(targetAsset.ParentID)
		if err != nil {
			return true, err
		}
	}

	return flag, nil
}

func (asset *assetService) ModifyAssetInfo(id uint, req define.ModifyAssetInfoReq) error {
	err := dao.AssetDao.UpdateByStruct(id, model.Asset{
		Name:        req.AssetName,
		Price:       req.Price,
		Description: req.Description,
		Position:    req.Position,
		ClassID:     req.ClassID,
		Type:        req.Type,
		Number:      req.Number,
	})
	if err != nil {
		return err
	}
	if req.ParentID != nil {
		if *req.ParentID != 0 {
			err = dao.AssetDao.Update(id, map[string]interface{}{
				"parent_id": *req.ParentID,
			})
		} else {
			err = dao.AssetDao.Update(id, map[string]interface{}{
				"parent_id": gorm.Expr("NULL"),
			})
		}
	}
	return err
}

func (asset *assetService) CreateAsset(req *define.CreateAssetReq, departmentID uint, parentID uint, userID uint) error {
	thisID, err := dao.AssetDao.CreateAndGetID(model.Asset{
		Name:         req.AssetName,
		Price:        req.Price,
		Description:  req.Description,
		Position:     req.Position,
		ClassID:      req.ClassID,
		Number:       req.Number,
		Type:         req.Type,
		DepartmentID: departmentID,
		UserID:       userID,
		ParentID:     parentID,
	})
	if err != nil {
		return err
	}
	for _, child := range req.Children {
		err = asset.CreateAsset(child, departmentID, thisID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (asset *assetService) ExpireAssets(assetIDs []uint) error {
	err := dao.AssetDao.ExpireAsset(assetIDs)
	return err
}

func (asset *assetService) TransferAssets(assetIDs []uint, userID uint, departmentID uint) error {
	subAssets, err := dao.AssetDao.GetSubAssetsByParents(assetIDs)
	if err != nil {
		return err
	}
	subAssetIDs := []uint{}

	for _, asset := range subAssets {
		subAssetIDs = append(subAssetIDs, asset.ID)
	}

	err = dao.AssetDao.AllUpdate(subAssetIDs, map[string]interface{}{
		"parent_id": gorm.Expr("NULL"),
	})
	if err != nil {
		return err
	}
	err = dao.AssetDao.AllUpdate(assetIDs, map[string]interface{}{
		"user_id":       userID,
		"parent_id":     gorm.Expr("NULL"),
		"department_id": departmentID,
	})
	return err
}

func (asset *assetService) GetAssetByUser(user_id uint) (assets []*define.AssetInfo, err error) {
	assetList, err := dao.AssetDao.GetAssetsByUser(user_id)
	if err != nil {
		return
	}
	err = copier.Copy(&assets, assetList)
	if err != nil {
		return
	}
	return
}

func (asset *assetService) GetDepartmentAssetsByIDs(ids []uint, departmentID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetDepartmentAssetsByIDs(ids, departmentID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) GetUserAssetsByIDs(ids []uint, userID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetUserAssetsByIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) GetDepartmentIdleAssets(ids []uint, departmentID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetDepartmentIdleAssetsByIDs(ids, departmentID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) AcquireAssets(ids []uint, userID uint) error {
	err := dao.AssetDao.ModifyAssetsUserAndState(ids, userID, 1)
	return err
}

func (asset *assetService) CancelAssets(ids []uint, userID uint) error {
	err := dao.AssetDao.ModifyAssetsUserAndState(ids, userID, 0)
	return err
}

func (asset *assetService) GetUserMaintainAssets(userID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetUserMaintainAssets(userID)
	return assetList, err
}

func (asset *assetService) ModifyAssetMaintainerAndState(assetIDs []uint, maintainerID uint) error {
	err := dao.AssetDao.ModifyAssetMaintainerAndState(assetIDs, maintainerID)
	return err
}
