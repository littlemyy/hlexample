/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

//
type Person struct {
	FirstName *string `json:"firstName,omitempty"`
	LastName	*string `json:"lastName,omitempty"`
	PersonalCode *string `json:"personalCode,omitempty"`
}

type Vehicle struct {
	Vin *string `json:"vin,omitempty"`
	Mark *string `json:"mark,omitempty"`
	Model *string `json:"model,omitempty"`
	RegistrationPlate *string `json:"registrationPlate,omitempty"`
}

// Define the sale appication structure.  Structure tags are used by encoding/json library
type SaleApplication struct {
	ApplicationId *string `json:"applicationId,omitempty"`
	Seller   *Person `json:"seller,omitempty"`
	Buyer  *Person `json:"buyer,omitempty"`
	Vehicle *Vehicle `json:"vehicle,omitempty"`
	Price  *float64 `json:"price,omitempty"`
	Status *string `json:"status,omitempty"`
}

const ACCEPTED string ="accepted"
const REJECTED string ="rejected"
const CANCELLED string ="cancelled"
const WAITING string="waiting"
const FINISHED string="finished"

/*
 * The Init method is called when the Smart Contract "sale application" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "sale application"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "makeApplication" {
		return s.makeApplication(APIstub, args)
	} else if function == "acceptApplication" {
		return s.acceptApplication(APIstub)
	} else if function == "rejectApplication" {
		return s.rejectApplication(APIstub, args)
	} else if function == "cancelApplication" {
		return s.cancelApplication(APIstub, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

// Query is our entry point for queries
func (t *SmartContract) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "getBuyerApplications" { //list of sale applications
		return t.getBuyerApplications(stub, args)
	} else if function =="getSellerApplications" {
		return t.getSellerApplications(stub, args)
	} else if function =="getInApplications" {
		return t.getInApplications(stub, args)
	} else if function =="getOutApplications" {
		return t.getOutApplications(stub, args)
	}

	fmt.Println("query did not find func: " + function)
  return nil, errors.New("Received unknown function query: " + function)
}

// Function is called to validate input
func (t *SmartContract) validateInput(args []string) (applicationIn SaleApplication, err error) {
	var applicationId string //application Id
	var saleApplication SaleApplication = SaleApplication{} //The calling function is expecting an object of type SaleApplication

	if len(args) !=1 {
		err = errors.New("Incorrect number of arguments. Expecting a json string with mandatory applicationId")
		return saleApplication, err

	}
	//was applicationId present
	if applicationIn.ApplicationId !=nil {
		applicationId = strings.TrimSpace(*applicationIn.applicationId)
		if applicationId =""{
			err = errors.New("ApplicationId not passed")
			return saleApplication, err
		}
	} else {
		err = errors.New("Application ID is mandatory in the input JSON data")
		return saleApplication, err
	}
	applicationIn.applicationId = &applicationId
	return applicationIn, nil
}

// Function is called in order to make a new application
func (t *SmartContract) makeApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("running makeApplication()")


	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}


	var saleApplication = SaleApplication{ApplicationId: args[0], Seller: args[1], Buyer: args[2], Vehicle: args[3], Price: args[4], Status: args[5]}

	/* Possible business rules
		- Vehicle must be provided
		- Vehicle must be registered in Vehicle Ledger
		- Seller has to have rights to initiate the sale
		- Sama auto kohta ei tohi olla teist taotlust
	*/

	applicationAsBytes, _ := json.Marshal(saleApplication)
	APIstub.PutState(args[0], applicationAsBytes)

	return shim.Success(nil)
	/*if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}*/

	/*assetAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(assetAsBytes)*/
}


// function is called to change application status
func (t *SmartContract) changeApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string, status string) sc.Response {
	fmt.Println("running changeApplicationStatus: " + status)
	/* Possible business rules
	*/
	return shim.Success(nil)
}

// function is called to change application status to Accepted
func (t *SmartContract) acceptApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
  fmt.Println("running acceptApplication()")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	t.changeApplicationStatus(args, ACCEPTED)
	/* Possible business rules
		- ApplicationID must be provided
		- New status is within allowed statuses

	*/
	/* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/

	return shim.Success(nil)
}

// function is called to change application status to Rejected
func (t *SmartContract) rejectApplication(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("running rejectApplication()")
	applicationIn, err = t.changeApplicationStatus(args, REJECTED)

	/* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/

	return shim.Success(nil)
}

// function is called to change application status to Cancelled
func (t *SmartContract) cancelApplication(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("running cancelApplication()")
	applicationIn, err = t.changeApplicationStatus(args, CANCELLED)

	/* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/

	return shim.Success(nil)
}

/* function returns applications made for concrete buyer */
func (t *SmartContract) getBuyerApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getBuyerApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	// Need to query all applications by Buyer

	return shim.Success(nil)

}

/* function returns applications made by concrete seller */
func (t *SmartContract) getSellerApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getSellerApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	// Need to query all applications by Seller


	return shim.Success(nil)

}

/* function returns all incoming applications */
func (t *SmartContract) getInApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getInApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	//Need to query all applications by Seller Leasing

	return shim.Success(nil)

}
/* function returns all outgoing applications */
func (t *SmartContract) getOutApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getOutApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	//Need to query all applications by Buyer Leasing

	return shim.Success(nil)

}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
