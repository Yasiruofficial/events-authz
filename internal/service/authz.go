package service

import (
	"context"
	"time"

	"github.com/spicedb/spicedb-go/internal/cache"
	"github.com/spicedb/spicedb-go/internal/model"
)

type permissionChecker interface {
	CheckPermission(ctx context.Context, req model.CheckRequest) (model.CheckResponse, error)
}

type AuthzService struct {
	spice permissionChecker
	cache *cache.Cache
}

func NewAuthzService(s permissionChecker, c *cache.Cache) *AuthzService {
	return &AuthzService{spice: s, cache: c}
}

func (s *AuthzService) Check(ctx context.Context, req model.CheckRequest) (model.CheckResponse, error) {
	key, err := req.CacheKey()
	if err == nil {
		if val, ok := s.cache.Get(key); ok {
			if cached, ok := val.(model.CheckResponse); ok {
				return cached, nil
			}
		}
	}

	allowed, err := s.spice.CheckPermission(ctx, req)
	if err != nil {
		return model.CheckResponse{}, err
	}

	if key != "" {
		s.cache.Set(key, allowed, 5*time.Second)
	}

	return allowed, nil
}
