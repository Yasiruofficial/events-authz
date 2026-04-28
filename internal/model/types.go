package model

type CheckRequest struct {
	Subject    string `json:"subject"`
	Resource   string `json:"resource"`
	Permission string `json:"permission"`
}

type CheckResponse struct {
	Allowed bool `json:"allowed"`
}
