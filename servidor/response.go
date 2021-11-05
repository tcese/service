package servidor

import (
	"github.com/go-chi/chi/middleware"
	"net/http"
)

// ServidorResponse is the response payload for the Servidor data model.
// See NOTE above in ServidorRequest as well.
//
// In the ServidorResponse object, first a Render() is called on itself,
// then the next field, and so on, all the way down the tree.
// Render is called in top-down order, like a http handler middleware chain.

type ServidorResponse struct {
	*Servidor
	// Returning the Request ID
	ReqID string `json:"reqid"`
}

func NewServidorResponse(servidor *Servidor) *ServidorResponse {
	return &ServidorResponse{Servidor: servidor}
}

func (sr *ServidorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		sr.ReqID = reqID
	}
	return nil
}

type ServidoresResponse struct {
	*Servidores
	// Returning the Request ID
	ReqID string `json:"reqid"`
}

func NewServidoresResponse(servidores *Servidores) *ServidoresResponse {
	return &ServidoresResponse{Servidores: servidores}
}

func (sr *ServidoresResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		sr.ReqID = reqID
	}
	return nil
}
