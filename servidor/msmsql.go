package servidor

import (
	"database/sql"
	"fmt"
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
	server string,
	schema string,
	user string,
	password string,
	logger log.Logger,
) Repository {

	dataSource := fmt.Sprintf("server=%s;user id=%s;password=%s;", server, user, password)
	db, err := sql.Open("mssql", dataSource)
	if err != nil {
		logger.Fatalln(" error open db:", err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		logger.Fatalln("cannot connect: ", err.Error())
		return nil
	}

	//defer db.Close()

	var (
		sqlversion string
	)

	rows, err := db.Query("select @@version")
	if err != nil {
		logger.Fatalln("error retrieving msmsql version", err.Error())
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&sqlversion)
		if err != nil {
			logger.Fatal(err)
		}
		logger.Println(sqlversion)
	}

	return &msmsqlRepository{db: db, server: server, schema: schema, user: user, logger: logger}
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
		m.logger.Println("FOUND: ", s.Matricula)
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
