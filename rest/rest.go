package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hootuu/gelato/configure"
	"github.com/hootuu/gelato/errors"
	"github.com/hootuu/gelato/sys"
	"net/http"
	"strings"
	"time"
)

type Rest struct {
	code       string
	httpServer *http.Server
	ginEngine  *gin.Engine
}

func New(code string) *Rest {
	ginEngine := gin.New()
	regCode := strings.ToUpper(code)
	cfgCode := strings.ToLower(code)
	addr := configure.GetString(fmt.Sprintf("rest.%s.addr", cfgCode), ":8080")
	readTimeout := configure.GetDuration(fmt.Sprintf("rest.%s.timeout.read", cfgCode), 30)
	readHeaderTimeout := configure.GetDuration(fmt.Sprintf("rest.%s.timeout.header.read", cfgCode), 30)
	writeTimeout := configure.GetDuration(fmt.Sprintf("rest.%s.timeout.write", cfgCode), 30)
	idleTimeout := configure.GetDuration(fmt.Sprintf("rest.%s.timeout.idle", cfgCode), 30)
	return &Rest{
		code: regCode,
		httpServer: &http.Server{
			Handler:           ginEngine,
			Addr:              addr,
			ReadTimeout:       readTimeout * time.Second,
			ReadHeaderTimeout: readHeaderTimeout * time.Second,
			WriteTimeout:      writeTimeout * time.Second,
			IdleTimeout:       idleTimeout * time.Second,
		},
		ginEngine: ginEngine,
	}
}

func (r *Rest) Code() string {
	return r.code
}

func (r *Rest) Startup() *errors.Error {
	sys.Info(" * Rest [", r.code, "] startup on: ", r.httpServer.Addr)
	go func() {
		err := r.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sys.Error(" * rest [", r.code, "] listen failed")
			sys.Exit(errors.System(err.Error()))
			return
		}
	}()

	return nil
}

func (r *Rest) Shutdown(ctx context.Context) *errors.Error {
	if err := r.httpServer.Shutdown(ctx); err != nil {
		sys.Error(" * Rest [", r.code, "] shutdown error: ", err.Error())
		return errors.System("rest shutdown error", err)
	}
	sys.Error(" * Rest [", r.code, "] exist: ", r.httpServer.Addr)
	return nil
}

func (r *Rest) Register(call func(router *gin.RouterGroup)) {
	call(&r.ginEngine.RouterGroup)
}
