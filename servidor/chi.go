package servidor

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
		logger.Println("buscando servidores")

		ss, err := service.ListarServidores()
		if err != nil { // || ss == nil
			logger.Println(err)
			http.Error(w, "erro ao buscar lista de servidores", http.StatusInternalServerError)
			return
		}

		logger.Println("servidores encontrados: ", len(*ss))
		logger.Printf("servidores: %v", ss)

		if err := render.Render(w, r, NewServidoresResponse(ss)); err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/{matricula}", func(w http.ResponseWriter, r *http.Request) {
		matricula := chi.URLParam(r, "matricula")
		logger.Println("buscando servidor: ", matricula)

		m, err := strconv.ParseInt(matricula, 10, 64)
		if err != nil {
			logger.Println(err)
			http.Error(w, "o parâmetro de matrícula informado deve ser um número inteiro", http.StatusBadRequest) // http.StatusText(http.StatusBadRequest)
			return
		}

		s, err := service.BuscarServidor(m)
		if err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if s == nil {
			logger.Println("servidor não encontrado com matricula: ", matricula)
			http.Error(w, "servidor não encontrado com essa matrícula", http.StatusNotFound)
			return
		}

		logger.Printf("encontrado servidor: %v", s)

		if err := render.Render(w, r, NewServidorResponse(s)); err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			//render.Render(w, r, ErrRender(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
