package pgx

import (
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/io/pagination"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func PagedOrderQuery[T any](
	db *gorm.DB,
	page *pagination.Page,
	order interface{},
	query interface{},
	cond ...interface{},
) (*pagination.Pagination[T], *errors.Error) {
	if page == nil {
		page = pagination.PageNormal()
	}
	var md T
	var arr []*T
	var count int64
	tx := db.Model(&md).Order(order).Where(query, cond...).Count(&count)
	if tx.Error != nil {
		logger.Logger.Error("PagedOrderQuery.Count err", zap.Error(tx.Error))
		return nil, errors.System("system error", tx.Error)
	}
	tx = db.Where(query, cond...).Limit(int(page.Size)).Offset(int((page.Numb - 1) * page.Size)).Find(&arr)
	if tx.Error != nil {
		logger.Logger.Error("PagedOrderQuery.Find err", zap.Error(tx.Error))
		return nil, errors.System("system error", tx.Error)
	}
	return pagination.NewPagination[T](pagination.PagingOfPage(page).WithCount(count), arr), nil
}
