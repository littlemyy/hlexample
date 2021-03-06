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


/*
peer chaincode invoke -o orderer.lyl-network.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/lyl-network.com/orderers/orderer.lyl-network.com/msp/tlscacerts/tlsca.lyl-network.com-cert.pem  -C $CHANNEL_NAME -n sacc -c '{"Args":["makeApplication","{\"applicationId\": \"LEP0000001\",\"seller\": {\"firstName\": \"Riita\",\"lastName\": \"Ratas\",\"personalCode\": \"123456789\"},\"buyer\": {\"firstName\": \"Mari\",\"lastName\": \"Maasikas\",\"personalCode\": \"123456779\"},\"vehicle\": {\"vin\":\"78347837483784\",\"mark\":\"Audi\",\"model\":\"A8\",\"registrationPlate\":\"123ABC\"},\"price\":\"30000.00\",\"status\":\"\"}"]}'

peer chaincode invoke -o orderer.lyl-network.com:7050  --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/lyl-network.com/orderers/orderer.lyl-network.com/msp/tlscacerts/tlsca.lyl-network.com-cert.pem  -C $CHANNEL_NAME -n sacc -c '{"Args":["readApplication","{\"applicationId\": \"LEP0000001\"}"]}'


*/

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	//"bytes"
	"encoding/json"
	"fmt"
	//"strconv"
  "errors"
	"strings"
	//"reflect"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
	//"github.com/icrowley/fake"
)

// Define the Application Contract structure
type ApplicationContract struct {
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
	Price  *string `json:"price,omitempty"`
	Status *string `json:"status,omitempty"`
}

const ACCEPTED string ="accepted"
const REJECTED string ="rejected"
const CANCELLED string ="cancelled"
const WAITING string="waiting"
const FINISHED string="finished"

/*
 * The Init method is called when the Application Contract "sale application" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (t *ApplicationContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	fmt.Println("Init")
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "sale application"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (t *ApplicationContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "makeTestData" {
		return t.makeTestData(APIstub)
	} else if function == "makeApplication" {
		return t.makeApplication(APIstub, args)
	} else if function == "acceptApplication" {
		return t.acceptApplication(APIstub, args)
	} else if function == "rejectApplication" {
		return t.rejectApplication(APIstub, args)
	} else if function == "cancelApplication" {
		return t.cancelApplication(APIstub, args)
	} else if function == "getBuyerApplications" { //list of sale applications
		return t.getBuyerApplications(APIstub, args)
	} else if function =="getSellerApplications" {
		return t.getSellerApplications(APIstub, args)
	} else if function =="getInApplications" {
		return t.getInApplications(APIstub, args)
	} else if function =="getOutApplications" {
		return t.getOutApplications(APIstub, args)
	} else if function =="readApplication" {
		return t.readApplication(APIstub, args)
	}

	fmt.Println("query did not find func: " + function)
  return shim.Error("Received unknown function query: " + function)
}

/*func (s *ApplicationContract) makeTestData2(APIstub shim.ChaincodeStubInterface) sc.Response {
	var seller_first_name string
	var seller_last_name string
	var seller_peronal_code string
  var seller Person
	var buyer Person
	var buyer_first_name string
	var buyer_last_name string
	var buyer_personal_code string
	var vehicle Vehicle
  var vehicle_vin string
	var vehicle_mark string
  var vehicle_model string
	var vehicle_registration_plate string
	var price string
	var status string

	applicationsIn := []SaleApplication

	j := 0

	for j < 50 {
		seller_first_name = fake.FirstName()
		seller_last_name = fake.lastName()
		seller_peronal_code = "123456791011"
		seller = Person{FirstName:&seller_first_name,LastName:&seller_last_name,PersonalCode:&seller_personal_code}

		buyer_first_name = fake.FirstName()
		buyer_last_name = fake.lastName()
		buyer_personal_code ="123456791011"
		buyer = Person{FirstName:&buyer_first_name,LastName:&buyer_last_name,PersonalCode:&buyer_personal_code}

		vehicle_vin ="12345678"
		vehicle_mark =fake.Brand()
		vehicle_model =fake.Model()
		vehicle_registration_plate="123ABS"
		price ="100000.00"
		status=WAITING

		vehicle = Vehicle{Vin:&vehicle_vin,Mark:&vehicle_mark,Model:&vehicle_model,RegistrationPlate:&vehicle_registration_plate}

		applicationsIn[i]:= SaleApplication{ApplicationId:&applicationId, Seller:&seller, Buyer:&buyer, Vehicle:&vehicle, Price:&price, Status:&status},

	}

	i := 0
	for i < len(applicationsIn) {
		fmt.Println("i is ", i)
		applicationAsBytes, _ := json.Marshal(applicationsIn[i])
		APIstub.PutState(*applicationsIn[i].ApplicationId, applicationAsBytes)
		fmt.Println("Added", applicationsIn[i])
		i = i + 1
	}

	return shim.Success(nil)
}*/

