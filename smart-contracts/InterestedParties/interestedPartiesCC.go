package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type IPChaincode struct {
}

func main() {
	err := shim.Start(new(IPChaincode))
	if err != nil {
		fmt.Printf("Error starting IP chaincode: %s", err)
	}
}

// Init resets all the things
func (t *IPChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}
	if function == "addInterestedParties" {
		return t.addInterestedParties(stub, args)
	}

	return nil, nil
}

// Invoke is your entry point to invoke a chaincode function
func (t *IPChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "addInterestedParties" {
		return t.addInterestedParties(stub, args)
	} else if function == "updateInterestedParties" {
		return t.updateInterestedParties(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *IPChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "queryInterestedParties" {
		return t.queryInterestedParties(stub, args)
	} else if function == "queryClient" {
		return t.queryClient(stub, args)
	} else if function == "queryInterestedPartiesByDate" {
		return t.queryInterestedPartiesByDate(stub, args)
	} else if function == "queryEffectiveDateList" {
		return t.queryEffectiveDateList(stub, args)
	}

	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query: " + function)
}

// read - query function to read key/value pair
func (t *IPChaincode) queryInterestedParties(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *IPChaincode) queryInterestedPartiesByDate(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var effectiveDate, companyID, jsonResp string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	//set effective date param
	effectiveDate = args[0]
	companyID = args[1]

	//validate string field effective date
	validateEffectiveDate := validate(effectiveDate, "queryDate", "query effective date", 8)
	if validateEffectiveDate != "true" {
		jsonResp = "{\"Error\":\"" + validateEffectiveDate + "\"}"
		return nil, errors.New(jsonResp)
	}
	//validate string field companyId
	validateCompanyID := validate(companyID, "string", "confirm", 6)
	if validateCompanyID != "true" {
		jsonResp = "{\"Error\":\"" + validateCompanyID + "\"}"
		return nil, errors.New(jsonResp)
	}

	ipBytes, err := stub.GetState(companyID + "_" + effectiveDate)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for interested parties\"}"
		return nil, errors.New(jsonResp)
	}

	return ipBytes, nil
}

func (t *IPChaincode) queryClient(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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

func (t *IPChaincode) queryEffectiveDateList(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
