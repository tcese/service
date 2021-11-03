package agendamento

type Repository interface {
	BuscarAgendamento(idAgendamento int64, a *Agendamento) error
	ListarAgendamentos(a *Agendamentos) error
	InserirAgendamento(a *Agendamento) error
	AtualizarAgendamento(idAgendamento int64, a *Agendamento) error
	DeletarAgendamento(idAgendamento int64) error
}
