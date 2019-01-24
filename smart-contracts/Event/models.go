package main

type Event struct {
	ID         string   `json:"id"`
	RecordDate string   `json:"recordDate"`
	EventDate  string   `json:"eventDate"`
	EventType  string   `json:"eventType"`
	Status     string   `json:"status"`
	Agenda     []Agenda `json:"agenda"`
}

type DLTEvent struct {
	ID            string `json:"id"`
	RecordDate    string `json:"recordDate"`
	EventDate     string `json:"eventDate"`
	EventType     string `json:"eventType"`
	Status        string `json:"status"`
	AgendaNbr     string `json:"agendaNbr"`
	EffectiveDate string `json:"effectiveDate"`
}

type Agenda struct {
	AgendaNum string     `json:"agendaNum"`
	Proposals []Proposal `json:"proposals"`
}

type DLTAgenda struct {
	AgendaNum     string     `json:"agendaNum"`
	Proposals     []Proposal `json:"proposals"`
	EffectiveDate string     `json:"effective_date"`
}

type Proposal struct {
	ProposalNum  string `json:"proposalNum"`
	ProposalText string `json:"proposalText"`
}
