package main

//use these structs to parse json coming from API
type APIRequest struct {
	Companies []Company `json:"companies"`
	CCID      string    `json:"eventCCID"`
}

type Company struct {
	CompanyID   string `json:"companyID"`
	CompanyType string `json:"companyType"`
}

//use these structs to marshal new data created from struct
type InterestedParties struct {
	Entitlements []IPCompany
}

type DLTIPCompany struct {
	ID            string
	ReceivedDate  string
	RecordDate    string
	Shares        string
	Status        Status
	EffectiveDate string
}

type IPEntitlement struct {
	ID           string
	ReceivedDate string
	RecordDate   string
	Shares       string
	Status       Status
}

type Status struct {
	Code    string
	Message string
}

//fabric values stored structs
type Entitlements struct {
	Entitlements []Entitlement `json:"Entitlements"`
}

type IPCompany struct {
	ID            string `json:"ID"`
	ReceivedDate  string `json:"receivedDate"`
	RecordDate    string `json:"recordDate"`
	Shares        string `json:"shares"`
	Status        Status `json:"status"`
	EffectiveDate string `json:"effectiveDate"`
}

type Entitlement struct {
	ID           string `json:"ID"`
	ReceivedDate string `json:"receivedDate"`
	RecordDate   string `json:"recordDate"`
	Shares       string `json:"shares"`
	Status       Status `json:"status"`
}

//creating new array from api json structs
type IPEntitlementsStr struct {
	Entitlements []Entitlement
}

//updated entitlement structs
type DLTIPEntitlements struct {
	Entitlements []IPCompany
}
