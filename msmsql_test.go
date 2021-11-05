package service

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"service/servidor"
	"testing"
)

var config Config
var db *sql.DB
var logFileName = "msmsql_test-log.txt"

func TestMain(m *testing.M) {
	log.SetPrefix("MSMSQLTest ")

	// Enviando os erros para o arquivo logFileName
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil || f == nil {
		log.Panicln("erro ao abrir arquivo de log", logFileName, err)
	}
	defer f.Close()
	log.SetOutput(f)

	fmt.Println("saida do log direcionada para o arquivo", logFileName)
	log.Println("iniciando os testes...")

	// Informando o nome do arquivo de configuração de testes
	viper.SetConfigName("config_test")
	// Informando o caminho para procurar pelo arquivo de configuração de testes
	viper.AddConfigPath(".")
	// Informando o tipo do arquivo de configuração testes
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln("erro do viper ao ler o arquivo de teste informado", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Panicln("erro do viper ao mapear a estrutura do arquivo de configuração informado", err)
	}

	log.Println("variaveis usadas ne teste informada no arquivo de configuração", logFileName)
	log.Println("DB server is\t\t", config.Database.Server)
	log.Println("DB name is\t\t", config.Database.Schema)
	log.Println("Port is\t\t\t", config.Server.Port)
	log.Println("Production is\t", config.Server.Production)
	log.Println("RepMode is\t\t", config.RepMode)

	dataSource := fmt.Sprintf("server=%s;user id=%s;password=%s;", config.Database.Server, config.Database.User, config.Database.Password)
	db, err = sql.Open("mssql", dataSource)
	if err != nil {
		log.Panicln("erro ao abrir a conexão com a base de dados", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Panicln("erro ao pingar a base de dados", err.Error())
	}
	defer db.Close()

	rows, err := db.Query("select @@version")
	if err != nil {
		log.Panicln("erro ao recuperar a versão da base de dados", err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		var sqlversion string
		err := rows.Scan(&sqlversion)
		if err != nil {
			log.Panicln(err)
		}
		log.Println("versão da base de dados:\n", sqlversion)
	}

	log.Println("rodando os testes...")
	code := m.Run()
	os.Exit(code)
}

func TestDatabaseServidor(t *testing.T) {

	sr := servidor.NewMsmsqlRepository(db, *log.New(log.Writer(), "TestServidorNewMsmsqlRepository ", log.Flags()))
	ss := servidor.NewInternalService(sr, *log.New(log.Writer(), "TestServidorNewInternalService ", log.Flags()))

	s, err := ss.BuscarServidor(2218)
	if err != nil {
		t.Error("error buscando servidores: ", err)
		return
	}
	if s == nil || s.Matricula != 2218 {
		//t.Errorf("error na matricula %T: %v", s.Matricula, s.Matricula)
		t.Errorf("servidor %v", s)
		t.Error("servidor nao encontrado")
		return
	}
}

func TestDatabaseServidores(t *testing.T) {

	sr := servidor.NewMsmsqlRepository(db, *log.New(log.Writer(), "TestServidoresNewMsmsqlRepository ", log.Flags()))
	ss := servidor.NewInternalService(sr, *log.New(log.Writer(), "TestServidoresNewInternalService ", log.Flags()))

	s, err := ss.ListarServidores()
	if err != nil {
		t.Error("error buscando servidores: ", err)
		return
	}
	if len(*s) < 1 {
		t.Error("a base de dados não retornou nenhum servidor na lsita")
		return
	}
}
