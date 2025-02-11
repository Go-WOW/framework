package dbRepository

import (
	"errors"
	"stncCms/app/domain/entity"
	"strings"

	"github.com/jinzhu/gorm"
)

type ModulesRepo struct {
	db *gorm.DB
}

func ModulesRepositoryInit(db *gorm.DB) *ModulesRepo {
	return &ModulesRepo{db}
}

//ModulesRepo implements the repository.ModulesRepository interface
// var _ services.ModulesAppInterface = &ModulesRepo{}

// GetAll all data
func (r *ModulesRepo) GetAll() ([]entity.Modules, error) {
	var datas []entity.Modules
	var err error
	err = r.db.Debug().Order("created_at desc").Find(&datas).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

// GetAll all data
func (r *ModulesRepo) GetAllModulesMerge() ([]entity.ModulesAndPermissionDTO, error) {
	var err error
	var datas []entity.ModulesAndPermissionDTO
	err = r.db.Debug().Table("modules").Order("created_at desc").Find(&datas).Error

	//TODO: nasil preload yapilir bakilacak
	// var datas []entity.Modules
	// err = r.db.Debug().Preload("Permission").Take(&datas).Error

	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

// GetAll all data
func (r *ModulesRepo) GetAllModulesMergePermission() ([]entity.ModulesAndPermissionRoleDTO, error) {
	var err error
	var datas []entity.ModulesAndPermissionRoleDTO
	err = r.db.Debug().Table("modules").Order("created_at desc").Find(&datas).Error

	//TODO: nasil preload yapilir bakilacak
	// var datas []entity.Modules
	// err = r.db.Debug().Preload("Permission").Take(&datas).Error

	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("data not found")
	}
	return datas, nil
}

// Save data
func (r *ModulesRepo) Save(cat *entity.Modules) (*entity.Modules, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&cat).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "post title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return cat, nil
}

// Update upate data
func (r *ModulesRepo) Update(cat *entity.Modules) (*entity.Modules, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&cat).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return cat, nil
}

// GetByID get data
func (r *ModulesRepo) GetByID(id uint64) (*entity.Modules, error) {
	var cat entity.Modules
	err := r.db.Debug().Where("id = ?", id).Take(&cat).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")

	}
	return &cat, nil
}

// GetAllPaginate paginatin list
func (r *ModulesRepo) GetAllPaginate(postsPerPage int, offset int) ([]entity.Modules, error) {
	var regionList []entity.Modules
	err := r.db.Debug().Order("id asc").Limit(postsPerPage).Offset(offset).Find(&regionList).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("post not found")
	}
	return regionList, nil
}

// GetAllPaginateCount
func (r *ModulesRepo) GetAllPaginateCount(returnValue *int64) {
	var table entity.Modules
	var count int64
	r.db.Debug().Model(table).Count(&count)
	*returnValue = count
}

// Delete delete data
func (r *ModulesRepo) Delete(id uint64) error {
	var cat entity.Modules
	err := r.db.Debug().Where("id = ?", id).Delete(&cat).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
