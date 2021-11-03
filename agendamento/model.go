package agendamento

import "time"

type Agendamentos []Agendamento

type Agendamento struct {
	IdAgendamento         int64      `json:"idagendamento"`
	Filial                int64      `json:"filial"`
	Matricula             int64      `json:"matricula"`
	Empresa               int64      `json:"empresa"`
	Sequencia             *int       `json:"sequencia"`
	DhAgendamento         time.Time  `json:"dhagendamento"`
	IdPrestador           int        `json:"idprestador"`
	FlCancelado           *string    `json:"flcancelado"` // char(1) 0 or 1
	IdReagendamento       *int64     `json:"idreagendamento"`
	IdMotivoReagendamento *int       `json:"idmotivoreagendamento"`
	TpFormaAgendamento    string     `json:"tpformaagendamento"` // P ou T
	DhCancelamento        *time.Time `json:"DhCancelamento"`     // ? can be null...
	FlAtendido            *string    `json:"Flatendido"`         // char(1) 0 or 1
}

// Retorna uma lista de Agendamentos já povoada para testes
func NewMockAgendamentos() *Agendamentos {
	return &Agendamentos{
		Agendamento{0, 0, 1, 0, new(int), time.Date(2020, time.July, 11, 9, 30, 0, 0, time.UTC), 0, new(string), new(int64), new(int), "P", nil, new(string)},
		Agendamento{1, 0, 3, 0, new(int), time.Date(2020, time.July, 15, 9, 00, 0, 0, time.UTC), 0, new(string), new(int64), new(int), "P", nil, new(string)},
		Agendamento{2, 0, 2, 0, new(int), time.Date(2020, time.July, 20, 10, 30, 0, 0, time.UTC), 0, new(string), new(int64), new(int), "P", nil, new(string)},
	}
}
