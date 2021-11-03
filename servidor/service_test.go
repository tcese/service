package servidor

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"service"
	"testing"
)

var config service.Config

func TestMain(m *testing.M) {
	// Enviando os erros para um arquivo
	f, err := os.OpenFile("service_test-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	if f != nil {
		log.SetOutput(f)
		log.Println("seting file as log output..")
	}
	log.Println("running tests")

	// Set the file name of the configurations file
	viper.SetConfigName("config_test")
	// Set the path to look for the configurations file
	viper.AddConfigPath("../.")
	// The type of the config file
	viper.SetConfigType("yml")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper error reading config file: %v", err)
		panic("review the config.yml file.")
	}

	// Set undefined variables
	//viper.SetDefault("database.dbname", "test_db")

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("viper unable to decode into struct: %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("DB server is\t", config.Database.Server)
	fmt.Println("DB name is\t", config.Database.Schema)
	fmt.Println("Port is\t\t\t", config.Server.Port)
	fmt.Println("Production is\t", config.Server.Production)
	fmt.Println("RepMode is\t\t", config.RepMode)
	fmt.Println("LogFile is\t\t", config.LogFile)

	code := m.Run()
	os.Exit(code)
}

func TestBuscarServidor(t *testing.T) {
	l := log.New(log.Writer(), "TestBuscarServidor ", log.Flags())
	l.Println("TODO")
}

func TestListarServidores(t *testing.T) {
}
