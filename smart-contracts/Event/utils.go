package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

func currentDate() string {
	//get current datetime
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

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func sTask() {
	fmt.Println("Updating event...")
}

func validate(value string, dataType string, fieldName string, length int) string {
	var response string
	dataLength := len(value)
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
