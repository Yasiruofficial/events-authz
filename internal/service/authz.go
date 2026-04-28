package service

import (
	"context"
	"time"

	"events-authz/internal/cache"
	"events-authz/internal/model"
	"events-authz/internal/spicedb"
)

type AuthzService struct {
	spice *spicedb.Client
	cache *cache.Cache
}

func NewAuthzService(s *spicedb.Client, c *cache.Cache) *AuthzService {
	return &AuthzService{spice: s, cache: c}
}

func (s *AuthzService) Check(ctx context.Context, req model.CheckRequest) (bool, error) {
	key := req.Subject + "|" + req.Resource + "|" + req.Permission

	if val, ok := s.cache.Get(key); ok {
		return val.(bool), nil
	}

	allowed, err := s.spice.CheckPermission(ctx, req.Subject, req.Resource, req.Permission)
	if err != nil {
		return false, err
	}

	s.cache.Set(key, allowed, 5*time.Second)
	return allowed, nil
}
