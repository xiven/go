package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type EventCC struct {
}

func main() {
	err := shim.Start(new(EventCC))
	if err != nil {
		fmt.Printf("Error starting Event chaincode: %s", err)
	}
}

func (t *EventCC) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	if function == "addEvent" {
		return t.addEvent(stub, args)
	}

	return nil, nil
}

// Invoke is your entry point to invoke a chaincode function
func (t *EventCC) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addEvent" {
		return t.addEvent(stub, args)
	} else if function == "updateEvent" {
		return t.updateEvent(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *EventCC) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "queryEvent" { //read a variable
		return t.queryEvent(stub, args)
	} else if function == "queryAgenda" {
		return t.queryAgenda(stub, args)
	} else if function == "queryAgendaByProposalText" {
		return t.queryAgendaByProposalText(stub, args)
	} else if function == "queryEventHistory" {
		return t.queryEventHistory(stub, args)
	} else if function == "queryAgendaHistory" {
		return t.queryAgendaHistory(stub, args)
	} else if function == "queryEffectiveDateList" {
		return t.queryEffectiveDateList(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

func (t *EventCC) queryEvent(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = "Event"
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *EventCC) queryAgenda(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = "Agenda"
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *EventCC) queryAgendaByProposalText(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *EventCC) queryEventHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the EffectiveDate to query")
	}

	key = args[0] //EffectiveDate

	//validate string field effective date
	validateEffectiveDate := validate(key, "queryDate", "query effective date", 8)
	if validateEffectiveDate != "true" {
		jsonResp = "{\"Error\":\"" + validateEffectiveDate + "\"}"
		return nil, errors.New(jsonResp)
	}

	valAsbytes, err := stub.GetState("Event_" + key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for Event\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *EventCC) queryAgendaHistory(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the EffectiveDate to query")
	}

	key = args[0] //EffectiveDate
	//validate string field effective date
	validateEffectiveDate := validate(key, "queryDate", "query effective date", 8)
	if validateEffectiveDate != "true" {
		jsonResp = "{\"Error\":\"" + validateEffectiveDate + "\"}"
		return nil, errors.New(jsonResp)
	}

	valAsbytes, err := stub.GetState("Agenda_" + key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for Event\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

func (t *EventCC) queryEffectiveDateList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]

	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for Effective Date List\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
