package agendamento

import (
	"errors"
	"fmt"
	"log"
)

type mockRepository struct {
	agendamentos  Agendamentos
	idAgendamento int64
	logger        log.Logger
}

func NewMockRepository(
	agendamentos *Agendamentos,
	idAgendamento int64,
	logger log.Logger,
) Repository {
	if idAgendamento == 0 {
		idAgendamento = int64(len(*agendamentos))
	}
	return &mockRepository{agendamentos: *agendamentos, idAgendamento: idAgendamento, logger: logger}
}

// retorna o indice do objeto ou -1 se não encontrado
func (m mockRepository) buscarAgendamento(idAgendamento int64, a *Agendamento) (error, int) {
	m.logger.Printf("buscando agendamento id: %d", idAgendamento)
	for i, agendamento := range m.agendamentos {
		if agendamento.IdAgendamento == idAgendamento {
			*a = agendamento
			m.logger.Printf("encontrado agendamento: %d - m: %d", agendamento.IdAgendamento, agendamento.Matricula)
			return nil, i
		}
	}
	return nil, -1
}

func (m *mockRepository) atualizarAgendamento(index int, a *Agendamento) error {
	m.logger.Printf("atualizando agendamento: %d", index)
	if index > len(m.agendamentos) {
		return errors.New("index maior que o tamanho da lista de agendamendos para atualizar")
	}
	m.agendamentos = append(m.agendamentos[:index], append([]Agendamento{*a}, m.agendamentos[index:]...)...)
	return nil
}

func (m *mockRepository) deletarAgendamento(index int) error {
	m.logger.Printf("deletando agendamento: %d", index)
	if index > len(m.agendamentos) {
		return fmt.Errorf("index %d maior que o tamanho %d da lista de agendamendos para deletar", index, len(m.agendamentos))
	}
	m.agendamentos = append(m.agendamentos[:index], m.agendamentos[index+1:]...)
	return nil
}

func (m mockRepository) BuscarAgendamento(idAgendamento int64, a *Agendamento) error {
	m.logger.Printf("buscar agendamento: %v", idAgendamento)
	err, _ := m.buscarAgendamento(idAgendamento, a) // descartando o índice
	return err
}

func (m mockRepository) ListarAgendamentos(a *Agendamentos) error {
	m.logger.Printf("mock buscando agendamentos qtd: %v", len(m.agendamentos))
	*a = m.agendamentos
	return nil
}

func (m *mockRepository) InserirAgendamento(a *Agendamento) error {
	a.IdAgendamento = m.idAgendamento
	m.idAgendamento = m.idAgendamento + 1
	m.logger.Printf("mock inserindo agendamento id, len: %d, %d", a.IdAgendamento, len(m.agendamentos))
	m.agendamentos = append(m.agendamentos, *a)
	m.logger.Printf("mock inserindo agendamento id, len: %d, %d", a.IdAgendamento, len(m.agendamentos))
	return nil
}

func (m *mockRepository) AtualizarAgendamento(idAgendamento int64, a *Agendamento) error {
	if idAgendamento != a.IdAgendamento {
		m.logger.Printf("idAgendamento diferente: %v != %v", idAgendamento, a.IdAgendamento)
		return fmt.Errorf("idAgendamento %d do caminho diferente do idAgendamento %d do objeto passado", idAgendamento, a.IdAgendamento)
	}
	m.logger.Printf("mock buscando agendamento: %d", idAgendamento)
	ag := Agendamento{}
	err, i := m.buscarAgendamento(idAgendamento, &ag)
	if err != nil {
		return err
	}
	if i == -1 {
		return fmt.Errorf("Agendamento %d não encontrado na base", idAgendamento)
	}
	return m.atualizarAgendamento(i, a) // atualizando a lista
}

func (m *mockRepository) DeletarAgendamento(idAgendamento int64) error {
	m.logger.Printf("mock buscando agendamento: %d", idAgendamento)
	ag := Agendamento{}
	err, i := m.buscarAgendamento(idAgendamento, &ag)
	if err != nil {
		return err
	}
	if i == -1 {
		return fmt.Errorf("Agendamento id %d não encontrado na base", idAgendamento)
	}
	return m.deletarAgendamento(i)
}
