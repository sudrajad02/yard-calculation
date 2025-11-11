package main

import (
	"log"
	"yard-calculation/config"
	"yard-calculation/handlers"
	"yard-calculation/models"
	"yard-calculation/repositories"
	"yard-calculation/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Migrate the schema
	config.DB.AutoMigrate(&models.Yard{}, &models.Block{}, &models.Container{}, &models.YardPlan{})

	// Initialize Repository
	containerRepo := repositories.NewContainerRepository(config.DB)

	// Initialize Service
	containerService := services.NewContainerService(containerRepo)

	// Initialize Handler
	containerHandler := handlers.NewContainerHandler(containerService)

	// Initialize Fiber App
	app := fiber.New()
	app.Use(logger.New())

	// Middleware
	app.Use(cors.New())

	// Define Routes
	app.Post("/suggestion", containerHandler.GetSuggestion)
	app.Post("/placement", containerHandler.PlaceContainer)
	app.Post("/pickup", containerHandler.PickupContainer)

	// GORM tidak otomatis membuat indeks unik untuk foreign key.
	// Kita tambahkan manual jika diperlukan untuk performa.
	// config.DB.Migrator().CreateIndex(&models.Block{}, "YardID") // Contoh
	// config.DB.Migrator().CreateIndex(&models.Container{}, "BlockID")
	// config.DB.Migrator().CreateIndex(&models.Container{}, "IsPlaced") // Berguna untuk query pickup
	// config.DB.Migrator().CreateIndex(&models.Container{}, "ContainerNumber") // Sudah ada dari tag gorm:"uniqueIndex"

	log.Fatal(app.Listen(":3003"))
}
