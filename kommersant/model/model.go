package model

// HealthRequest collects the request parameters for the Health method.
type HealthRequest struct{}

// HealthResponse collects the response values for the Health method.
type HealthResponse struct {
	Healthy bool  `json:"healthy,omitempty"`
	Err     error `json:"err,omitempty"`
}

// Failed implements Failer.
func (r HealthResponse) Failed() error { return r.Err }

// CreateRequest collects the request parameters for the Greeting method.
type CreateRequest struct {
	AdNum string `json:"ad_num,omitempty"`
}

// CreateResponse collects the response values for the Greeting method.
type CreateResponse struct {
	Status  int32  `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Err     error  `json:"err,omitempty"`
}

// Failed implements Failer.
func (r CreateResponse) Failed() error { return r.Err }
