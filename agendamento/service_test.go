package agendamento

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"service"
	"testing"
	"time"
)

var config service.Config
var mock = &Agendamentos{
	Agendamento{1, 0, 1, 0, 1, time.Date(2020, time.July, 11, 9, 30, 0, 0, time.UTC), 0, "0", 0, 0, "P", nil, "0"},
	Agendamento{2, 0, 1, 0, 2, time.Date(2020, time.July, 15, 9, 00, 0, 0, time.UTC), 0, "0", 0, 0, "P", nil, "0"},
	Agendamento{3, 0, 1, 0, 3, time.Date(2020, time.July, 20, 10, 30, 0, 0, time.UTC), 0, "0", 0, 0, "P", nil, "0"},
}

func TestMain(m *testing.M) {
	// Enviando os erros para um arquivo
	f, err := os.OpenFile("agendamento_test-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
	fmt.Println("Database server is\t", config.Database.Server)
	fmt.Println("Database name is\t", config.Database.Schema)
	fmt.Println("Port is\t\t\t", config.Server.Port)
	fmt.Println("Production is\t", config.Server.Production)
	fmt.Println("RepMode is\t\t", config.RepMode)
	fmt.Println("LogFile is\t\t", config.LogFile)

	code := m.Run()
	os.Exit(code)
}

func TestBuscarAgendamento(t *testing.T) {
	m := mock
	l := log.New(log.Writer(), "TestBuscarAgendamentos", log.Flags())
	r := NewMockRepository(m, 4, *l)
	if r == nil {
		t.Error("fatal error opening mock db")
	}
	s := NewInternalService(r, *l)

	a, err := s.BuscarAgendamento(2)
	if err != nil {
		t.Error("error buscando agendamento: ", err)
		return
	}
	if a.IdAgendamento != 2 {
		t.Errorf("no agendaemnto id  2 != %d", a.IdAgendamento)
		return
	}
}

func TestListarAgendamentos(t *testing.T) {
	m := mock
	l := log.New(log.Writer(), "TestBuscarAgendamentos", log.Flags())
	r := NewMockRepository(m, 4, *l)
	if r == nil {
		t.Error("fatal error opening mock db")
	}
	s := NewInternalService(r, *l)

	a, err := s.ListarAgendamentos()
	if err != nil {
		t.Error("error buscando agendamento: ", err)
		return
	}
	if len(a) != 3 {
		t.Errorf("tamanho da lista de agendamento  3 != %d", len(a))
		return
	}
}

func TestInserirAgendamento(t *testing.T) {
	m := mock
	l := log.New(log.Writer(), "TestBuscarAgendamentos", log.Flags())
	r := NewMockRepository(m, 4, *l)
	if r == nil {
		t.Error("fatal error opening mock db")
	}
	s := NewInternalService(r, *l)

	na := Agendamento{0, 0, 3, 0, 1, time.Date(2020, time.July, 11, 9, 30, 0, 0, time.UTC), 0, "0", 0, 0, "P", nil, "0"}
	err := s.InserirAgendamento(&na)
	if err != nil {
		t.Error("error inserindo agendamento: ", err)
		return
	}
	if na.IdAgendamento != 4 {
		t.Errorf("id do agendamento inserido 4 != %d", na.IdAgendamento)
		return
	}
}

func TestAtualizarAgendamento(t *testing.T) {
	m := mock
	l := log.New(log.Writer(), "TestBuscarAgendamentos", log.Flags())
	r := NewMockRepository(m, 4, *l)
	if r == nil {
		t.Error("fatal error opening mock db")
	}
	s := NewInternalService(r, *l)

	na := Agendamento{2, 0, 3, 0, 1, time.Date(2020, time.July, 11, 9, 30, 0, 0, time.UTC), 0, "0", 0, 0, "P", nil, "0"}
	err := s.AtualizarAgendamento(2, &na)
	if err != nil {
		t.Error("error inserindo agendamento: ", err)
		return
	}

	a, err := s.BuscarAgendamento(2)
	if err != nil {
		t.Error("erro buscando agendamento: ", err)
		return
	}
	if a.Matricula != 3 {
		t.Errorf("no agendaemnto matricula  3 != %d", a.Matricula)
		return
	}
}

func TestDeletarAgendamento(t *testing.T) {
	m := mock
	l := log.New(log.Writer(), "TestBuscarAgendamentos", log.Flags())
	r := NewMockRepository(m, 4, *l)
	if r == nil {
		t.Error("fatal error opening mock db")
	}
	s := NewInternalService(r, *l)

	err := s.DeletarAgendamento(2)
	if err != nil {
		t.Error("error inserindo agendamento: ", err)
		return
	}

	a, err := s.ListarAgendamentos()
	if err != nil {
		t.Error("erro buscando agendamento: ", err)
		return
	}
	if len(a) != 2 {
		t.Errorf("tamanho da lista de agendamento  2 != %d", len(a))
		return
	}
}
