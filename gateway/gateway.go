// gateway is a reverse proxy that routes requests to different backends based on defined
// middlewares. you can use bellow template for future parser middlewares
//
// type DefaultParserConfig struct{}
//
//	func (g *gateway) DefaultParser(config DefaultParserConfig) echo.MiddlewareFunc {
//		return func(next echo.HandlerFunc) echo.HandlerFunc {
//			return func(c echo.Context) (err error) {
//				if checkBypassAll(c) {
//					return next(c)
//				}
//				// logic
//				return next(c)
//			}
//		}
//	}
package gateway

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v4"
)

const bypassall = "bypassall"

type gateway struct {
	proxy  *httputil.ReverseProxy
	logger *log.Logger
	repo   Repository
}

func New(proxy *httputil.ReverseProxy, logger *log.Logger, repo Repository) *gateway {
	return &gateway{
		proxy:  proxy,
		logger: logger,
		repo:   repo,
	}
}

type Repository interface {
	GetTerminal(terminal string) (addr string, err error)
}

func changeProxyDirector(p *httputil.ReverseProxy, u *url.URL) (err error) {
	if u.Scheme == "" || u.Host == "" {
		return errors.New("URL Scheme or Host cannot be empty")
	}
	p.Director = func(req *http.Request) {
		req.URL.Scheme = u.Scheme
		req.URL.Host = u.Host
		req.Host = u.Host
	}
	return
}

// copyRequestBody duplicates request's Body and preserves the initial Body of request
// because reading body will drain it
func copyRequestBody(r *http.Request) io.ReadCloser {
	buf, _ := io.ReadAll(r.Body)
	reader1 := io.NopCloser(bytes.NewBuffer(buf))
	reader2 := io.NopCloser(bytes.NewBuffer(buf))
	r.Body = reader2
	return reader1
}

// setBypassAll sets a flag in context
// all parsers must first check this value toi see if further process is needed
func setBypassAll(c echo.Context) {
	c.Set(bypassall, bypassall)
}

func checkBypassAll(c echo.Context) bool {
	bypass, ok := c.Get(bypassall).(string)
	if !ok || bypass != bypassall {
		return false
	}
	return true
}

// Handler
func (g *gateway) ProxyHandler(c echo.Context) error {
	g.proxy.ServeHTTP(c.Response(), c.Request())
	return nil
}
