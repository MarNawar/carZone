package service

import (
	"context"

	"github.com/MarNawar/carZone/models"
)

type CarServiceInterface interface {
	GetCarById(context.Context, string) (*models.Car, error)
	GetCarsByBrand(context.Context, string, bool)([]models.Car, error)
	CreateCar(context.Context, *models.CarRequest)(*models.Car, error)
	UpdateCar(context.Context, string, *models.CarRequest)(*models.Car, error)
	DeleteCar(context.Context, string)(*models.Car, error)
}

type EngineServiceInterface interface{
	GetEngineByID(context.Context, string)(*models.Engine, error)
	CreateEngine(context.Context, *models.EngineRequest)(*models.Engine, error)
	UpdateEngine(context.Context, *models.EngineRequest, string)(*models.Engine, error)
	DeleteEngine(context.Context, string)(*models.Engine, error)
}