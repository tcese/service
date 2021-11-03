package servidor

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

type msmsqlRepository struct {
	db     *sql.DB
	server string
	schema string
	user   string
	logger log.Logger
}

func NewMsmsqlRepository(
	db *sql.DB,
	logger log.Logger,
) Repository {
	return &msmsqlRepository{db: db, logger: logger}
}

func (m msmsqlRepository) BuscarServidor(matricula int64, s *Servidor) error {
	m.logger.Printf("buscando servidor: %v", matricula)

	rows, err := m.db.Query("SELECT TOP (1) [MATRICULA] FROM [dbSCM].[scm].[SCMFUNCIONARIOS] WHERE MATRICULA =  ?", matricula)
	if err != nil {
		m.logger.Println("Cannot query: ", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&s.Matricula)
		if err != nil {
			m.logger.Println("error scanning: ", err)
			continue
		}
		m.logger.Println("Servidor encontrado: ", s.Matricula)
	}

	return nil
}

func (m msmsqlRepository) ListarServidores(s *Servidores) error {
	m.logger.Println("mock buscando servidores ...")

	rows, err := m.db.Query("SELECT s.MATRICULA FROM s dbSCM.SCMFUNCIONARIOS")
	if err != nil {
		m.logger.Println("Cannot query: ", err.Error())
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var val []interface{}
		err = rows.Scan(val...)
		if err != nil {
			m.logger.Println(err)
			continue
		}
		m.logger.Println(val)
	}

	return nil
}
