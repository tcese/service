package main

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"service"
	"service/agendamento"
	"service/middleware"
	"service/servidor"

	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	var config service.Config
	var db *sql.DB

	fmt.Println("Iniciando...")

	// Configurando o nome, o caminho e a extensão do arquivo de configuração
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	// Autorizando o VIPER a ler as Variáveis de Ambiente
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Panicln("erro do viper ao ler o arquivo config.yml: ", err)
	}

	// Set undefined variables - NÃO UTILIZADO
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Panicln("erro do viper ao ler a estrutura do arquivo config.yml: ", err)
	}

	// Enviando o log para o arquivo informado no config.LogFile
	if len(config.LogFile) > 0 {
		// Enviando os erros para um arquivo
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panicln("error ao criar/abrir o arquivo de log: ", err)
		} else {
			defer f.Close()
			if f != nil {
				fmt.Println("configurando como saida do log o arquivo ", config.LogFile)
				log.SetOutput(f)
			}
		}
	}

	log.Println("Iniciando o serviço...")
	log.Println("Variáveis de configuração carregadas do arquivo config.yml")
	log.Println("Database server \t", config.Database.Server)
	log.Println("Database name \t", config.Database.Schema)
	log.Println("Port \t\t\t", config.Server.Port)
	log.Println("Production \t\t", config.Server.Production)
	log.Println("RepMode \t\t", config.RepMode)
	log.Println("LogFile \t\t", config.LogFile)

	// Abrindo a conexão com a Base de Dados
	// Bases implementadas: msmsql ou mock
	if config.RepMode == "msmsql" { // Configurando a base MSMSQL
		// Conectando com a Base de Dados
		dataSource := fmt.Sprintf("server=%s;user id=%s;password=%s;", config.Database.Server, config.Database.User, config.Database.Password)
		db, err = sql.Open("mssql", dataSource)
		if err != nil {
			log.Panicln("erro conectando com a base:", err.Error())
		}

		// Testando a conexão com a base
		err = db.Ping()
		if err != nil {
			log.Panicln("falha na conexão com a base: ", err.Error())
		}

		// Adiar o fechamento da conexão quando essa função 'main' retornar
		defer db.Close()

		// Testando a conexão com a base buscando a versão do MSSQL
		rows, err := db.Query("select @@version")
		if err != nil {
			log.Panicln("erro recuperando a versão da base: ", err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var sqlversion string
			err := rows.Scan(&sqlversion)
			if err != nil {
				log.Panicln("erro ao ler a versão da base", err)
			}
			log.Println("Versão da Base de Dados")
			log.Println(sqlversion)
			log.Println("Fim da Versão da Base de Dados")
		}

	} else if config.RepMode == "mock" { // Caso utilize objetos MOCK simular a base de dados
		log.Println("Utilizando objetos MOCK para simular a base de dados")

	} else { // Base configurada no arquivo config.yml não implementada
		log.Panicln("Base de dados descrita no arquivo config.yml não implementada: ", config.RepMode)
	}

	// Subindo o roteador
	r := chi.NewRouter()
	// Anexando os Middlewares
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogger(log.Writer()))
	r.Use(chiMiddleware.Recoverer)

	// ---------------- ROTA PRINCIPAL ----------------
	// TODO Rota principal
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /")
		w.Write([]byte("=)")) // Retorno ainda não implementado
	})

	// ---------------- SERVIDOR ----------------
	// instanciar repositório
	var sr servidor.Repository
	{
		logger := log.New(log.Writer(), "servidor.Repository ", log.Flags())
		switch config.RepMode {
		case "mock":
			sr = servidor.NewMockRepository(servidor.NewMockServidores(), *logger)
		case "msmsql":
			sr = servidor.NewMsmsqlRepository(db, *logger)
		default:
			log.Panicln("parametro repomode utilizado nao suportado para o serviço: servidor")
		}
	}
	// instanciar servico
	var ss servidor.Service
	{
		logger := log.New(log.Writer(), "servidor.Service ", log.Flags())
		ss = servidor.NewInternalService(sr, *logger)
	}
	// instanciar rota
	{
		logger := log.New(log.Writer(), "servidor.ChiController ", log.Flags())
		r.Mount("/servidor/", servidor.NewChiController(ss, *logger))
	}
	// -------------- FIM SERVIDOR --------------

	// ---------------- AGENDAMENTO ----------------
	// instanciar repositório
	var ar agendamento.Repository
	{
		logger := log.New(log.Writer(), "servidor.Repository ", log.Flags())
		switch config.RepMode {
		case "mock":
			ar = agendamento.NewMockRepository(agendamento.NewMockAgendamentos(), 0, *logger)
		case "msmsql":
			ar = agendamento.NewMsmsqlRepository(db, *logger)
		default:
			log.Panicln("parametro repomode utilizado nao suportado para o serviço: servidor")
		}
	}
	// instanciar servico
	var as agendamento.Service
	{
		logger := log.New(log.Writer(), "servidor.Service ", log.Flags())
		as = agendamento.NewInternalService(ar, *logger)
	}
	// instanciar rota
	{
		logger := log.New(log.Writer(), "servidor.ChiController ", log.Flags())
		r.Mount("/agendamento/", agendamento.NewChiController(as, *logger))
	}
	// -------------- FIM SERVIDOR --------------

	// ----------- EXEMPLO DE ROTA LENTA -------------
	// Esta rota é utilizada para testes
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		log.Println("iniciando rota lenta")
		time.Sleep(15 * time.Second)
		log.Println("terminando rota lenta")
		w.Write([]byte(fmt.Sprintf("tudo feito!\n")))
	})

	// --------------- ABRINDO A ESCUTA DA PORTA config.Server.Port --------------
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.Port),
		Handler: r,
	}
	go func() {
		fmt.Println("Ouvindo a porta ", config.Server.Port) // na saida padrao
		log.Println("Ouvindo a porta ", config.Server.Port) // no arquivo
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("err ao abrir a porta do serviço:", err) // na saida padrao
			log.Panicln("err ao abrir a porta do serviço:", err) // no arquivo
		}
	}()

	// -------- DERRUBANDO O SERVIDOR EDUCADAMENTE -------
	// sinais como ^C, recebendo eles e finalizando educadamente
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("interropendo o serviço...") // na saida padrao
	log.Println("interropendo o serviço...") // no arquivo
	// criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("serviço interrompido:", err) // na saida padrao
		log.Println("serviço interrompido:", err) // no arquivo
	}
	fmt.Println("fim.") // na saida padrao
	log.Println("fim.") // no arquivo
}
