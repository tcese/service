package agendamento

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"log"
	"net/http"
	"strconv"
)

func NewChiController(
	service Service,
	logger log.Logger,
) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("-> agendamentos/")

		as, err := service.ListarAgendamentos()
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		logger.Printf("agendamentos encontrados: %d", len(as))

		if err := render.Render(w, r, NewAgendamentosResponse(&as)); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/{idAgendamento}", func(w http.ResponseWriter, r *http.Request) {
		idAgendamento := chi.URLParam(r, "idAgendamento")
		logger.Printf("buscando agendamento id: %s", idAgendamento)

		id, err := strconv.ParseInt(idAgendamento, 10, 64)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		a, err := service.BuscarAgendamento(id)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var resp *AgendamentoResponse
		if a.IdAgendamento > 0 {
			logger.Printf("encontrado agendamento: %d", a.IdAgendamento)
			w.WriteHeader(http.StatusOK)
			resp = NewAgendamentoResponse(&a)
		} else {
			w.WriteHeader(http.StatusNotFound)
			resp = NewAgendamentoResponse(nil)
		}

		if err := render.Render(w, r, resp); err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("inserindo novo agendamento...")
		var a Agendamento
		err := render.DecodeJSON(r.Body, &a)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = service.InserirAgendamento(&a)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Put("/{idAgendamento}", func(w http.ResponseWriter, r *http.Request) {
		idAgendamento := chi.URLParam(r, "idAgendamento")
		logger.Printf("atualizando agendamento id: %s", idAgendamento)
		id, err := strconv.ParseInt(idAgendamento, 10, 64)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var a Agendamento
		err = render.DecodeJSON(r.Body, &a)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = service.AtualizarAgendamento(id, &a)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Delete("/{idAgendamento}", func(w http.ResponseWriter, r *http.Request) {
		idAgendamento := chi.URLParam(r, "idAgendamento")
		logger.Printf("removendo agendamento id: %s", idAgendamento)
		id, err := strconv.ParseInt(idAgendamento, 10, 64)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		err = service.DeletarAgendamento(id)
		if err != nil {
			logger.Println(err)
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
