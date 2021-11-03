package agendamento

import "log"

type Service interface {
	BuscarAgendamento(idAgendamento int64) (Agendamento, error)
	ListarAgendamentos() (Agendamentos, error)
	InserirAgendamento(a *Agendamento) error
	AtualizarAgendamento(idAgendamento int64, a *Agendamento) error
	DeletarAgendamento(idAgendamento int64) error
}

type internalService struct {
	repository Repository
	logger     log.Logger
}

func (i internalService) BuscarAgendamento(idAgendamento int64) (Agendamento, error) {
	a := Agendamento{}
	i.logger.Printf("buscando agendamento: %d", idAgendamento)
	err := i.repository.BuscarAgendamento(idAgendamento, &a)
	i.logger.Printf("encontrado agendamento: %d - m: %d", a.IdAgendamento, a.Matricula)
	return a, err
}

func (i internalService) ListarAgendamentos() (Agendamentos, error) {
	as := Agendamentos{}
	i.logger.Printf("buscando servidores")
	err := i.repository.ListarAgendamentos(&as)
	i.logger.Printf("buscando servidores len: %d", len(as))
	return as, err
}

func (i internalService) InserirAgendamento(a *Agendamento) error {
	i.logger.Printf("inserindo agendamento do servidor: %d", a.Matricula)
	err := i.repository.InserirAgendamento(a)
	i.logger.Printf("agendamento inserido id: %d", a.IdAgendamento)
	return err
}

func (i internalService) AtualizarAgendamento(idAgendamento int64, a *Agendamento) error {
	i.logger.Printf("atualizando agendamento id: %d", idAgendamento)
	err := i.repository.AtualizarAgendamento(idAgendamento, a)
	i.logger.Printf("agendamento atualizado id: %d - m: %d", a.IdAgendamento, a.Matricula)
	return err
}

func (i internalService) DeletarAgendamento(idAgendamento int64) error {
	i.logger.Printf("deletando agendamento id: %d", idAgendamento)
	err := i.repository.DeletarAgendamento(idAgendamento)
	i.logger.Printf("removido agendamento id: %d", idAgendamento)
	return err
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
