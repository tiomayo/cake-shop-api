package main

import (
	"cake-store/internal/cakes"
	"cake-store/internal/middlewares"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"time"

	_ "cake-store/docs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var dbInstance *sql.DB

func initDB() (*sql.DB, error) {
	conf := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  os.Getenv("DB_NET"),
		Addr:                 os.Getenv("DB_ADDRESS"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
		ParseTime:            true,
	}
	if dbInstance == nil {
		newInstance, err := sql.Open("mysql", conf.FormatDSN())
		if err != nil {
			return nil, err
		}
		if err = newInstance.Ping(); err != nil {
			return nil, err
		}
		dbInstance = newInstance
		dbInstance.SetConnMaxIdleTime(3)
		dbInstance.SetMaxOpenConns(10)
		dbInstance.SetConnMaxLifetime(time.Hour)
		return dbInstance, nil
	}
	return dbInstance, nil
}

// @title Cake Store API
// @version 1.0
// @description Cake store API for testing purposes.

func main() {
	godotenv.Load(".env")
	e := echo.New()

	e.Use(middleware.Recover())

	db, err := initDB()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders: []string{echo.HeaderContentLength, echo.HeaderContentType, "Pagination-Rows", "Pagination-Page", "Pagination-Limit"},
	}))
	middlewares.UseCustomValidatorHandler(e)
	e.Use(middleware.Logger())

	// Init Repo
	cakesRepo := cakes.NewRepository(db)

	// Init Handler
	cakesHandler := cakes.NewHandler(cakesRepo)

	// Routes
	e.GET("/cakes", cakesHandler.List)
	e.GET("/cakes/:id", cakesHandler.Get)
	e.POST("/cakes", cakesHandler.Create)
	e.PATCH("/cakes/:id", cakesHandler.Update)
	e.DELETE("/cakes/:id", cakesHandler.Delete)

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]interface{}{"message": "API OK"})
	})
	e.GET("/docs/*", echoSwagger.WrapHandler)
	e.Logger.Fatal(e.Start(":8080"))
}
