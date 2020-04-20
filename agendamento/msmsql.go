package agendamento

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type msmsqlRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewMsmsqlRepository(
	db *sql.DB,
	logger log.Logger,
) Repository {
	return &msmsqlRepository{db: db, logger: logger}
}

func (m msmsqlRepository) BuscarAgendamento(idAgendamento int64, a *Agendamento) error {
	return nil
}

func (m msmsqlRepository) ListarAgendamentos(a *Agendamentos) error {
	return nil
}

func (m msmsqlRepository) InserirAgendamento(a *Agendamento) error {
	return nil
}

func (m msmsqlRepository) AtualizarAgendamento(idAgendamento int64, a *Agendamento) error {
	return nil
}

func (m msmsqlRepository) DeletarAgendamento(idAgendamento int64) error {
	return nil
}
