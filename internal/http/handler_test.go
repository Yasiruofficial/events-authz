package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"events-authz/internal/model"
	"events-authz/internal/spicedb"

	"github.com/gin-gonic/gin"
)

type fakeService struct {
	response model.CheckResponse
	err      error
}

func (f *fakeService) Check(_ context.Context, _ model.CheckRequest) (model.CheckResponse, error) {
	return f.response, f.err
}

func TestHandlerCheckReturnsValidationErrorsAsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&fakeService{
		err: fmt.Errorf("%w: subject must be in type:id format", spicedb.ErrInvalidCheckRequest),
	})
	router := gin.New()
	router.POST("/v1/check", handler.Check)

	body := []byte(`{"subject":"alice","resource":"document:budget-2026","permission":"view"}`)
	request := httptest.NewRequest(stdhttp.MethodPost, "/v1/check", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestHandlerCheckReturnsSpiceDBResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := NewHandler(&fakeService{
		response: model.CheckResponse{
			Allowed:        true,
			Permissionship: "has_permission",
			CheckedAt:      "zed-token",
		},
	})
	router := gin.New()
	router.POST("/v1/check", handler.Check)

	body := []byte(`{"subject":"user:alice","resource":"document:budget-2026","permission":"view"}`)
	request := httptest.NewRequest(stdhttp.MethodPost, "/v1/check", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var payload model.CheckResponse
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if !payload.Allowed {
		t.Fatal("expected allowed response")
	}
	if payload.Permissionship != "has_permission" {
		t.Fatalf("expected permissionship has_permission, got %q", payload.Permissionship)
	}
	if payload.CheckedAt != "zed-token" {
		t.Fatalf("expected checked_at zed-token, got %q", payload.CheckedAt)
	}
}
