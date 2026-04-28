package model

import "encoding/json"

type CheckRequest struct {
	Subject     string         `json:"subject" binding:"required"`
	Resource    string         `json:"resource" binding:"required"`
	Permission  string         `json:"permission" binding:"required"`
	Context     map[string]any `json:"context,omitempty"`
	Consistency string         `json:"consistency,omitempty"`
	ZedToken    string         `json:"zed_token,omitempty"`
}

func (r CheckRequest) CacheKey() (string, error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return "", err
	}

	return string(payload), nil
}

type CheckResponse struct {
	Allowed        bool   `json:"allowed"`
	Permissionship string `json:"permissionship,omitempty"`
	CheckedAt      string `json:"checked_at,omitempty"`
}
