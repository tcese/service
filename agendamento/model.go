package agendamento

import "time"

type Agendamentos []Agendamento

type Agendamento struct {
	IdAgendamento         int64      `json:"idagendamento"`
	Filial                int64      `json:"filial"`
	Matricula             int64      `json:"matricula"`
	Empresa               int64      `json:"empresa"`
	Sequencia             int        `json:"sequencia"`
	DhAgendamento         time.Time  `json:"dhagendamento"`
	IdPrestador           int        `json:"idprestador"`
	FlCancelado           string     `json:"flcancelado"` // char(1) 0 or 1
	IdReagendamento       int64      `json:"idreagendamento"`
	IdMotivoReagendamento int        `json:"idmotivoreagendamento"`
	TpFormaAgendamento    string     `json:"tpformaagendamento"` // P ou T
	DhCancelamento        *time.Time `json:"DhCancelamento"`     // ? can be null...
	FlAtendido            string     `json:"Flatendido"`         // char(1) 0 or 1
}
