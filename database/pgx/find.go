package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/io/pagination"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Exist[T any](db *gorm.DB, query interface{}, args ...interface{}) (bool, *errors.Error) {
	var model T
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

func Get[T any](db *gorm.DB, cond ...interface{}) (*T, *errors.Error) {
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

func MustGet[T any](db *gorm.DB, cond ...interface{}) (*T, *errors.Error) {
	model, err := Get[T](db, cond...)
	if err != nil {
		return nil, err
	}
	if model == nil {
		return nil, errors.Verify("no such record")
	}
	return model, nil
}

func PagedOrderFind[T any](
	db *gorm.DB,
	page *pagination.Page,
	order interface{},
	query interface{},
	cond ...interface{},
) (*[]*T, *pagination.Paging, *errors.Error) {
	var md T
	var models []*T
	var count int64
	tx := db.Model(&md).Order(order).Where(query, cond...).Count(&count)
	if tx.Error != nil {
		return nil, nil, errors.System("system error", tx.Error)
	}
	tx = db.Where(query, cond...).Limit(int(page.Size)).Offset(int((page.Numb - 1) * page.Size)).Find(&models)
	if tx.Error != nil {
		return nil, nil, errors.System("system error", tx.Error)
	}
	return &models, pagination.PagingOfPage(page).WithCount(count), nil
}
