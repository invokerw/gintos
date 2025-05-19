package main

import (
	"flag"
	"github/invokerw/gintos/config"
	"github/invokerw/gintos/config/file"
	"github/invokerw/gintos/demo/internal/conf"
	"github/invokerw/gintos/demo/internal/g"
	"github/invokerw/gintos/demo/internal/initialize"
	"github/invokerw/gintos/log"
	"github/invokerw/gintos/log/zap"
	"os"

	"github.com/gin-gonic/gin"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagConf is the config flag.
	flagConf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagConf, "conf", "./configs", "config path, eg: -conf config.yaml")
}

type App struct {
	engine *gin.Engine
}

func (a *App) Run(addr ...string) error {
	return a.engine.Run(addr...)
}

func newApp(_ *initialize.InitRet, engine *gin.Engine) *App {
	return &App{engine: engine}
}

func main() {
	flag.Parse()
	logger := log.With(zap.NewLogger(zap.NewZapLogger()))
	_ = logger
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	g.Config = &bc
	g.Log = log.NewHelper(log.With(logger, "module", "global"))

	app, cleanup, err := wireApp(bc.Server, bc.Data, bc.File, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(bc.Server.Http.Addr); err != nil {
		panic(err)
	}
}
