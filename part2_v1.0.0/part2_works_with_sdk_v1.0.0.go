/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var tradeIndexStr = "_tradeindex"				//name for the key/value that will store a list of all known trades

type Trade struct {
	TradeDate string `json:"tradedate"`
	ValueDate string `json:"valuedate"`
	Operation string `json:"operation"`
	Quantity int `json:"quantity"`
	Security string `json:"security"`
	Price string `json:"price"`
	Counterparty string `json:"counterparty"`
	User string `json:"user"`
	Timestamp int64 `json:"timestamp"`			// utc timestamp of creation
	Settled int `json:"settled"`				// enriched & settled
	NeedsRevision int `json:"needsrevision"`	// returned to client for revision
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}
	
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(tradeIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invokcations
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													// initialize the chaincode state, used as reset
		return t.init(stub, args)
	} else if function == "write" {											// writes a value to the chaincode state
		return t.write(stub, args)
	} else if function == "create_and_submit_trade" {								// create and submit a new trade
		return t.create_and_submit_trade(stub, args)
	} 

	fmt.Println("run did not find func: " + function)						// error

	return nil, errors.New("Received unknown function invocation")

}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")

}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	var timestamp, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting timestamp of the var to query")
	}

	timestamp = args[0]
	valAsbytes, err := stub.GetState(timestamp)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + timestamp + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil													//send it onward
}

// ============================================================================================================================

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	var timestamp, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. timestamp of the variable and value to set")
	}

	timestamp = args[0]													
	value = args[1]
	err = stub.PutState(timestamp, []byte(value))								//write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil

}

// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) create_and_submit_trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	// reference
	// chaincode.invoke.init_trade([data.tradedate, data.valuedate, data.operation, data.quantity, data.security, data.price, data.counterparty, data.user, data.timestamp, data.settled, data.needsrevision], cb_invoked);				//create a new trade

	var err error

	if len(args) != 11 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}

	fmt.Println("- start create_and_submit_trade")

	// if len(args[0]) <= 0 {
	// 	return nil, errors.New("1st argument must be a non-empty string")
	// }
	// if len(args[1]) <= 0 {
	// 	return nil, errors.New("2nd argument must be a non-empty string")
	// }
	// if len(args[2]) <= 0 {
	// 	return nil, errors.New("3rd argument must be a non-empty string")
	// }
	// if len(args[3]) <= 0 {
	// 	return nil, errors.New("4th argument must be a non-empty string")
	// }
	
	tradedate := strings.ToLower(args[0])
	valuedate := strings.ToLower(args[1])
	operation := strings.ToLower(args[2])

	quantity, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("4th argument must be a numeric string")
	}

	security := strings.ToLower(args[4])
	price := strings.ToLower(args[5])
	counterparty := strings.ToLower(args[6])
	user := strings.ToLower(args[7])
	
	timestamp := makeTimestamp()
	timestampAsString := strconv.FormatInt(timestamp, 10)

	settled, err := strconv.Atoi(args[9])
	if err != nil {
		return nil, errors.New("9th argument must be a numeric string, either 0 or 1")
	}

	needsrevision, err := strconv.Atoi(args[10])
	if err != nil {
		return nil, errors.New("10th argument must be a numeric string, either 0 or 1")
	}

	// reference
	// chaincode.invoke.init_trade([data.tradedate, data.valuedate, data.operation, data.quantity, data.security, data.price, data.counterparty, data.user, data.timestamp, data.settled, data.needsrevision], cb_invoked);				//create a new trade

	str := `{"tradedate": "` + tradedate + `", "valuedate": "` + valuedate + `", "operation": "` + operation + `", "quantity": ` + strconv.Itoa(quantity) + `, "security": "` + security + `", "price": "` + price + `", "counterparty": "` + counterparty + `", "user": "` + user + `", "timestamp": "` + timestampAsString + `", "settled": "` + strconv.Itoa(settled) + `", "needsrevision": "` + strconv.Itoa(needsrevision) + `"}`

	err = stub.PutState(timestampAsString, []byte(str))							// store trade with timestamp as key
	if err != nil {
		return nil, err
	}
		
	tradesAsBytes, err := stub.GetState(tradeIndexStr)					//get the trade index
	if err != nil {
		return nil, errors.New("Failed to get trade index")
	}

	var tradeIndex []string
	json.Unmarshal(tradesAsBytes, &tradeIndex)							// un stringify it aka JSON.parse()
	
	//append
	tradeIndex = append(tradeIndex, timestampAsString)					// add trade timestampAsString to index list
	fmt.Println("! trade index: ", tradeIndex)
	jsonAsBytes, _ := json.Marshal(tradeIndex)
	err = stub.PutState(tradeIndexStr, jsonAsBytes)						// store name of trade

	fmt.Println("- end create_and_submit_trade")
	return nil, nil

}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
