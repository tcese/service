package agendamento

import (
	"github.com/go-chi/chi/middleware"
	"net/http"
)

type AgendamentoResponse struct {
	*Agendamento
	ReqID string `json:"reqid"`
}

func NewAgendamentoResponse(a *Agendamento) *AgendamentoResponse {
	return &AgendamentoResponse{Agendamento: a}
}

func (sr *AgendamentoResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		sr.ReqID = reqID
	}
	return nil
}

type AgendamentosResponse struct {
	*Agendamentos
	ReqID string `json:"reqid"`
}

func NewAgendamentosResponse(a *Agendamentos) *AgendamentosResponse {
	return &AgendamentosResponse{Agendamentos: a}
}
func (sr *AgendamentosResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		sr.ReqID = reqID
	}
	return nil
}
