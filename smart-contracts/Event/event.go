package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *EventCC) addEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//var agendaId string
	var jsonResp string

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value")
	}

	//save full event json to dlt for comparison on update event
	fullEvent := stub.PutState("FullEvent", []byte(args[1]))
	if fullEvent != nil {
		return nil, fullEvent
	}

	var event Event
	//map eventjson to model
	err := json.Unmarshal([]byte(args[1]), &event)
	if err != nil {
		return nil, err
	}

	//event validations
	validateID := validate(event.ID, "string", "ID", 10)
	if validateID != "true" {
		jsonResp = "{\"Error\":\"" + validateID + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateRecordDate := validate(event.RecordDate, "date", "record date", 10)
	if validateRecordDate != "true" {
		jsonResp = "{\"Error\":\"" + validateRecordDate + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateEventDate := validate(event.EventDate, "date", "event date", 10)
	if validateEventDate != "true" {
		jsonResp = "{\"Error\":\"" + validateEventDate + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateEventType := validate(event.EventType, "string", "event type", 50)
	if validateEventType != "true" {
		jsonResp = "{\"Error\":\"" + validateEventType + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateStatus := validate(event.Status, "string", "status", 20)
	if validateStatus != "true" {
		jsonResp = "{\"Error\":\"" + validateStatus + "\"}"
		return nil, errors.New(jsonResp)
	}
	//validate agenda fields
	for i := 0; i < len(event.Agenda); i++ {
		j := event.Agenda[i]
		validateAgendaNum := validate(j.AgendaNum, "string", "AgendaNum", 9)
		if validateAgendaNum != "true" {
			jsonResp = "{\"Error\":\"" + validateAgendaNum + "\"}"
			return nil, errors.New(jsonResp)
		}
	}
	//validate proposal fields
	if len(event.Agenda) > 0 {
		for i := 0; i < len(event.Agenda[0].Proposals); i++ {
			j := event.Agenda[0].Proposals[i]
			validateProposalNum := validate(j.ProposalNum, "string", "ProposalNum", 3)
			if validateProposalNum != "true" {
				jsonResp = "{\"Error\":\"" + validateProposalNum + "\"}"
				return nil, errors.New(jsonResp)
			}
			validateProposalText := validate(j.ProposalText, "string", "ProposalText", 1500)
			if validateProposalText != "true" {
				jsonResp = "{\"Error\":\"" + ProposalText + "\"}"
				return nil, errors.New(jsonResp)
			}
		}
	}

	//end event validations

	//write effectivedate (current datetime) to blockchain
	ef := stub.PutState("EffectiveDate", []byte(currentDate()))
	if ef != nil {
		return nil, ef
	}

	var edList []string
	edList = append(edList, edCurrentDate())
	MEventEfList, _ := json.Marshal(edList)
	EventEfList := stub.PutState("Event_Effective_Date_List", []byte(MEventEfList))
	if EventEfList != nil {
		return nil, EventEfList
	}
	MAgendaEfList, _ := json.Marshal(edList)
	AgendaEfList := stub.PutState("Agenda_Effective_Date_List", []byte(MAgendaEfList))
	if AgendaEfList != nil {
		return nil, AgendaEfList
	}

	AgendaOut, _ := json.Marshal(event.Agenda)
	//add event to blockchain
	err = stub.PutState("Agenda", []byte(AgendaOut))
	if err != nil {
		return nil, err
	}

	err = stub.PutState("Agenda_"+edCurrentDate(), []byte(AgendaOut))
	if err != nil {
		return nil, err
	}

	//set agenda id for event json
	eventStr := `{"ID": "` + event.ID + `", "recordDate": "` + event.RecordDate + `", "eventDate": "` + event.EventDate + `", "eventType": "` + event.EventType + `", "status": "` + event.Status + `", "agenda": ` + string(AgendaOut) + `, "effectiveDate": "` + currentDate() + `"}`

	//add event to blockchain
	err = stub.PutState("Event", []byte(eventStr))
	if err != nil {
		return nil, err
	}

	//add event state with date
	err = stub.PutState("Event_"+edCurrentDate(), []byte(eventStr))
	if err != nil {
		return nil, err
	}

	//add individual proposals to blockchain
	if len(event.Agenda) > 0 {
		for i := 0; i < len(event.Agenda[0].Proposals); i++ {
			j := event.Agenda[0].Proposals[i]
			proposalStr := `{"proposalNumber": "` + j.ProposalNum + `", "proposalText": "` + j.ProposalText + `"}`
			//add agenda to blockchain
			err = stub.PutState(j.ProposalNum, []byte(proposalStr))
			if err != nil {
				return nil, err
			}

		}
	}

	//parse out Agendas and Proposals from Event json
	//This is for getting Agenda by EffectiveDate
	for i := 0; i < len(event.Agenda); i++ {
		j := event.Agenda[i]
		//agendaId = j.AssocAgenda
		propOut, _ := json.Marshal(j.Proposals)
		agendaStr := `{"agendaNum": "` + j.AgendaNum + `", "proposals": ` + string(propOut) + `, "effective_date": "` + currentDate() + `"}`
		//add agenda to blockchain with the current date
		err = stub.PutState("ED_Agenda", []byte(agendaStr))
		if err != nil {
			return nil, err
		}

	}

	return nil, nil
}

func (t *EventCC) updateEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. The eventJson")
	}

	//Get current event
	eventBytes, err := stub.GetState("FullEvent")
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for FullEvent\"}"
		return nil, errors.New(jsonResp)
	}

	//Remove '\' from State to enable string comparison
	str := []byte(eventBytes)
	eventEscaped := bytes.Replace(str, []byte("\\"), []byte(""), -1)

	//unmarshal json from fabric
	var event Event
	//map to model
	err1 := json.Unmarshal([]byte(eventEscaped), &event)

	if err1 != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal json from fabric\"}"
		return nil, errors.New(jsonResp)
	}

	fabEvent, fabEventErr := json.Marshal(event)
	if fabEventErr != nil {
		jsonResp = "{\"Error\":\"Failed to marshal json from fabric\"}"
		return nil, errors.New(jsonResp)
	}

	//unmarshal json from api
	var apiEvent Event
	//map to model
	err2 := json.Unmarshal([]byte(args[0]), &apiEvent)

	if err2 != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal json from api\"}"
		return nil, errors.New(jsonResp)
	}

	//event validations
	validateID := validate(apiEvent.ID, "string", "ID", 10)
	if validateID != "true" {
		jsonResp = "{\"Error\":\"" + validateID + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateRecordDate := validate(apiEvent.RecordDate, "date", "record date", 10)
	if validateRecordDate != "true" {
		jsonResp = "{\"Error\":\"" + validateRecordDate + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateEventDate := validate(apiEvent.EventDate, "date", "event date", 10)
	if validateEventDate != "true" {
		jsonResp = "{\"Error\":\"" + validateEventDate + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateEventType := validate(apiEvent.EventType, "string", "event type", 50)
	if validateEventType != "true" {
		jsonResp = "{\"Error\":\"" + validateEventType + "\"}"
		return nil, errors.New(jsonResp)
	}
	validateStatus := validate(apiEvent.Status, "string", "status", 20)
	if validateStatus != "true" {
		jsonResp = "{\"Error\":\"" + validateStatus + "\"}"
		return nil, errors.New(jsonResp)
	}
	//validate agenda fields
	for i := 0; i < len(apiEvent.Agenda); i++ {
		j := apiEvent.Agenda[i]
		validateAgendaNum := validate(j.AgendaNum, "string", "AgendaNum", 9)
		if validateAgendaNum != "true" {
			jsonResp = "{\"Error\":\"" + validateAgendaNum + "\"}"
			return nil, errors.New(jsonResp)
		}
	}
	//validate proposal fields
	if len(apiEvent.Agenda) > 0 {
		for i := 0; i < len(apiEvent.Agenda[0].Proposals); i++ {
			j := apiEvent.Agenda[0].Proposals[i]
			validateProposalNumber := validate(j.ProposalNum, "string", "ProposalNum", 3)
			if validateProposalNumber != "true" {
				jsonResp = "{\"Error\":\"" + validateProposalNumber + "\"}"
				return nil, errors.New(jsonResp)
			}
			validateProposalText := validate(j.ProposalText, "string", "ProposalText", 1500)
			if validateProposalText != "true" {
				jsonResp = "{\"Error\":\"" + validateProposalText + "\"}"
				return nil, errors.New(jsonResp)
			}
		}
	}

	//end event validations

	mrshlEvent, mrshlEventErr := json.Marshal(apiEvent)
	if mrshlEventErr != nil {
		jsonResp = "{\"Error\":\"Failed to marshal json from api\"}"
		return nil, errors.New(jsonResp)
	}

	//compare latest event data with dlt event data
	if string(mrshlEvent) != string(fabEvent) {
		fmt.Println("Event do not match. Update!")
		//save full event json to dlt for comparison on update event
		fullUpdatedEvent := stub.PutState("FullEvent", []byte(mrshlEvent))
		if fullUpdatedEvent != nil {
			return nil, fullUpdatedEvent
		}
		//write effectivedate (current datetime) to blockchain
		ef := stub.PutState("EffectiveDate", []byte(currentDate()))
		if ef != nil {
			return nil, ef
		}

		var edList []string
		//get the effective date list
		GSEventEfList, gsEvErr := stub.GetState("Event_Effective_Date_List")
		if gsEvErr != nil {
			return nil, gsEvErr
		}
		//unmarshal event effective date list
		umEventEfList := json.Unmarshal([]byte(GSEventEfList), &edList)
		if umEventEfList != nil {
			return nil, umEventEfList
		}
		//get the effective date list
		GSAgendaEfList, gsAgErr := stub.GetState("Event_Effective_Date_List")
		if gsAgErr != nil {
			return nil, gsAgErr
		}
		//unmarshal agenda effective date list
		GSAgendaEfList, gsAgEfErr := stub.GetState("Event_Effective_Date_List")
		if gsAgEfErr != nil {
			return nil, gsAgEfErr
		}
		umAgendaEfList := json.Unmarshal([]byte(GSAgendaEfList), &edList)
		if umAgendaEfList != nil {
			return nil, umAgendaEfList
		}
		//add the current effective date to array list
		edList = append(edList, edCurrentDate())

		//marshal lists
		EventEffDateList, _ := json.Marshal(edList)
		AgendaEffDateList, _ := json.Marshal(edList)

		EventEfList := stub.PutState("Event_Effective_Date_List", []byte(EventEffDateList))
		if EventEfList != nil {
			return nil, EventEfList
		}

		AgendaEfList := stub.PutState("Agenda_Effective_Date_List", []byte(AgendaEffDateList))
		if AgendaEfList != nil {
			return nil, AgendaEfList
		}

		AgendaOut, _ := json.Marshal(apiEvent.Agenda)
		//add event to blockchain
		err = stub.PutState("Agenda", []byte(AgendaOut))
		if err != nil {
			return nil, err
		}

		//add index event for queryAgendaHistory
		err = stub.PutState("Agenda_"+edCurrentDate(), []byte(AgendaOut))
		if err != nil {
			return nil, err
		}

		//set agenda id for event json
		eventStr := `{"ID": "` + apiEvent.ID + `", "recordDate": "` + apiEvent.RecordDate + `", "eventDate": "` + apiEvent.EventDate + `", "eventType": "` + apiEvent.EventType + `", "status": "` + apiEvent.Status + `", "agenda": ` + string(AgendaOut) + `, "effectiveDate": "` + currentDate() + `"}`

		//add event to blockchain
		err = stub.PutState("Event", []byte(eventStr))
		if err != nil {
			return nil, err
		}

		//add index event for queryEventHistory
		err = stub.PutState("Event_"+edCurrentDate(), []byte(eventStr))
		if err != nil {
			return nil, err
		}

		//add individual proposals to blockchain
		if len(apiEvent.Agenda) > 0 {
			for i := 0; i < len(apiEvent.Agenda[0].Proposals); i++ {
				j := apiEvent.Agenda[0].Proposals[i]
				dirOut, _ := json.Marshal(j.Director)
				proposalStr := `{"proposalNumber": "` + j.ProposalNum + `", "proposalText": "` + j.ProposalText + `"}`
				//add agenda to blockchain
				err = stub.PutState(j.ProposalNum, []byte(proposalStr))
				if err != nil {
					return nil, err
				}
			}
		}

		//parse out Agendas and Proposals from Event json
		//This is for getting Agenda by EffectiveDate
		for i := 0; i < len(apiEvent.Agenda); i++ {
			j := apiEvent.Agenda[i]
			//agendaId = j.AssocAgenda
			propOut, _ := json.Marshal(j.Proposals)
			agendaStr := `{"agendaNum": "` + j.AgendaNum + `", "proposals": ` + string(propOut) + `, "effective_date": "` + currentDate() + `"}`
			//add agenda to blockchain
			err = stub.PutState("ED_Agenda", []byte(agendaStr))
			if err != nil {
				return nil, err
			}
		}

		return []byte("event updated!"), nil
	} else {
		return []byte("event data matches, no update made!"), nil
	}

	return nil, nil
}
