package agendamento

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Enviando os erros para um arquivo
	f, err := os.OpenFile("service_test-log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("erro ao abrir arquivo de log: ", err)
	}
	defer f.Close()
	if f != nil {
		log.SetOutput(f)
		// Imprimindo a mensagem na saída padrão
		fmt.Println("Saida do log direcionada para o arquivo service_test-log.txt")
		// Imprimindo a mensagem no arquivo de log de test
	}
	log.New(log.Writer(), "- AgendamentoService ", log.Flags()).Println("rodando os testes...")
	code := m.Run()
	os.Exit(code)
}

func TestBuscarAgendamento(t *testing.T) {
	m := NewMockAgendamentos()
	l := log.New(log.Writer(), "TestBuscarAgendamentos ", log.Flags())
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
	m := NewMockAgendamentos()
	l := log.New(log.Writer(), "TestListarAgendamentos ", log.Flags())
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
	m := NewMockAgendamentos()
	l := log.New(log.Writer(), "TestInserirAgendamento ", log.Flags())
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
	m := NewMockAgendamentos()
	l := log.New(log.Writer(), "TestAtualizarAgendamento ", log.Flags())
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
	m := NewMockAgendamentos()
	l := log.New(log.Writer(), "TestDeletarAgendamento ", log.Flags())
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
