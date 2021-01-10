package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/api/service/configuration", handlerCheckCaddy)
	e.GET("/api/redirect/add/safe", handlerAddSafeRedirect)
	e.GET("/api/redirect/add/stable", handlerAddStableRedirect)
	e.GET("/api/redirect/add/no-sniff", handlerAddNoSniffRedirect)

	// Start server
	initCaddy()
	e.Logger.Fatal(e.Start("127.0.0.1:1323"))
}

// Handler
func handlerCheckCaddy(c echo.Context) error {
	return c.String(http.StatusOK, checkCaddy())
}

func handlerAddSafeRedirect(c echo.Context) error {
	go addSafeRedirect(c.QueryParam("from"), c.QueryParam("to"))
	return c.String(http.StatusOK, "ok")
}

func handlerAddStableRedirect(c echo.Context) error {
	go addStableRedirect(c.QueryParam("from"), c.QueryParam("to"))
	return c.String(http.StatusOK, "ok")
}

func handlerAddNoSniffRedirect(c echo.Context) error {
	go addNoSniffRedirect(c.QueryParam("from"), c.QueryParam("to"), c.QueryParam("canonical"))
	return c.String(http.StatusOK, "ok")
}
