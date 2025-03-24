package eggmq

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/eggcone/eggmq/modelx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/io/pagination"
	"go.uber.org/zap"
	"time"
)

func MessageConvertTo(model *modelx.MessageM) *Message {
	return &Message{
		ID:         model.ID,
		Topic:      model.Topic,
		Payload:    model.Payload,
		RetryCount: model.RetryCount,
		Status:     model.Status,
	}
}

func MessageConvertFrom(msg *Message) *modelx.MessageM {
	return &modelx.MessageM{
		ID:         msg.ID,
		Topic:      msg.Topic,
		Payload:    msg.Payload,
		RetryCount: 0,
		Status:     msg.Status,
	}
}

func MessageCreate(msg *Message) *errors.Error {
	bExist, err := MessageExist(msg.ID)
	if err != nil {
		return err
	}
	if bExist {
		return errors.Verify("duplicate message id")
	}
	model := MessageConvertFrom(msg)
	tx := eggdbx.Egg().DB().Create(&model)
	if tx.Error != nil {
		eggdbx.Logger.Error("create eggmq.message failed", zap.Error(tx.Error), zap.Any("paras", model))
		return errors.System("db err", tx.Error)
	}
	return nil
}

func MessageToProcessing(msg *Message) *errors.Error {
	return pgx.Update[modelx.MessageM](
		eggdbx.EggPgDB(),
		map[string]interface{}{
			"status": PROCESSING,
		},
		"id = ?",
		msg.ID,
	)
}

func MessageToPending(msg *Message) *errors.Error {
	return pgx.Update[modelx.MessageM](
		eggdbx.EggPgDB(),
		map[string]interface{}{
			"status":      PENDING,
			"retry_count": msg.RetryCount + 1,
		},
		"id = ?",
		msg.ID,
	)
}

func MessageToFailed(msg *Message) *errors.Error {
	return pgx.Update[modelx.MessageM](
		eggdbx.EggPgDB(),
		map[string]interface{}{
			"status":      FAILED,
			"retry_count": msg.RetryCount + 1,
		},
		"id = ?",
		msg.ID,
	)
}

func MessageToCompleted(msg *Message) *errors.Error {
	return pgx.Update[modelx.MessageM](
		eggdbx.EggPgDB(),
		map[string]interface{}{
			"status": COMPLETED,
		},
		"id = ?",
		msg.ID,
	)
}

func MessageExist(id string) (bool, *errors.Error) {
	return pgx.Exist[modelx.MessageM](
		eggdbx.Egg().DB(),
		"id = ?",
		id,
	)
}

func MessageLoadPending(callback func(msg *Message)) *errors.Error {
	page := &pagination.Page{
		Size: 100,
		Numb: pagination.FirstPage,
	}

	processingInTime := time.Now().Add(0 - time.Duration(10*time.Minute))
	for {
		arr, paging, err := pgx.PagedOrderFind[modelx.MessageM](
			eggdbx.EggPgDB(),
			page,
			"auto_id DESC",
			"status IN (?) OR updated_at < ?",
			[]int{PENDING, FAILED, PROCESSING},
			processingInTime,
		)
		if err != nil {
			return err
		}
		if paging.Count == 0 {
			break
		}
		for _, model := range *arr {
			if model.Status == PROCESSING {
				if time.Now().Sub(model.UpdatedAt).Abs() < 10*time.Minute {
					continue
				}
			}
			msg := MessageConvertTo(model)
			callback(msg)
		}
		page.Numb = paging.Numb + 1
	}
	return nil
}
