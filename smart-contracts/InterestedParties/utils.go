package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

//adds company string to array
func (entitlement *IPEntitlements) AddIPEntitlement(item IPCompany) []IPCompany {
	entitlement.Entitlements = append(entitlement.Entitlements, item)
	return entitlement.Entitlements
}

func (entitlement *DLTIPEntitlements) AddIPEntitlementUpdate(item DLTIPCompany) []DLTIPCompany {
	entitlement.Entitlements = append(entitlement.Entitlements, item)
	return entitlement.Entitlements
}

//gets current date
func currentDate() string {
	currentDateTime := time.Now()
	const layout = "2006-01-02"
	return currentDateTime.Format(layout)
}

func edCurrentDate() string {
	//get current datetime
	currentDateTime := time.Now()
	const layout = "20060102"
	return currentDateTime.Format(layout)
}

func validate(value string, dataType string, fieldName string, length int) string {
	var response string
	dataLength := len(value)
	//fields that need to be a certain length
	if fieldName == "CCID" {
		if dataLength != length {
			response = "Data length is incorrect for " + fieldName
			return response
		}
	}
	if dataLength > length {
		response = "Data length is incorrect for " + fieldName
		return response
	} else if dataType == "string" {
		response = "true"
		return response
	}

	//validate decimal data type
	if dataType == "decimal" {
		//check if value is a decimal by converting to an int
		flt, fltErr := strconv.Atoi(value)
		//if the string converts to an int then its not a decimal
		if fltErr != nil {
			//check to see if decimal has 5 leading numbers after decimal point
			decStr := strings.Split(value, ".")
			lzStr := len(decStr[1])
			if lzStr != 5 {
				response = "Invalid format for " + fieldName
				return response
			}
			//value is not an int, convert to decimal
			data, decErr := decimal.NewFromString(value)
			if decErr != nil {
				response = "Invalid format for " + fieldName
				return response
			} else {
				response = "true"
				//make sure decimal has 5 trailing zeros
				fmt.Println("data value:", data.StringFixed(5))
				return response
			}
		} else {
			response = "Data is not decimal format for " + fieldName
			fmt.Println("data value:", flt)
			return response
		}

	}
	if dataType == "int" {
		//convert string to int
		voteableInt, vIntErr := strconv.Atoi(value)
		if vIntErr != nil {
			response = "Failed to convert string to int, incorrect format for " + fieldName
			return response
		} else {
			response = "true"
			fmt.Println("valid voteable int:", voteableInt)
			return response
		}
	}
	//validate date data type
	if dataType == "date" {
		const layout = "2006-01-02"
		d, dErr := time.Parse(layout, value)
		if dErr != nil {
			response = "Date in incorrect format for " + fieldName
			return response
		} else {
			fmt.Println("valid date:", d.Format(layout))
			response = "true"
			return response
		}
	}

	//validate date data type
	if dataType == "queryDate" {
		const layout = "20060102"
		qd, qdErr := time.Parse(layout, value)
		if qdErr != nil {
			response = "Date in incorrect format for " + fieldName
			return response
		} else {
			fmt.Println("valid date:", qd.Format(layout))
			response = "true"
			return response
		}
	}

	return response
}
