package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"example.com/cmd/config"
	"example.com/pkg/consoleview"
	"example.com/pkg/controller"
	"example.com/pkg/model"
	migrations "example.com/pkg/model/migrate"
	"example.com/pkg/webview"
)

var (
	isConsoleView bool
)

func main() {
	//parsing dbconfig
	dbconf, err := config.ParseDBConnConfig(`configs/dbconf.yaml`)
	if err != nil {
		log.Fatalln("error parsing config: ", err)
	}

	//applying migrations to db
	migrateurl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	pathToMigrations := "pkg/model/migrate/sql" //"../../pkg/db/migrations/sql"
	err = migrations.MigrateUp(pathToMigrations, migrateurl)
	if err != nil {
		log.Fatalf("error migrating db: %s\n", err)
	}

	//opening UserDB
	dburl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		dbconf.User, dbconf.Password, dbconf.Host, dbconf.Port, dbconf.Name)
	fmt.Println(dburl)
	userRepo, err := model.Open(context.Background(), dburl)
	if err != nil {
		log.Fatalf("error opening db: %s", err)
	}
	defer userRepo.Close()

	//creating UserService
	uservice, err := controller.NewController(userRepo)
	if err != nil {
		log.Fatalf("error creating service: %s", err)
	}

	flag.BoolVar(&isConsoleView, "console", false, "run as a console app")
	flag.Parse()

	if !isConsoleView {
		//creating server
		server, err := webview.NewServer(uservice)
		if err != nil {
			log.Fatalf("error creating server: %s", err)
		}

		//starting server
		log.Fatal(server.ListenAndServe())
	} else {
		app, err := consoleview.NewCMDApp(uservice)
		if err != nil {
			log.Fatalf("error creating cmd app: %s", err)
		}
		app.Run()
	}
}
