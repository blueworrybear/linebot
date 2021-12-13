package main

import (
	"fmt"
	"log"
	"os"

	"github.com/blueworrybear/lineBot/config"
	"github.com/blueworrybear/lineBot/core"
	"github.com/blueworrybear/lineBot/routers"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectDatabase(cfg *config.Config) *gorm.DB {
	var x *gorm.DB
	var err error
	switch cfg.Database.Driver {
	case "sqlite3":
		x, err = gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	case "cloudrun":
		x, err = gorm.Open(
			postgres.Open(
				fmt.Sprintf(
					"user=%s password=%s database=%s host=%s/%s",
					cfg.Database.User,
					cfg.Database.Password,
					cfg.Database.Name,
					cfg.Database.Socket,
					cfg.Database.Instance,
				)), &gorm.Config{})
	default:
		log.Fatal("database driver not support")
	}
	if err != nil {
		log.Fatal(err)
	}
	return x
}

func Run(c *cli.Context) error{
	cfg, err := config.Environ()
	if err != nil {
		return err
	}
	app, err := InitializeApplication(cfg, connectDatabase(cfg))
	if err != nil {
		return err
	}

	if cfg.Database.AutoMigrate {
		go func ()  {
			if err := app.db.Migrate(); err != nil {
				log.Panic(err)
			}
			log.Panicln("migration done")
		}()
	}

	r := gin.Default()
	app.routers.RegisterRoutes(r)
	port := os.Getenv("PORT")
	if port == "" {
			port = "8080"
			log.Printf("defaulting to port %s", port)
	}
	r.Run(fmt.Sprintf(":%s", port))
	return nil
}

func main()  {
	app := &cli.App{
		Name: "lineBot",
		Version: "0.1",
		Action: Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type application struct {
	routers *routers.Routers
	db core.DatabaseService
}

func newApplication(routers *routers.Routers, db core.DatabaseService) application {
	return application{
		routers: routers,
		db: db,
	}
}
