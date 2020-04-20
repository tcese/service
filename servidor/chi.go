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
		logger.Printf("buscando servidores")

		s, err := service.ListarServidores()
		if err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		logger.Printf("encontrado servidores")

		// Assume if we've reach this far, we can access the article
		// context because this handler is a child of the ArticleCtx
		// middleware. The worst case, the recoverer middleware will save us.
		//article := r.Context().Value("servidor").(*Servidor)
		if err := render.Render(w, r, NewServidoresResponse(&s)); err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			//render.Render(w, r, ErrRender(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/{matricula}", func(w http.ResponseWriter, r *http.Request) {
		matricula := chi.URLParam(r, "matricula")
		logger.Printf("buscando servidor:%v", matricula)

		m, err := strconv.ParseInt(matricula, 10, 64)
		if err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		s, err := service.BuscarServidor(m)
		if err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		logger.Printf("encontrado servidor: %v", s.Nome)

		// Assume if we've reach this far, we can access the SErvidor
		// context because this handler is a child of the ServidorCtx
		// middleware. The worst case, the recoverer middleware will save us.
		//article := r.Context().Value("servidor").(*Servidor)
		if err := render.Render(w, r, NewServidorResponse(&s)); err != nil {
			logger.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			//render.Render(w, r, ErrRender(err))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
