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
	m.logger.Printf("buscando agendamento: %d", idAgendamento)
	rows, err := m.db.Query(
		`SELECT TOP (1) [IDAGENDAMENTO],	[FILIAL], [MATRICULA], [EMPRESA], [SEQUENCIA], [DHAGENDAMENTO],
[IDPRESTADOR], [FLCANCELADO], [IDREAGENDAMENTO], [IDMOTIVOREAGENDAMENTO], [TPFORMAAGENDAMENTO],
[DHCANCELAMENTO], [FLATENDIDO] FROM [dbSCM].[scm].[SCMAGENDAMENTO] WHERE [IDAGENDAMENTO] = ?`, idAgendamento)
	if err != nil {
		m.logger.Println("erro ao executar a query: ", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&a.IdAgendamento, &a.Filial, &a.Matricula, &a.Empresa, &a.Sequencia, &a.DhAgendamento,
			&a.IdPrestador, &a.FlCancelado, &a.IdReagendamento, &a.IdMotivoReagendamento, &a.TpFormaAgendamento,
			&a.DhCancelamento, &a.FlAtendido)
		if err != nil {
			m.logger.Println("erro escaneando a resposta: ", err)
			continue
		}
		m.logger.Println("encontrado: ", a.IdAgendamento)
	}

	return nil
}

func (m msmsqlRepository) ListarAgendamentos(as *Agendamentos) error {
	m.logger.Println("buscando agendamentos")
	rows, err := m.db.Query(
		`SELECT TOP (100) [IDAGENDAMENTO],	[FILIAL], [MATRICULA], [EMPRESA], [SEQUENCIA], [DHAGENDAMENTO],
[IDPRESTADOR], [FLCANCELADO], [IDREAGENDAMENTO], [IDMOTIVOREAGENDAMENTO], [TPFORMAAGENDAMENTO],
[DHCANCELAMENTO], [FLATENDIDO] FROM [dbSCM].[scm].[SCMAGENDAMENTO]`)
	if err != nil {
		m.logger.Println("Cannot query: ", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var a = Agendamento{}
		err = rows.Scan(&a.IdAgendamento, &a.Filial, &a.Matricula, &a.Empresa, &a.Sequencia, &a.DhAgendamento,
			&a.IdPrestador, &a.FlCancelado, &a.IdReagendamento, &a.IdMotivoReagendamento, &a.TpFormaAgendamento,
			&a.DhCancelamento, &a.FlAtendido)
		if err != nil {
			m.logger.Println("erro escaneando a resposta: ", err)
			continue
		}
		*as = append(*as, a)
	}
	m.logger.Println("agendamentos encontrado: ", len(*as))
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
