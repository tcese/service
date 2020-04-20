package servidor

import (
	"log"
)

type mockRepository struct {
	servidores Servidores
	logger     log.Logger
}

func NewMockRepository(
	servidores *Servidores,
	logger log.Logger,
) Repository {
	return &mockRepository{servidores: *servidores, logger: logger}
}

func (m mockRepository) BuscarServidor(matricula int64, s *Servidor) error {
	m.logger.Printf("buscando servidor: %v", matricula)
	for _, servidor := range m.servidores {
		if servidor.Matricula == matricula {
			*s = servidor
			m.logger.Printf("xencontrado servidor: %v", servidor.Nome)
			return nil
		}
	}
	return nil
}

func (m mockRepository) ListarServidores(s *Servidores) error {
	m.logger.Printf("mock buscando servidores len: %v", len(m.servidores))
	*s = m.servidores
	return nil
}
