package service

import (
	"Timos-API/Vuement/persistence"
	"context"
)

type ComponentService struct {
	p *persistence.ComponentPersistor
}

func NewComponentService(p *persistence.ComponentPersistor) *ComponentService {
	return &ComponentService{p}
}

func (s *ComponentService) Create(ctx context.Context, component persistence.Component) (*persistence.Component, error) {
	// TODO: Validate, Clean
	return s.p.Create(ctx, component)
}

func (s *ComponentService) Update(ctx context.Context, id string, component persistence.Component) (*persistence.Component, error) {
	// TODO: Validate, Clean
	return s.p.Update(ctx, id, component)
}

func (s *ComponentService) Delete(ctx context.Context, id string) (bool, error) {
	return s.p.Delete(ctx, id)
}

func (s *ComponentService) GetById(ctx context.Context, id string) (*persistence.Component, error) {
	return s.p.GetById(ctx, id)
}

func (s *ComponentService) GetAll(ctx context.Context) (*[]persistence.Component, error) {
	return s.p.GetAll(ctx)
}
