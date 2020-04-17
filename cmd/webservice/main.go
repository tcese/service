package main

import (
	"context"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/spf13/viper"
	"service"
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

	// Set the file name of the configurations file
	viper.SetConfigName("config")
	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	// The type of the config file
	viper.SetConfigType("yml")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	var config service.Config

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper error reading config file: %v", err)
		panic("review the config.yml file.")
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("viper unable to decode into struct: %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("Database server is\t", config.Database.Server)
	fmt.Println("Database name is\t", config.Database.Schema)
	fmt.Println("Port is\t\t\t", config.Server.Port)
	fmt.Println("Production is\t", config.Server.Production)
	fmt.Println("RepMode is\t\t", config.RepMode)
	fmt.Println("LogFile is\t\t", config.LogFile)

	// Setting log mode
	if len(config.LogFile) > 0 {
		// Enviando os erros para um arquivo
		f, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()
		if f != nil {
			log.Println("seting file as log output")
			log.SetOutput(f)
		} else {
			log.Println("seting terminal as log output")
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
			//sr = servidor.NewMockRepository(nil, *logger)
			panic("parametro repomode=mock nao suportado no momento")
		case "msmsql":
			sr = servidor.NewMsmsqlRepository(config.Database.Server, config.Database.Schema, config.Database.User, config.Database.Password, *logger)
			if sr == nil {
				panic("fatal error opening msmsql db")
			}
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

	// ----------- EXEMPLO DE ROTA LENTA -------------
	r.Get("/slow", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("starting slow")
		time.Sleep(15 * time.Second)
		fmt.Println("finishing slow")
		w.Write([]byte(fmt.Sprintf("all done.\n")))
	})

	// --------------- SUBINDO O SERVIDOR --------------
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.Server.Port),
		Handler: r,
	}

	go func() {
		fmt.Println("server starting...") // log to console
		log.Println("server starting...") // log to file
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("err listen:", err) // log to console
			log.Fatalln("err listen:", err) // log to file
		}
	}()

	// -------- DERRUBANDO O SERVIDOR EDUCADAMENTE -------

	// sig is a ^C, handle it
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("shutting down...") // log to console
	log.Println("shutting down...") // log to file

	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Println("server shutdown:", err) // log to console
		log.Fatalln("server shutdown:", err) // log to file
	}

	fmt.Println("server exiting") // log to console
	log.Println("server exiting") // log to file
}
