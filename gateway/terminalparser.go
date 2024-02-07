package gateway

import (
	"encoding/json"
	"net/url"

	"github.com/labstack/echo/v4"
)

type TerminalParserConfig struct {
	TerminalFieldName string
	DefaultTarget     *url.URL
}

func (g *gateway) TerminalParser(config TerminalParserConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if checkBypassAll(c) {
				return next(c)
			}
			var parsedBackendURL *url.URL
			req := c.Request()
			reader := copyRequestBody(req)
			defer req.Body.Close()
			defer reader.Close()

			var requestBody map[string]interface{}
			decoder := json.NewDecoder(reader)
			if err = decoder.Decode(&requestBody); err != nil {
				g.logger.Println("Error: decoding request body")
				return
			}
			terminal, ok := requestBody[config.TerminalFieldName].(string)
			if !ok {
				// 'terminal' not found bypass this middleware
				return next(c)
			}
			backendURL, err := g.repo.GetTerminal(terminal)
			if err != nil {
				return
			}
			if backendURL == "" {
				g.logger.Println("'terminal' not found. revert to default")
				parsedBackendURL = config.DefaultTarget

			} else {
				parsedBackendURL, err = url.Parse(backendURL)
				if err != nil {
					return
				}
			}
			changeProxyDirector(g.proxy, parsedBackendURL)
			setBypassAll(c)
			return next(c)
		}
	}
}
