package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func (t *IPChaincode) addInterestedParties(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var jsonResp string

	//check for length of arguments passed in
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value")
	}

	var request APIRequest
	//map intParties json to model
	err := json.Unmarshal([]byte(args[1]), &request)

	if err != nil {
		return nil, err
	}

	//validate companies fields
	for i := 0; i < len(request.Companies); i++ {
		j := request.Companies[i]
		validateCompanyID := validate(j.CompanyID, "string", "CompanyID", 6)
		if validateCompanyID != "true" {
			jsonResp = "{\"Error\":\"" + validateCompanyID + "\"}"
			return nil, errors.New(jsonResp)
		}
		validateCompanyType := validate(j.CompanyType, "string", "CompanyType", 10)
		if validateCompanyType != "true" {
			jsonResp = "{\"Error\":\"" + validateCompanyType + "\"}"
			return nil, errors.New(jsonResp)
		}
	}
	//validate CCID
	validateCCID := validate(request.CCID, "string", "CCID", 64)
	if validateCCID != "true" {
		jsonResp = "{\"Error\":\"" + validateCCID + "\"}"
		return nil, errors.New(jsonResp)
	}

	//declare array list of client Ids
	var list_of_clientIds []string

	//set arrays for newly formed JSON string for Companies List
	entArray := []IPCompany{}
	entList := IPEntitlements{entArray}

	//loop through company list in interested parties json
	for i := 0; i < len(request.Companies); i++ {
		j := request.Companies[i]
		//add companyId(clientId) to array list for querying list of clientIds
		list_of_clientIds = append(list_of_clientIds, j.CompanyID)

		//build new string for each company to include Entitlements (blank initially) and EffectiveDate
		entitlement := IPCompany{ID: j.CompanyID, ReceivedDate: "2017-01-01", RecordDate: "2017-01-01", Shares: "100", Status: Status{Code: "0", Message: "initial message"}, EffectiveDate: currentDate()}
		//add to CompanyList array
		entList.AddIPEntitlement(entitlement)

		ent, err := json.Marshal(entList)
		err = stub.PutState(j.CompanyID, []byte(ent))
		if err != nil {
			return nil, err
		}

		var edList []string
		edList = append(edList, edCurrentDate())
		MEfList, _ := json.Marshal(edList)
		EfList := stub.PutState(j.CompanyID+"_Effective_Date_List", []byte(MEfList))
		if EfList != nil {
			return nil, EfList
		}

	}

	//encode new string into JSON format using marshal
	entitlements, err := json.Marshal(entList)
	if err != nil {
		return nil, err
	}

	//set Entitlements key and Marshaled entitlements compList as value
	//use for querying list of all Entitlements
	err = stub.PutState("Entitlements", entitlements)
	if err != nil {
		return nil, err
	}

	//encode clientId list array into JSON format using marshal
	clientIDs, cltErr := json.Marshal(list_of_clientIds)
	if cltErr != nil {
		return nil, cltErr
	}

	//add ClientID as key and Entitlements as value
	//used for querying list of all ClientIds
	err = stub.PutState("ClientID", []byte(clientIDs))
	if err != nil {
		return nil, err
	}

	//add Event Chaincode ID
	//used for querying the Chaincode ID for the Event associated with Entitlement
	err = stub.PutState("Event", []byte(request.CCID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *IPChaincode) updateInterestedParties(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var jsonResp string

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1. The api entitlement Json")
	}

	//unmarshal json from api request
	var request APIRequest
	//map entjson to model
	err2 := json.Unmarshal([]byte(args[0]), &request)

	if err2 != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal request\"}"
		return nil, errors.New(jsonResp)
	}

	//Build string in DLT format for comparison
	entArray := []IPCompany{}
	entList := IPEntitlements{entArray}

	for i := 0; i < len(request.Entitlement); i++ {
		j := request.Entitlement[i]

		//validate string field ParticipantID
		validateID := validate(j.ID, "string", "ID", 6)
		if validateID != "true" {
			jsonResp = "{\"Error\":\"" + validateID + "\"}"
			return nil, errors.New(jsonResp)
		}

		//validate date field received date
		validateReceivedDate := validate(j.ReceivedDate, "date", "ReceivedDate", 10)
		if validateReceivedDate != "true" {
			jsonResp = "{\"Error\":\"" + validateReceivedDate + "\"}"
			return nil, errors.New(jsonResp)
		}

		//validate date field record date
		validateRecordDate := validate(j.RecordDate, "date", "RecordDate", 10)
		if validateRecordDate != "true" {
			jsonResp = "{\"Error\":\"" + validateRecordDate + "\"}"
			return nil, errors.New(jsonResp)
		}

		ipCompany := IPCompany{ID: j.ID, ReceivedDate: j.ReceivedDate, RecordDate: j.RecordDate, Shares: j.Shares, Status: Status{Code: j.Status.Code, Message: j.Status.Message}}
		entList.AddEntitlement(ipCompany)
	}

	//marshal ta entitlement
	mrshlDepEnt, mrshlIPErr := json.Marshal(entList)
	if mrshlDepErr != nil {
		jsonResp = "{\"Error\":\"Failed to marshal to IPEntitlement\"}"
		return nil, errors.New(jsonResp)
	}

	//Get current entitlement
	entBytes, err := stub.GetState(request.Entitlement[0].ID)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for ID\"}"
		return nil, errors.New(jsonResp)
	}

	//Remove '\' from State to enable string comparison
	str := []byte(entBytes)
	entEscaped := bytes.Replace(str, []byte("\\"), []byte(""), -1)

	//unmarshal json from fabric
	var fabEntitlement Entitlements
	//map entjson to model
	err1 := json.Unmarshal([]byte(entEscaped), &fabEntitlement)

	if err1 != nil {
		jsonResp = "{\"Error\":\"Failed to unmarshal to fabEntitlement\"}"
		return nil, errors.New(jsonResp)
	}

	fabDepEnt, fabDepEntErr := json.Marshal(fabEntitlement)
	if fabDepEntErr != nil {
		jsonResp = "{\"Error\":\"Failed to marshal to fabEntitlement\"}"
		return nil, errors.New(jsonResp)
	}

	//compare Entitlements
	if string(mrshlDepEnt) != string(fabDepEnt) {
		fmt.Println("IPEntitlements do not match. Update!")

		entArray := []IPCompany{}
		entList := IPEntitlements{entArray}

		for i := 0; i < len(request.Entitlement); i++ {
			j := request.Entitlement[i]

			ipCompanyUpdate := IPCompany{ID: j.ID, ReceivedDate: j.ReceivedDate, RecordDate: j.RecordDate, Shares: j.Shares, Status: Status{Code: j.Status.Code, Message: j.Status.Message}, EffectiveDate: currentDate()}
			entList.AddEntitlement(ipCompanyUpdate)

		}

		//marshal updated entitlement back up
		mrshlEnt, mrshlCompErr := json.Marshal(entList)
		if mrshlCompErr != nil {
			return nil, mrshlCompErr
		}

		newEnt := stub.PutState(request.Entitlement[0].ID, []byte(mrshlEnt))
		if newEnt != nil {
			return nil, newEnt
		}

		//add entitlement key with date for querying by date
		newEDEnt := stub.PutState(request.Entitlement[0].ID+"_"+edCurrentDate(), []byte(mrshlEnt))
		if newEDEnt != nil {
			return nil, newEDEnt
		}

		//add list of dates that the entitlement was updated---------
		var edList []string
		//get the effective date list
		GSEfList, gsErr := stub.GetState(request.Entitlement[0].ID + "_Effective_Date_List")
		if gsErr != nil {
			return nil, gsErr
		}
		//unmarshal entitlement effective date list
		umEfList := json.Unmarshal([]byte(GSEfList), &edList)
		if umEfList != nil {
			return nil, umEfList
		}
		// append to the list only if the first item in the list is not equal to the
		// current date (to prevent dupes)
		if edList[0] != edCurrentDate() {
			//add the current effective date to array list
			edList = append(edList, edCurrentDate())
			//marshal lists
			EffDateList, _ := json.Marshal(edList)

			EfList := stub.PutState(request.Entitlement[0].ID+"_Effective_Date_List", []byte(EffDateList))
			if EfList != nil {
				return nil, EfList
			}
		}

	} else {
		return nil, nil
	}

	return nil, nil
}
