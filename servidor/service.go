package servidor

import "log"

type Service interface {
	BuscarServidor(matricula int64) (Servidor, error)
	ListarServidores() (Servidores, error)
}

type internalService struct {
	repository Repository
	logger     log.Logger
}

func (i internalService) BuscarServidor(matricula int64) (Servidor, error) {
	s := Servidor{}
	i.logger.Printf("buscando servidor:%v", matricula)
	err := i.repository.BuscarServidor(matricula, &s)
	i.logger.Printf("1encontrado servidor: %v", s.Nome)
	return s, err
}

func (i internalService) ListarServidores() (Servidores, error) {
	s := Servidores{}
	i.logger.Printf("buscando servidores")
	err := i.repository.ListarServidores(&s)
	i.logger.Printf("buscando servidores len: %v", len(s))
	return s, err
}

// NewInternalService
func NewInternalService(
	repository Repository,
	logger log.Logger,
) Service {
	return &internalService{
		repository: repository,
		logger:     logger,
	}
}
