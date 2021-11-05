package servidor

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"service/tools"
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

func (m msmsqlRepository) BuscarServidor(matricula int64) (*Servidor, error) {
	m.logger.Println("buscando servidor: ", matricula)
	rows, err := m.db.Query("SELECT TOP (1) [MATRICULA], [NOME], [SEXO], [DTNASCIMENTO], [EMPRESA], [SETOR], [FONE], [FILIAL], [SITUACAO], [SANGUE], [DOADOR], [EC], [ENDERECO], [NUMERO], [BAIRRO], [CIDADE], [CEP] FROM [dbSCM].[scm].[SCMFUNCIONARIOS] WHERE MATRICULA =  ?", matricula)
	if err != nil {
		return nil, &tools.DataBaseError{err}
	}
	defer rows.Close()
	for rows.Next() {
		s := &Servidor{}
		err = rows.Scan(&s.Matricula, &s.Nome, &s.Sexo, &s.Dtnascimento, &s.Empresa, &s.Setor, &s.Fone, &s.Filial, &s.Situacao, &s.Sangue, &s.Doador, &s.Ec, &s.Endereco, &s.Numero, &s.Bairro, &s.Cidade, &s.Cep)
		if err != nil {
			return nil, &tools.DataBaseError{err}
		}
		return s, nil // Servidor encontrado
	}
	return nil, &tools.EntityNotFoundError{} // Servidor não encontrado
}

func (m msmsqlRepository) ListarServidores() (*Servidores, error) {
	m.logger.Println("buscando lista de servidores ...")
	rows, err := m.db.Query("SELECT TOP (10) [MATRICULA], [NOME], [SEXO], [DTNASCIMENTO], [EMPRESA], [SETOR], [FONE], [FILIAL], [SITUACAO], [SANGUE], [DOADOR], [EC], [ENDERECO], [NUMERO], [BAIRRO], [CIDADE], [CEP] FROM [dbSCM].[scm].[SCMFUNCIONARIOS]") // SELECT s.MATRICULA FROM s dbSCM.SCMFUNCIONARIOS
	if err != nil {
		return nil, &tools.DataBaseError{err}
	}
	defer rows.Close()

	ss := &Servidores{}
	for rows.Next() {
		s := Servidor{}
		err = rows.Scan(&s.Matricula, &s.Nome, &s.Sexo, &s.Dtnascimento, &s.Empresa, &s.Setor, &s.Fone, &s.Filial, &s.Situacao, &s.Sangue, &s.Doador, &s.Ec, &s.Endereco, &s.Numero, &s.Bairro, &s.Cidade, &s.Cep)
		if err != nil {
			return nil, &tools.DataBaseError{err}
		}
		*ss = append(*ss, s)
	}

	return ss, nil
}
