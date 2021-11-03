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

func TestMain(m *testing.M) {
	l := log.New(log.Writer(), "- MSMSQL ", log.Flags())

	// Enviando os erros para um arquivo
	f, err := os.OpenFile("msmsql_test-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("erro ao abrir arquivo de log: ", err)
	}
	defer f.Close()
	if f != nil {
		log.SetOutput(f)
		// Imprimindo a mensagem na saída padrão
		fmt.Println("Saida do log direcionada para o arquivo service_test-log.txt")
	}

	// Set the file name of the configurations file
	viper.SetConfigName("config_test")
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	// The type of the config file
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		l.Fatalln("viper erro lendo arquivo config: ", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		l.Fatalln("viper erro lendo a estrutura do arquivo config: ", err)
	}

	fmt.Println("Variaveis usadas no teste do banco de dados...")
	fmt.Println("DB server is\t", config.Database.Server)
	fmt.Println("DB name is\t", config.Database.Schema)
	fmt.Println("Port is\t\t\t", config.Server.Port)
	fmt.Println("Production is\t", config.Server.Production)
	fmt.Println("RepMode is\t\t", config.RepMode)
	fmt.Println("LogFile is\t\t", config.LogFile)

	dataSource := fmt.Sprintf("server=%s;user id=%s;password=%s;", config.Database.Server, config.Database.User, config.Database.Password)
	db, err = sql.Open("mssql", dataSource)
	if err != nil {
		l.Fatalln(" error open db:", err.Error())
	}

	err = db.Ping()
	if err != nil {
		l.Fatalln("cannot connect: ", err.Error())
	}

	//defer db.Close()

	var (
		sqlversion string
	)

	rows, err := db.Query("select @@version")
	if err != nil {
		l.Fatalln("error retrieving msmsql version", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&sqlversion)
		if err != nil {
			l.Fatal(err)
		}
		l.Println(sqlversion)
	}

	l.Println("rodando os testes...")
	code := m.Run()
	os.Exit(code)
}

func TestDatabaseServidor(t *testing.T) {

	l := log.New(log.Writer(), "TestDatabaseServidor ", log.Flags())

	sr := servidor.NewMsmsqlRepository(db, *l)
	if sr == nil {
		t.Error("erro ao abrir o servidor")
	}

	ss := servidor.NewInternalService(sr, *l)

	s, err := ss.BuscarServidor(2218)
	if err != nil {
		t.Error("error buscando servidores: ", err)
		return
	}
	if s.Matricula != 2218 {
		t.Errorf("error na matricula %T: %v", s.Matricula, s.Matricula)
		return
	}
}
