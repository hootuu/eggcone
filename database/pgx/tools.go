package pgx

import (
	"fmt"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/io/pagination"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func PgExists(db *gorm.DB, model interface{}, query interface{}, args ...interface{}) (bool, *errors.Error) {
	var exists bool
	err := db.Model(model).
		Select("1").
		Where(query, args...).
		Limit(1).
		Find(&exists).
		Error
	if err != nil {
		return false, errors.System("db error", err)
	}
	return exists, nil
}

func PgGet[T any](db *gorm.DB, cond ...interface{}) (*T, *errors.Error) {
	var model T
	tx := db.First(&model, cond...)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Logger.Error("system error", zap.Any("cond", cond), zap.Error(tx.Error))
		return nil, errors.System("db err", tx.Error)
	}
	return &model, nil
}

func PgLoad[T any](db *gorm.DB, cond ...interface{}) (*T, *errors.Error) {
	md, err := PgGet[T](db, cond...)
	if err != nil {
		return nil, err
	}
	if md == nil {
		return nil, errors.E("no_such", fmt.Sprintf("no such data %v", cond))
	}
	return md, nil
}

func PgPageFind[T any](db *gorm.DB, page *pagination.Page, query interface{}, cond ...interface{}) (*[]*T, *pagination.Paging, *errors.Error) {
	var md T
	var models []*T
	var count int64
	tx := db.Model(&md).Where(query, cond...).Count(&count)
	if tx.Error != nil {
		return nil, nil, errors.System("system error", tx.Error)
	}
	tx = db.Where(query, cond...).Limit(int(page.Size)).Offset(int((page.Numb - 1) * page.Size)).Find(&models)
	if tx.Error != nil {
		return nil, nil, errors.System("system error", tx.Error)
	}
	return &models, pagination.PagingOfPage(page).WithCount(count), nil
}

func PgFind[T any](db *gorm.DB, query interface{}, cond ...interface{}) ([]*T, *errors.Error) {
	var models []*T
	tx := db.Where(query, cond...).Find(&models)
	if tx.Error != nil {
		return nil, errors.System("db error", tx.Error)
	}
	return models, nil
}
