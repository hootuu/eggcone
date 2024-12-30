package once

import (
	"fmt"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/hootuu/gelato/sys"
	"go.uber.org/zap"
	"time"
)

func Do(code string, doFunc func() *errors.Error) *errors.Error {
	m, err := Get(code)
	if err != nil {
		return err
	}

	if m != nil {
		if m.DoStatus == FAILED {
			return errors.System(fmt.Sprintf("a single task failed to execute on other machines: %s", code))
		}
		return nil
	}

	m = &Once{
		Code:        code,
		DoServerID:  sys.ServerID,
		DoStatus:    EXECUTING,
		DoStartTime: time.Now(),
	}

	err = Create(m)
	if err != nil {
		return err
	}

	err = doFunc()
	if err != nil {
		logger.Logger.Error("a single task execute failed", zap.String("code", code), zap.Error(err))
		igErr := SetEnd(m, FAILED)
		if igErr != nil {
			logger.Logger.Error("a single task execute failed, and sync task status failed", zap.String("code", code), zap.Error(igErr))
		}
		return err
	}

	igErr := SetEnd(m, SUCCESS)
	if igErr != nil {
		logger.Logger.Error("a single task execute failed, and sync task status failed", zap.String("code", code), zap.Error(igErr))
	}

	return nil
}
