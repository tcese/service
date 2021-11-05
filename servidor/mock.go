package servidor

import (
	"log"
	"service/tools"
)

type mockRepository struct {
	servidores *Servidores
	logger     log.Logger
}

func NewMockRepository(
	servidores *Servidores,
	logger log.Logger,
) Repository {
	return &mockRepository{servidores: servidores, logger: logger}
}

func (m mockRepository) BuscarServidor(matricula int64) (*Servidor, error) {
	m.logger.Printf("mock buscando servidor: %v", matricula)
	for _, s := range *m.servidores {
		if s.Matricula == matricula {
			m.logger.Printf("mock encontrado servidor: %v", s.Nome)
			return &s, nil
		}
	}
	return nil, &tools.EntityNotFoundError{}
}

func (m mockRepository) ListarServidores() (*Servidores, error) {
	m.logger.Printf("mock buscando servidores len: %v", len(*m.servidores))
	return m.servidores, nil
}
