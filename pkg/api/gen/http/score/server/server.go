// Code generated by goa v3.4.2, DO NOT EDIT.
//
// Score HTTP server
//
// Command:
// $ goa gen sunflower/pkg/api/design -o pkg/api/

package server

import (
	"context"
	"net/http"
	score "sunflower/pkg/api/gen/score"

	goahttp "goa.design/goa/v3/http"
	goa "goa.design/goa/v3/pkg"
)

// Server lists the Score service endpoint HTTP handlers.
type Server struct {
	Mounts      []*MountPoint
	Upload      http.Handler
	ScoreList   http.Handler
	ScoreDetail http.Handler
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the design.
type ErrorNamer interface {
	ErrorName() string
}

// MountPoint holds information about the mounted endpoints.
type MountPoint struct {
	// Method is the name of the service method served by the mounted HTTP handler.
	Method string
	// Verb is the HTTP method used to match requests to the mounted handler.
	Verb string
	// Pattern is the HTTP request path pattern used to match requests to the
	// mounted handler.
	Pattern string
}

// New instantiates HTTP handlers for all the Score service endpoints using the
// provided encoder and decoder. The handlers are mounted on the given mux
// using the HTTP verb and path defined in the design. errhandler is called
// whenever a response fails to be encoded. formatter is used to format errors
// returned by the service methods prior to encoding. Both errhandler and
// formatter are optional and can be nil.
func New(
	e *score.Endpoints,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) *Server {
	return &Server{
		Mounts: []*MountPoint{
			{"Upload", "POST", "/api/score/upload"},
			{"ScoreList", "GET", "/api/score"},
			{"ScoreDetail", "GET", "/api/score/{id}"},
		},
		Upload:      NewUploadHandler(e.Upload, mux, decoder, encoder, errhandler, formatter),
		ScoreList:   NewScoreListHandler(e.ScoreList, mux, decoder, encoder, errhandler, formatter),
		ScoreDetail: NewScoreDetailHandler(e.ScoreDetail, mux, decoder, encoder, errhandler, formatter),
	}
}

// Service returns the name of the service served.
func (s *Server) Service() string { return "Score" }

// Use wraps the server handlers with the given middleware.
func (s *Server) Use(m func(http.Handler) http.Handler) {
	s.Upload = m(s.Upload)
	s.ScoreList = m(s.ScoreList)
	s.ScoreDetail = m(s.ScoreDetail)
}

// Mount configures the mux to serve the Score endpoints.
func Mount(mux goahttp.Muxer, h *Server) {
	MountUploadHandler(mux, h.Upload)
	MountScoreListHandler(mux, h.ScoreList)
	MountScoreDetailHandler(mux, h.ScoreDetail)
}

// MountUploadHandler configures the mux to serve the "Score" service "Upload"
// endpoint.
func MountUploadHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("POST", "/api/score/upload", f)
}

// NewUploadHandler creates a HTTP handler which loads the HTTP request and
// calls the "Score" service "Upload" endpoint.
func NewUploadHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeUploadRequest(mux, decoder)
		encodeResponse = EncodeUploadResponse(encoder)
		encodeError    = EncodeUploadError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "Upload")
		ctx = context.WithValue(ctx, goa.ServiceKey, "Score")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		data := &score.UploadRequestData{Payload: payload.(*score.UploadPayload), Body: r.Body}
		res, err := endpoint(ctx, data)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountScoreListHandler configures the mux to serve the "Score" service
// "ScoreList" endpoint.
func MountScoreListHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/score", f)
}

// NewScoreListHandler creates a HTTP handler which loads the HTTP request and
// calls the "Score" service "ScoreList" endpoint.
func NewScoreListHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeScoreListRequest(mux, decoder)
		encodeResponse = EncodeScoreListResponse(encoder)
		encodeError    = EncodeScoreListError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "ScoreList")
		ctx = context.WithValue(ctx, goa.ServiceKey, "Score")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}

// MountScoreDetailHandler configures the mux to serve the "Score" service
// "ScoreDetail" endpoint.
func MountScoreDetailHandler(mux goahttp.Muxer, h http.Handler) {
	f, ok := h.(http.HandlerFunc)
	if !ok {
		f = func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		}
	}
	mux.Handle("GET", "/api/score/{id}", f)
}

// NewScoreDetailHandler creates a HTTP handler which loads the HTTP request
// and calls the "Score" service "ScoreDetail" endpoint.
func NewScoreDetailHandler(
	endpoint goa.Endpoint,
	mux goahttp.Muxer,
	decoder func(*http.Request) goahttp.Decoder,
	encoder func(context.Context, http.ResponseWriter) goahttp.Encoder,
	errhandler func(context.Context, http.ResponseWriter, error),
	formatter func(err error) goahttp.Statuser,
) http.Handler {
	var (
		decodeRequest  = DecodeScoreDetailRequest(mux, decoder)
		encodeResponse = EncodeScoreDetailResponse(encoder)
		encodeError    = EncodeScoreDetailError(encoder, formatter)
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), goahttp.AcceptTypeKey, r.Header.Get("Accept"))
		ctx = context.WithValue(ctx, goa.MethodKey, "ScoreDetail")
		ctx = context.WithValue(ctx, goa.ServiceKey, "Score")
		payload, err := decodeRequest(r)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		res, err := endpoint(ctx, payload)
		if err != nil {
			if err := encodeError(ctx, w, err); err != nil {
				errhandler(ctx, w, err)
			}
			return
		}
		if err := encodeResponse(ctx, w, res); err != nil {
			errhandler(ctx, w, err)
		}
	})
}
