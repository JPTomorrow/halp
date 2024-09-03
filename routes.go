package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/JPTomorrow/halp/config"
	"github.com/JPTomorrow/halp/db"
	"github.com/coder/websocket"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/*
This is where all the the routes are defined for the API.
*/

type WsContext struct {
	echo.Context
	chat *ChatHub
}

func initRoutes(e *echo.Echo) {
	// middleware
	ctxChat := NewChat()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &WsContext{
				c,
				ctxChat,
			}
			return next(cc)
		}
	})
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// backend routes
	e.GET("/", alive)
	e.GET("/customer-connect", customerConnect)
	e.GET("/support-connect", supportRepConnect)

	// debug only routes
	if config.DEBUG {
		e.GET("/update-schema", updateDbSchema)
	}
}

func alive(c echo.Context) error {
	return c.String(http.StatusOK, "Backend API of HALP! We are alive!")
}

func updateDbSchema(c echo.Context) error {
	schema, err := db.SchemaString(db.Customer{}, db.SupportRepresentative{}, db.SupportTicket{})
	for _, table := range schema {

		if err != nil {
			return c.String(http.StatusBadRequest, table+"\n\n"+err.Error())
		} else {
			_, err := db.Exec(table)
			if err != nil {
				return c.String(http.StatusBadRequest, table+"\n\n"+err.Error())
			}

		}
	}

	msg := strings.Join(schema, "\n\n") + "\nTables created successfully!!!\n\n"
	fmt.Println(msg)
	return c.String(http.StatusOK, msg)
}

// Connect a user to the next available support representative
func customerConnect(c echo.Context) error {
	wc := c.(*WsContext)
	conn, err := websocket.Accept(wc.Response().Writer, wc.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "Websocket connection to support closed...")

	id := wc.Request().Header.Get("id")
	if id == "" {
		conn.Close(websocket.StatusProtocolError, "no id provided.")
	}

	s, err := strconv.Atoi(id)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "invalid id.")
	}

	// print chat length
	ch, ok := wc.chat.NextAvailableRep(s, conn)
	if !ok {
		log.Printf("no available rep for id: %v", id)
		conn.Close(websocket.StatusProtocolError, "no available support rep.")
		return wc.String(http.StatusOK, "no available support rep...")
	}

	wc.chat.Run(ch, wc, conn)
	return wc.String(http.StatusOK, "Connection to support closed...")
}

func supportRepConnect(c echo.Context) error {
	wc := c.(*WsContext)
	conn, err := websocket.Accept(wc.Response().Writer, wc.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close(websocket.StatusNormalClosure, "Websocket connection to customer closed...")

	id := wc.Request().Header.Get("id")
	if id == "" {
		conn.Close(websocket.StatusProtocolError, "no id provided.")
	}

	s, err := strconv.Atoi(id)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "invalid id (during str conversion).")
	}

	ch := wc.chat.RegisterSupportRep(s, conn)
	wc.chat.Run(ch, wc, conn)
	return wc.String(http.StatusOK, "Connection to customer closed...")
}
