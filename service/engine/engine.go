package engine

import (
	"context"

	"github.com/MarNawar/carZone/models"
	"github.com/MarNawar/carZone/store"
)
 
type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) *EngineService {
	return &EngineService{
		store: store,
	}
}


func (s *EngineService) GetEngineByID(ctx context.Context, id string)(*models.Engine, error){
	engine, err := s.store.EngineById(ctx, id)

	if err != nil{
		return nil, err
	}
	return &engine, nil
}

func (s *EngineService)CreateEngine(ctx context.Context, engineReq *models.EngineRequest)(*models.Engine, error){
	err := models.ValidateEngineRequest(*engineReq)
	if err != nil{
		return nil, err
	}
	
	engine, err := s.store.CreateEngine(ctx, engineReq)
	
	if err != nil{
		return nil, err
	}
	return &engine, nil
}

func (s *EngineService)UpdateEngine(ctx context.Context, engineReq *models.EngineRequest, id string)(*models.Engine, error){
	err := models.ValidateEngineRequest(*engineReq)
	if err != nil{
		return nil, err
	}
	
	engine, err := s.store.EngineUpdate(ctx, id, engineReq)
	
	if err != nil{
		return nil, err
	}
	return &engine, nil
}

func (s *EngineService)DeleteEngine(ctx context.Context, id string)(*models.Engine, error){	
	engine, err := s.store.EngineDelete(ctx, id)
	
	if err != nil{
		return nil, err
	}
	return &engine, nil
}
