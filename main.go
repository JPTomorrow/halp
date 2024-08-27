/*
The halp support chat API is a simple yet robust implementation of a chat support system in Golang!
*/
package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/JPTomorrow/halp/config"
	"github.com/JPTomorrow/halp/db"
	dir_tree "github.com/JPTomorrow/halp/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	e.HideBanner = true
	fmt.Printf("---------------------\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("- DataMux File Ingest\n")
	fmt.Printf("----------------------------\n")
	fmt.Printf("- System: Backend API\n")
	fmt.Printf("- Author: Justin Morrow\n")
	fmt.Printf("- Year: 2024\n")
	fmt.Printf("-------------------------\n")
	fmt.Printf("-----------------------\n\n")

	if config.DEBUG {
		fmt.Println("App is in DEBUG mode")
	} else {
		fmt.Println("App is in PRODUCTION mode")
	}

	cdw, _ := os.Getwd()
	fmt.Println("EXEC PATH: " + cdw + "/" + os.Args[0])
	dir_tree.Print()

	err := godotenv.Load()
	if err != nil {
		e.Logger.Fatal("Error loading .env file")
	} else {
		e.Logger.Info("Loaded .env file successfully")
	}

	db.InitDb()
	defer db.CloseDB()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.ERROR)

	initRoutes(e)

	go func() {
		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
	}()

	// any async stuff here

	// wait
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
