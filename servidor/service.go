package servidor

import "log"

type Service interface {
	BuscarServidor(matricula int64) (*Servidor, error)
	ListarServidores() (*Servidores, error)
}

type internalService struct {
	repository Repository
	logger     log.Logger
}

func (i internalService) BuscarServidor(matricula int64) (*Servidor, error) {
	return i.repository.BuscarServidor(matricula)
}

func (i internalService) ListarServidores() (*Servidores, error) {
	return i.repository.ListarServidores()
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
