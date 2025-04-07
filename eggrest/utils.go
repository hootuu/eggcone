package eggrest

import (
	"github.com/gin-gonic/gin"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/idx"
	"github.com/hootuu/gelato/io/rest"
	"github.com/hootuu/gelato/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
)

var gLogger = logger.GetLogger("eggrest")

func Handle[REQ any, RESP any](ctx *gin.Context, callback func(req *rest.Request[REQ]) *rest.Response[RESP]) {
	bodyBytes, nErr := ctx.GetRawData()
	if nErr != nil {
		gLogger.Error("get request body data failed", zap.Error(nErr))
		ctx.JSON(
			http.StatusOK,
			rest.FailResponse[any](idx.New(), errors.System("get request body data failed")),
		)
		return
	}
	req, err := rest.UnmarshalRequest[REQ](bodyBytes)
	if err != nil {
		gLogger.Error("unmarshal request failed", zap.Error(err.Native()))
		ctx.JSON(
			http.StatusOK,
			rest.FailResponse[RESP](idx.New(), errors.System("unmarshal request failed")),
		)
		return
	}
	bReqVerified := true
	err = Guard(req.GuardID, func(pubKey []byte) {
		innerErr := req.Verify(pubKey)
		if innerErr != nil {
			gLogger.Error("guard.Verify failed", zap.String("req", req.JSON()), zap.Error(innerErr))
			bReqVerified = false
		}
	})
	if err != nil {
		gLogger.Error("guard failed", zap.Error(err.Native()))
		ctx.JSON(
			http.StatusOK,
			rest.FailResponse[RESP](req.ID, errors.System("guard failed")),
		)
		return
	}
	if !bReqVerified {
		ctx.JSON(
			http.StatusOK,
			rest.FailResponse[RESP](req.ID, errors.System("invalid signature")),
		)
		return
	}
	resp := callback(req)
	if gLogger.Level() <= zapcore.InfoLevel {
		gLogger.Info("["+ctx.Request.Method+"]",
			zap.String("URL", ctx.Request.URL.String()),
			zap.Any("req", req),
			zap.Any("resp", resp))
	}
	ctx.JSON(http.StatusOK, resp)
}
