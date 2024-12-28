package car

import (
	"context"

	"github.com/MarNawar/carZone/models"
	"github.com/MarNawar/carZone/store"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) *CarService {
	return &CarService{
		store: store,
	}
}

func (s *CarService) GetCarById(ctx context.Context, id string)(*models.Car, error){
	car, err := s.store.GetCarById(ctx, id)
	if err != nil{
		return nil, err
	}
	return &car, nil
}

func (s *CarService)GetCarsByBrand(ctx context.Context, brand string, isEngine bool)([]models.Car, error){
	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)

	if err != nil{
		return nil, err
	}
	return cars, nil
}


func (s *CarService)CreateCar(ctx context.Context, car *models.CarRequest)(*models.Car, error){
	if err := models.ValidateRequest(*car); err != nil{
		return nil, err
	}
	createdCar, err := s.store.CreateCar(ctx, car)
	if err != nil{
		return nil, err
	}
	return &createdCar, nil
}


func (s *CarService)UpdateCar(ctx context.Context, id string, car *models.CarRequest)(*models.Car, error){
	if err := models.ValidateRequest(*car); err != nil{
		return nil, err
	}
	updatedCar, err := s.store.UpdateCar(ctx, id, car)
	if err != nil{
		return nil, err
	}
	return &updatedCar, nil
}

func ( s *CarService) DeleteCar(ctx context.Context, id string)(*models.Car, error){
	deletedCar, err := s.store.DeleteCar(ctx, id)
	if err != nil{
		return nil, err
	}
	return &deletedCar, nil
}