package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/log"
	kommendpoint "microsrv/kommersant/endpoint"
	"microsrv/kommersant/model"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

// NewHTTPHandler returns an HTTP handler that makes a set of endpoints
// available on predefined paths.
func NewHTTPHandler(endpoints kommendpoint.Endpoints, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerErrorLogger(logger),
	}

	// GET /health         retrieves service heath information

	// m := http.NewServeMux()
	// m.Handle("/sum", httptransport.NewServer(
	m := mux.NewRouter()
	m.Methods("GET").Path("/health").Handler(httptransport.NewServer(
		endpoints.HealthEndpoint,
		DecodeHTTPHealthRequest,
		EncodeHTTPGenericResponse,
		options...,
	))
	m.Methods("POST").Path("/create").Handler(httptransport.NewServer(
		endpoints.CreateEndpoint,
		DecodeHTTPKommersantRequest,
		EncodeHTTPGenericResponse,
		options...,
	))
	m.Methods("GET").Path("/result").Handler(httptransport.NewServer(
		endpoints.ResultEndpoint,
		DecodeHTTPKommersantRequest,
		EncodeHTTPGenericResponse,
		options...,
	))
	return m
}

// DecodeHTTPHealthRequest method.
func DecodeHTTPHealthRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	return model.HealthRequest{}, nil
}

// DecodeHTTPKommersantRequest method.
func DecodeHTTPKommersantRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req model.CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func err2code(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer
func EncodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(kommendpoint.Failer); ok && f.Failed() != nil {
		encodeError(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
