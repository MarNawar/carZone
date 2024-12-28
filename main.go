package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/MarNawar/carZone/driver"
	carHandler "github.com/MarNawar/carZone/handler/car"
	engineHandler "github.com/MarNawar/carZone/handler/engine"
	loginHandler "github.com/MarNawar/carZone/handler/login"
	"github.com/MarNawar/carZone/middleware"
	carService "github.com/MarNawar/carZone/service/car"
	engineService "github.com/MarNawar/carZone/service/engine"
	carStore "github.com/MarNawar/carZone/store/car"
	engineStore "github.com/MarNawar/carZone/store/engine"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	driver.InitDB()
	defer driver.CloseDB()

	db := driver.GetDB()
	carStore := carStore.New(db)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)

	router := gin.New()
	router.Use(gin.Logger())


	schemaFile := "./store/schema.sql"
	if err := executeSchemaFile(db, schemaFile); err != nil{
		log.Fatal("error while executing the schema file")
	}

	//login
	router.POST("/login", loginHandler.Login)
	router.Use(middleware.AuthMiddleware())

	// car router
	router.GET("/car/:id", carHandler.HandleGetCarByID)
	router.GET("/cars", carHandler.HandleGetCarByBrand)
	router.POST("/car", carHandler.HandleCreateCar)
	router.PUT("/car/:id", carHandler.HandleUpdateCar)
	router.DELETE("/car/:id", carHandler.HandleDeleteCar)

	// engine router
	router.GET("/engine/:id", engineHandler.HandleGetEngineByID)
	router.POST("/engine", engineHandler.HandleCreateEngine)
	router.PUT("/engine/:id", engineHandler.HandleUpdateEngine)
	router.DELETE("/engine/:id", engineHandler.HandleDeleteEngine)

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func executeSchemaFile(db *sql.DB, fileName string)error{
	sqlFile, err := os.ReadFile(fileName)
	if err != nil{
		log.Printf("Error reading file: %v", err)
		return err;
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil{
		log.Printf("Error executing SQL: %v", err)
		return err
	}
	return nil
}