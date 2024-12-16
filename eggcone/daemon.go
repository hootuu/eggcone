package eggcone

import (
	"context"
	"fmt"
	"github.com/hootuu/gelato/configure"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/logger"
	"github.com/hootuu/gelato/sys"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"strings"
	"time"
)

type Daemon interface {
	Code() string
	Startup() *errors.Error
	Shutdown(ctx context.Context) *errors.Error
}

var gDaemons map[string]Daemon

func doRegister(daemon Daemon) *errors.Error {
	_, exists := gDaemons[daemon.Code()]
	if exists {
		return errors.System(fmt.Sprintf("The same daemon has been register.[ code=%s ]", daemon.Code()))
	}
	gDaemons[daemon.Code()] = daemon
	return nil
}

func Register(daemon ...Daemon) *errors.Error {
	if len(daemon) == 0 {
		return errors.System("no any daemon")
	}
	for _, s := range daemon {
		err := doRegister(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func Startup() {
	sys.Info("# Startup all registered daemons ...... #")
	for code, daemon := range gDaemons {
		sys.Info("## Startup the daemon: ", code, " ...... #")
		err := daemon.Startup()
		if err != nil {
			logger.Logger.Error("start daemon failed", zap.String("code", code), zap.Error(err))
			sys.Error("# Start daemon exception: ", code, " #")
			return
		}
		sys.Success("## Startup the daemon ", code, " [OK] #")
	}

	sys.Info("## Display All Used Configure Items ...... #")
	configure.Dump(func(key string, val any) {
		if strings.Index(strings.ToLower(key), "password") > -1 {
			sys.Info("  # ", key, " # ==>> ", "**********")
		} else {
			sys.Info("  # ", key, " # ==>> ", val)
		}
	})
	sys.Success("## Display All Used Configure Items [OK] #")
	sys.Success("# Startup all registered daemons [OK] #")

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	sys.Warn("# Shutting down the system ...... #")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for code, daemon := range gDaemons {
		sys.Info(" * Shutdown daemon: ", code, " ......")
		err := daemon.Shutdown(ctx)
		if err != nil {
			logger.Logger.Error("shutdown daemon failed", zap.String("code", code), zap.Error(err))
			sys.Error(" * Shutdown service ", code, " exception")
		}
		sys.Success(" * Shutdown daemon ", code, " [OK]")
	}
	sys.Success("# Shutting down the system [OK] #")
}

func init() {
	gDaemons = make(map[string]Daemon)
}
