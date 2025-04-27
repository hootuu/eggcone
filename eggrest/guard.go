package eggrest

import (
	"github.com/hootuu/eggcone/database/pgx"
	"github.com/hootuu/eggcone/eggdbx"
	"github.com/hootuu/eggcone/eggdbx/basic"
	"github.com/hootuu/gelato/crtpto/ed25519x"
	"github.com/hootuu/gelato/crtpto/hexx"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/idx"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type GuardM struct {
	basic.Basic

	ID           string    `gorm:"column:id;primaryKey;not null;size:100"`
	Biz          string    `gorm:"column:biz;index;not null;size:100"`
	PubKey       string    `gorm:"column:pub_key;index;not null;size:200"`
	Usage        int64     `gorm:"column:usage;"`
	LstUsageTime time.Time `gorm:"column:lst_usage_time;"`
}

func (model *GuardM) TableName() string {
	return "eggrest_guard"
}

func BuildGuard(biz string) (string, []byte, []byte, *errors.Error) {
	pub, pri, err := ed25519x.NewRandom()
	if err != nil {
		return "", nil, nil, err
	}
	guardM := &GuardM{
		ID:           idx.New(),
		Biz:          biz,
		PubKey:       hexx.Encode(pub),
		Usage:        0,
		LstUsageTime: time.Now(),
	}
	nErr := eggdbx.EggPgDB().Model(guardM).Create(guardM).Error
	if nErr != nil {
		gLogger.Error(
			"eggdbx.EggPgDB().Model(guardM).Create(guardM).Error",
			zap.Error(nErr),
		)
		return "", nil, nil, errors.System("GuardM Create Err", nErr)
	}
	return guardM.ID, pub, pri, nil
}

func GuardExistByBiz(biz string) (bool, *errors.Error) {
	exist, err := pgx.Exist[GuardM](eggdbx.EggPgDB(), "biz = ?", biz)
	if err != nil {
		logger.Error.Error("guard check biz err", zap.Error(err))
		return false, err
	}
	return exist, nil
}

func Guard(id string, callback func(pubKey []byte)) *errors.Error {
	guardM, err := pgx.MustGet[GuardM](eggdbx.EggPgDB(), "id = ?", id)
	if err != nil {
		return err
	}
	bPubKey, err := hexx.Decode(guardM.PubKey)
	if err != nil {
		return err
	}
	callback(bPubKey)
	_ = pgx.Update[GuardM](
		eggdbx.EggPgDB(),
		map[string]interface{}{
			"usage":          gorm.Expr("usage + 1"),
			"lst_usage_time": gorm.Expr("CURRENT_TIMESTAMP"),
		},
		"id = ?",
		id,
	)
	return nil
}

//
//const HeaderGuardID = "guard-id"
//const HeaderGuardRequestID = "guard-request-id"
//const HeaderGuardTimestamp = "guard-timestamp"
//const HeaderGuardNonce = "guard-nonce"
//const HeaderGuardSignature = "guard-signature"
//
//const REQUEST = "egg_request"
//
//type GuardMid struct {
//	guardPrefix string
//}
//
//func NewGuardMid(guardPrefix string) *GuardMid {
//	return &GuardMid{guardPrefix: guardPrefix}
//}
//
//func (g *GuardMid) Handle() gin.HandlerFunc {
//	headerGuardID := fmt.Sprintf("%s-%s", g.guardPrefix, HeaderGuardID)
//	headerGuardRequestID := fmt.Sprintf("%s-%s", g.guardPrefix, HeaderGuardRequestID)
//	headerGuardTimestamp := fmt.Sprintf("%s-%s", g.guardPrefix, HeaderGuardTimestamp)
//	headerGuardNonce := fmt.Sprintf("%s-%s", g.guardPrefix, HeaderGuardNonce)
//	headerGuardSignature := fmt.Sprintf("%s-%s", g.guardPrefix, HeaderGuardSignature)
//	return func(c *gin.Context) {
//		guardRequestID := c.GetHeader(headerGuardRequestID)
//		guardID := c.GetHeader(headerGuardID)
//		guardTimestamp := c.GetHeader(headerGuardTimestamp)
//		guardNonce := c.GetHeader(headerGuardNonce)
//		guardSignature := c.GetHeader(headerGuardSignature)
//		if len(guardRequestID) == 0 {
//			g.badResp(headerGuardRequestID, idx.New(), c)
//			return
//		}
//		if len(guardID) == 0 {
//			g.badResp(headerGuardID, guardRequestID, c)
//			return
//		}
//		if len(guardTimestamp) == 0 {
//			g.badResp(headerGuardTimestamp, guardRequestID, c)
//			return
//		}
//		if len(guardNonce) == 0 {
//			g.badResp(headerGuardNonce, guardRequestID, c)
//			return
//		}
//		if len(guardSignature) == 0 {
//			g.badResp(headerGuardSignature, guardRequestID, c)
//			return
//		}
//		bodyBytes, nErr := c.GetRawData()
//		if nErr != nil {
//			c.JSON(
//				http.StatusBadRequest,
//				rest.FailResponse(guardRequestID, errors.E("902", "load row data bytes failed")),
//			)
//			return
//		}
//		req := rest.OfRequest(guardID, strs.ToInt64(guardTimestamp), strs.ToInt64(guardNonce), guardSignature, bodyBytes)
//		c.Set(REQUEST, req)
//		c.Next()
//	}
//}
//
//func (g *GuardMid) badResp(headerKey string, reqID string, c *gin.Context) {
//	c.JSON(http.StatusBadRequest, rest.FailResponse(reqID, errors.E("901", "require %s", headerKey)))
//}
