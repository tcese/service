package servidor

import (
	"time"
)

type Servidores []Servidor

type Servidor struct {
	Matricula    int64      `json:"matricula"`
	Nome         *string    `json:"nome"`
	Sexo         *string    `json:"sexo"`
	Dtnascimento *time.Time `json:"dtnascimento"`
	Empresa      int64      `json:"empresa"`
	Setor        *int       `json:"setor"`
	Fone         *string    `json:"fone"`
	Filial       int64      `json:"filial"`
	Situacao     *string    `json:"situacao"`
	Sangue       *string    `json:"sangue"`
	Doador       *string    `json:"doador"`
	Ec           *string    `json:"ec"`
	Endereco     *string    `json:"endereco"`
	Numero       *string    `json:"numero"`
	Bairro       *string    `json:"bairro"`
	Cidade       *int       `json:"cidade"`
	Cep          *int       `json:"cep"`
}

type Perfil struct {
	Matricula       int64  `json:"matricula"`
	FlAdministrador string `json:"nome"`
	FlCadastro      string `json:"sexo"`
	FlConsulta      string `json:"dtnascimento"`
	FlRelatorio     string `json:"empresa"`
	FlAgendamento   string `json:"setor"`
	FlMedico        string `json:"fone"`
	FlDentista      string `json:"filial"`
	FlPsicologo     string `json:"situacao"`
	FlEnfermeiro    string `json:"sangue"`
	FlAssistSocial  string `json:"doador"`
}

// Retorna uma lista de Servidores já povoada para testes
func NewMockServidores() *Servidores {
	sangue := "A+"
	texto := "Teste"
	numero := 123
	data := time.Now()
	return &Servidores{
		Servidor{1, &texto, &texto, &data, 1, &numero, &texto, 1, &texto, &sangue, &texto, &texto, &texto, &texto, &texto, &numero, &numero},
	}
}