func (s *ApplicationContract) makeTestData(APIstub shim.ChaincodeStubInterface) sc.Response {
	var applicationId string = "100000"
	var seller_first_name string ="Ulvi"
	var seller_last_name string ="Sädem"
  var seller_personal_code string="49104231234"
  var seller Person

	seller = Person{FirstName:&seller_first_name,LastName:&seller_last_name,PersonalCode:&seller_personal_code}
	var buyer Person
	var buyer_first_name string="Pilvi"
	var buyer_last_name string="Sädem"
  var buyer_personal_code string="47712121234"
	buyer = Person{FirstName:&buyer_first_name,LastName:&buyer_last_name,PersonalCode:&buyer_personal_code}

	var vehicle Vehicle
  var vehicle_vin string="12345678"
	var vehicle_mark string="audi"
  var vehicle_model string="a8"
	var vehicle_registration_plate string="123ABS"
	var price string="100000.00"
	var status string=WAITING

  vehicle = Vehicle{Vin:&vehicle_vin,Mark:&vehicle_mark,Model:&vehicle_model,RegistrationPlate:&vehicle_registration_plate}


	//applicationId = "100000"

	applicationsIn := []SaleApplication{
		SaleApplication{ApplicationId:&applicationId, Seller:&seller, Buyer:&buyer, Vehicle:&vehicle, Price:&price, Status:&status},
		}

	i := 0
	for i < len(applicationsIn) {
		fmt.Println("i is ", i)
		applicationAsBytes, _ := json.Marshal(applicationsIn[i])
		APIstub.PutState(*applicationsIn[i].ApplicationId, applicationAsBytes)
		fmt.Println("Added", applicationsIn[i])
		i = i + 1
	}

	return shim.Success(nil)
}


// Function is called to validate input
func (t *ApplicationContract) validateInput(args []string) (applicationIn SaleApplication, err error) {
	var applicationId string //application Id
	var saleApplication SaleApplication = SaleApplication{} //The calling function is expecting an object of type SaleApplication
  fmt.Println("validateInput")
	// Assumes that we have json input
	if len(args) !=1 {
		err = errors.New("Incorrect number of arguments. Expecting a json string with mandatory applicationId")
		return saleApplication, err

	}

	jsonData:=args[0]
	applicationId=""
	stateJSON :=[]byte(jsonData)
	//fmt.Println("State json"+stateJSON)
	err = json.Unmarshal(stateJSON,&applicationIn)

	if err!=nil {
		err = errors.New("Unable to unmarshal input JSON data"+fmt.Sprint(err))
		return saleApplication, err
		}

	if applicationIn.ApplicationId !=nil {
		applicationId = strings.TrimSpace(*applicationIn.ApplicationId)
		if applicationId=="" {
			err = errors.New("ApplicationId not passed")
			return saleApplication, err
		}
	} else {
		err = errors.New("Application ID is mandatory in the input JSON data")
		return saleApplication, err
	}
	applicationIn.ApplicationId = &applicationId
	return applicationIn, nil
}

