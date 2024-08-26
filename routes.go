package main

import (
	"fmt"
	"net/http"

	"github.com/JPTomorrow/halp/config"
	"github.com/JPTomorrow/halp/db"
	"github.com/gorilla/websocket"

	"github.com/JPTomorrow/halp/auth"

	"github.com/labstack/echo/v4"
)

/*
This is where all the the routes are defined for the API.
*/

func initRoutes(e *echo.Echo) {
	// backend routes
	e.GET("/", alive)
	e.GET("/connect-to-support", connectWithNextSupportRep)
	e.GET("/queue-for-support", supportRepQueue)
	// e.GET("/add-file-profile", addIngestFileProfile)
	// e.GET("/login", tokenLogin)
	// e.POST("/create-email-account", createEmailPasswordAccount)
	// e.GET("/resource", resource)

	// debug only routes
	if config.DEBUG {
		e.GET("/update-schema", updateDbSchema)
	}
}

func alive(c echo.Context) error {
	return c.String(http.StatusOK, "Backend API of HALP! We are alive!")
}

const (
	customerPoolSize = 10000
	salesRepPoolSize = 10000
)

var (
	upgrader     = websocket.Upgrader{}
	customerPool = make([]db.User, customerPoolSize)
	salesRepPool = make([]db.User, customerPoolSize)
)

// Connect a user to the next available support representative
func connectWithNextSupportRep(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	isCustomer := c.Request().Header.Get("is-customer")
	if isCustomer == "" {
		return c.String(http.StatusBadRequest, "Missing is-customer header")
	} else if isCustomer == "true" {
		fmt.Println("YOU ARE A CUSTOMER")
	} else if isCustomer == "false" {
		fmt.Println("YOU ARE A SUPPORT REP")
	}

	exit := false

	for exit {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}

	return c.String(http.StatusOK, "Websocket closed. Support chat finished!")
}

func supportRepQueue(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	isCustomer := c.Request().Header.Get("is-customer")
	if isCustomer == "" {
		return c.String(http.StatusBadRequest, "Missing is-customer header")
	} else if isCustomer == "true" {
		fmt.Println("YOU ARE A CUSTOMER")
	} else if isCustomer == "false" {
		fmt.Println("YOU ARE A SUPPORT REP")
	}

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
	return c.String(http.StatusOK, "Websocket closed. Support chat finished!")
}

// func addIngestFileProfile(c echo.Context) error {
// 	profile := db.IngestFileProfile{
// 		Name:             "",
// 		Description:      "",
// 		FirstCreated:     time.Now(),
// 		LastUpdated:      time.Now(),
// 		FilePath:         "",
// 		AiOcrTextSummary: "",
// 	}

// 	json.NewDecoder(c.Request().Body).Decode(&profile)
// 	_, insert_err := db.Insert(profile)
// 	if insert_err != nil {
// 		return c.String(http.StatusOK, "ADD PROFILE ERROR: database insert\n\n"+insert_err.Error())
// 	}

// 	return c.String(http.StatusOK, "File orofile '"+profile.Name+"' was added successfully!!!")
// }

func updateDbSchema(c echo.Context) error {
	test_str, err := db.PushSchema(db.SupportTicket{}, db.User{})
	if err != nil {
		return c.String(http.StatusOK, test_str+"\n\n"+err.Error())
	} else {
		msg := test_str + "\n\nTable created successfully!!!\n\n"
		fmt.Println(msg)
		return c.String(http.StatusOK, msg)
	}
}

// func createEmailPasswordAccount(c echo.Context) error {
// 	user := db.User{
// 		Username: "",
// 		Password: "",
// 		Email:    "",
// 		Phone:    "",
// 	}
// 	json.NewDecoder(c.Request().Body).Decode(&user)

// 	if user.Username == "" || user.Password == "" || user.Email == "" || user.Phone == "" {
// 		return c.String(http.StatusBadRequest, "Missing required fields")
// 	}

// 	_, qerr := db.Query(user, "username")
// 	if qerr != nil {
// 		return c.String(http.StatusInternalServerError, "Email account creation error: username "+user.Username+" already exists\n\n"+qerr.Error())
// 	}

// 	_, err := db.Insert(user)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, "Email account creation error: database insert\n\n"+err.Error())
// 	}

// 	return c.String(http.StatusOK, "Account Created: "+user.Username)
// }

// func tokenLogin(c echo.Context) error {
// 	// @TODO: get username and password from request and then look to see if it is in the database already
// 	un := c.Request().FormValue("username")
// 	pw := c.Request().FormValue("password")

// 	if un == "" || pw == "" {
// 		return c.String(http.StatusBadRequest, "Missing username or password")
// 	}

// 	rows, qerr := db.Query(db.IngestFileProfile{})
// 	if qerr == nil {
// 		// found the user
// 		for rows.Next() {
// 			profile := db.IngestFileProfile{
// 				Name: un,
// 			}
// 			serr := rows.Scan(&profile.Id, &profile.Name, &profile.FirstCreated, &profile.LastUpdated)
// 			if serr != nil {
// 				fmt.Printf("SCAN ERROR: %v\n\n", serr)
// 			}
// 			fmt.Printf("PROFILE: %v\n\n", profile)
// 		}
// 	} else {
// 		// did not find the user
// 		fmt.Printf("USER NOT FOUND: %v\n\n", qerr)
// 		fmt.Printf("CREATING USER: %v\n\n", un)
// 		// db.Insert(db.User{})
// 	}

// 	bearerToken := c.Request().Header.Get("Authorization")
// 	if !auth.IsValidLoginToken(bearerToken) {
// 		return c.String(http.StatusUnauthorized, "Unauthorized")
// 	}

// 	return c.String(http.StatusOK, "Authorized")
// }

func resource(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")
	if !auth.IsValidLoginToken(bearerToken) {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	return c.String(http.StatusOK, "Authorized")
}
