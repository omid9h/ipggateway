package gateway

import (
	"fmt"
	"net/url"

	"github.com/labstack/echo/v4"
)

type TokenParserConfig struct {
	DefaultTarget *url.URL
}

func (g *gateway) TokenParser(config TokenParserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			fmt.Println("TokenParser Begin")
			if checkBypassAll(c) {
				fmt.Println("checkBypassAll")
				return next(c)
			}
			fmt.Println("TokenParser Logic")
			return next(c)
		}
	}
}
