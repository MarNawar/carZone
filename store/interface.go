package store

import (
	"context"

	"github.com/MarNawar/carZone/models"
)

type CarStoreInterface interface {
	GetCarById(context.Context, string) (models.Car, error)
	GetCarByBrand(context.Context, string, bool) ([]models.Car, error)
	CreateCar(context.Context, *models.CarRequest) (models.Car, error)
	UpdateCar(context.Context, string, *models.CarRequest) (models.Car, error)
	DeleteCar(context.Context, string) (models.Car, error)
}

type EngineStoreInterface interface{
	EngineById(context.Context, string) (models.Engine, error)
	CreateEngine(context.Context, *models.EngineRequest) (models.Engine, error)
	EngineUpdate(context.Context, string, *models.EngineRequest) (models.Engine, error) 
	EngineDelete(context.Context, string) (models.Engine, error)
}

