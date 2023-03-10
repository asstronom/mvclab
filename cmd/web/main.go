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

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

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