// Function is called to read asset information
func (t *ApplicationContract) readApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var applicationID string
	var err error
	var saleApplication SaleApplication

	fmt.Println("running readApplication()")


	applicationIn,err:=t.validateInput(args)
	if err!=nil {
		return shim.Error("Couldn't find the application")
	}

	applicationID = *applicationIn.ApplicationId
	applicationAsBytes, err := APIstub.GetState(applicationID)
	if err!=nil || len(applicationAsBytes)==0{
		return shim.Error("Unable to get application state from the ledger")
	}

	err = json.Unmarshal(applicationAsBytes, &saleApplication)
	if err!=nil {
		return shim.Error("Unable to unmarshal application data received from the ledger")
	}
	fmt.Println("Application status " + *saleApplication.Status)
	return shim.Success(applicationAsBytes)

}


// Function is called in order to make a new application
func (t *ApplicationContract) makeApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	var err error
	var applicationId string
	//var applicationIn SaleApplication
  var applicationStub SaleApplication
	//var applicationBytes []byte

	fmt.Println("running makeApplication()")

	applicationIn,err:=t.validateInput(args)
	if err!=nil {
		return shim.Error(fmt.Sprint(err))
	}
	/*if len(args) !=1 {
		err = errors.New("Incorrect number of arguments. Expecting a json string with mandatory applicationId")
		return shim.Error(fmt.Sprint(err))

	}

	jsonData:=args[0]
	stateJSON :=[]byte(jsonData)
	//fmt.Println("State json"+stateJSON)
	err = json.Unmarshal(stateJSON,&applicationIn)
	*/

  applicationId = *applicationIn.ApplicationId
	applicationAsBytes, err := APIstub.GetState(applicationId)
	if err != nil || len(applicationAsBytes)==0 {
		applicationStub = applicationIn //The record that goes into stub is the one that came in
	} else {
		fmt.Println("else case: update")
	}

	/* Possible business rules
		- Vehicle must be provided
		- Vehicle must be registered in Vehicle Ledger
		- Seller has to have rights to initiate the sale
		- Sama auto kohta ei tohi olla teist taotlust
	*/

  *applicationStub.Status = WAITING
	applicationJSON, err := json.Marshal(applicationStub)
	if err != nil {
		return shim.Error("Marshal failed for contract state" + fmt.Sprint(err))
	}
	//get existing state from the stub

	//Write the new state to the ledger
	err = APIstub.PutState(applicationId, applicationJSON)
	if err != nil {
		return shim.Error("Put ledger state failed: "+ fmt.Sprint(err))
	}


	return shim.Success(nil)
	/*if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}*/

	/*assetAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(assetAsBytes)*/
}


// function is called to change application status
func (t *ApplicationContract) changeApplicationStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("running changeApplicationStatus: " + args[1])
	/* Possible business rules
	*/

	



	return shim.Success(nil)
}

// function is called to change application status to Accepted
func (t *ApplicationContract) acceptApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
  fmt.Println("running acceptApplication()")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
  args[1]=ACCEPTED
	t.changeApplicationStatus(APIstub,args)
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
func (t *ApplicationContract) rejectApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("running rejectApplication()")
	args[1]=REJECTED
	t.changeApplicationStatus(APIstub,args)

	/* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/

	return shim.Success(nil)
}

// function is called to change application status to Cancelled
func (t *ApplicationContract) cancelApplication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	fmt.Println("running cancelApplication()")
  args[1]=CANCELLED
	t.changeApplicationStatus(APIstub, args)

	/* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/

	return shim.Success(nil)
}

/* function returns applications made for concrete buyer */
func (t *ApplicationContract) getBuyerApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getBuyerApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	// Need to query all applications by Buyer

	return shim.Success(nil)

}

/* function returns applications made by concrete seller */
func (t *ApplicationContract) getSellerApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getSellerApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	// Need to query all applications by Seller


	return shim.Success(nil)

}

/* function returns all incoming applications */
func (t *ApplicationContract) getInApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	fmt.Println("running getInApplications()")
  /* if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}*/
	//Need to query all applications by Seller Leasing

	return shim.Success(nil)

}
/* function returns all outgoing applications */
func (t *ApplicationContract) getOutApplications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

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
	err := shim.Start(new(ApplicationContract))
	if err != nil {
		fmt.Printf("Error creating new Application Contract: %s", err)
	}
}
