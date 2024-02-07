package main

import (
	"flag"
	"fmt"
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/omid9h/ipggateway/gateway"
	"github.com/omid9h/ipggateway/gateway/repo"
	"github.com/omid9h/ipggateway/pkg/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/natefinch/lumberjack"
)

var (
	conf = flag.String("config", "config.json", "path to config file")
)

func main() {

	flag.Parse()

	// config
	cfg, err := gateway.NewConfig(*conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// log
	if cfg.LogDir == "" {
		logPath, err := util.GetDirectory("logs")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfg.LogDir = logPath
	}
	gatewayLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, "app.log"),
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}
	gatewayLogger := log.New(gatewayLogFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	accessLogFile := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.LogDir, "access.log"),
		MaxSize:    50,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	// reverse proxy
	defaultURL, err := url.Parse(cfg.DefaultTarget)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	proxy := httputil.NewSingleHostReverseProxy(defaultURL)

	repo := repo.New(cfg.DBPath)

	// gateway
	ipggateway := gateway.New(proxy, gatewayLogger, repo)

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: accessLogFile,
	}))
	e.Use(middleware.Recover())

	e.Pre(ipggateway.TerminalParser(gateway.TerminalParserConfig{
		TerminalFieldName: cfg.TerminalFieldName,
		DefaultTarget:     defaultURL,
	}))

	e.Pre(ipggateway.TokenParser(gateway.TokenParserConfig{
		DefaultTarget: defaultURL,
	}))

	// Routes
	e.Any("/*", ipggateway.ProxyHandler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
