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

	// Configurando o nome, o caminho e a extensão do arquivo config
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	// Habilitando VIPER para ler as Variáveis de Ambiente
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("viper erro lendo arquivo config: ", err)
		panic("revise o arquivo config")
	}

	// Set undefined variables - NÃO USADO AINDA
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("viper erro ao ler a estrutura do arquivo de configuração: ", err)
		panic("revise a estrutura do arquivo config.")
	}

	fmt.Println("Lendo as variáveis contidas no arquivo de configuração")
	fmt.Println("Database server \t", config.Database.Server)
	fmt.Println("Database name \t", config.Database.Schema)
	fmt.Println("Port \t\t\t", config.Server.Port)
	fmt.Println("Production \t", config.Server.Production)
	fmt.Println("RepMode \t\t", config.RepMode)
	fmt.Println("LogFile \t\t", config.LogFile)

	// Enviando o log para o arquivo config.LogFile
	if len(config.LogFile) > 0 {
		// Enviando os erros para um arquivo
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("error abrindo o arquivo de log: ", err)
		} else {
			defer f.Close()
			if f != nil {
				fmt.Println("configurando como saida do log o arquivo ", config.LogFile)
				log.SetOutput(f)
			}
		}
	}

	// Abrindo a conexão com a Base de Dados
	if config.RepMode == "msmsql" {
		// Conectando com a Base de Dados
		dataSource := fmt.Sprintf("server=%s;user id=%s;password=%s;", config.Database.Server, config.Database.User, config.Database.Password)
		db, err = sql.Open("mssql", dataSource)
		if err != nil {
			log.Fatalln(" erro conectando com a base:", err.Error())
		}

		err = db.Ping()
		if err != nil {
			log.Fatalln("falha na conexão com a base: ", err.Error())
		}

		defer db.Close()

		// Testando a conexão com a base buscando a versão do MSSQL
		rows, err := db.Query("select @@version")
		if err != nil {
			log.Fatalln("erro recuperando a versão da base: ", err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var sqlversion string
			err := rows.Scan(&sqlversion)
			if err != nil {
				log.Fatalln("erro ao ler a versão da base", err)
			}
			log.Println(sqlversion)
		}
	}

	// Subindo o roteador
	r := chi.NewRouter()
	// Anexando os Middlewares
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.RequestLogger(log.Writer()))
	r.Use(chiMiddleware.Recoverer)

	// TODO Rota principal
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Get /")
		w.Write([]byte("=)"))
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
			panic("parametro repomode utilizado nao suportado para o serviço: servidor")
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
		cr := servidor.NewChiController(ss, *logger)
		r.Mount("/servidor/", cr)

	}
	// -------------- END SERVIDOR --------------

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
			panic("parametro repomode utilizado nao suportado para o serviço: servidor")
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
		ar := agendamento.NewChiController(as, *logger)
		r.Mount("/agendamento/", ar)
	}
	// -------------- END SERVIDOR --------------

	// ----------- EXEMPLO DE ROTA LENTA -------------
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("iniciando rota lenta")
		time.Sleep(15 * time.Second)
		fmt.Println("terminando rota lenta")
		w.Write([]byte(fmt.Sprintf("tudo feito!\n")))
	})

	// --------------- SUBINDO O SERVIDOR --------------
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.Port),
		Handler: r,
	}

	go func() {
		fmt.Println("iniciando servidor...") // na saida padrao
		log.Println("iniciando servidor...") // no arquivo
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("err ouvindo:", err) // na saida padrao
			log.Fatalln("err ouvindo:", err) // no arquivo
		}
	}()

	// -------- DERRUBANDO O SERVIDOR EDUCADAMENTE -------
	// signais como ^C
	// recebendo eles e saindo educadamente
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("desligando...") // na saida padrao
	log.Println("desligando...") // no arquivo
	// criar contexto com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("servidor desligado:", err) // na saida padrao
		log.Fatalln("servidor desligado:", err) // no arquivo
	}
	fmt.Println("fim.") // na saida padrao
	log.Println("fim.") // no arquivo
}
