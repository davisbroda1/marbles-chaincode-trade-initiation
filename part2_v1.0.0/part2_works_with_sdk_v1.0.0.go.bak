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
	Quantity int `json:"quantity,string"`
	Security string `json:"security"`
	Price string `json:"price"`
	Counterparty string `json:"counterparty"`
	User string `json:"user"`
	Timestamp string `json:"timestamp"`			// utc timestamp of creation, use JS/jQuery timestamp as string
	Settled int `json:"settled,string"`			// enriched & settled
	NeedsRevision int `json:"needsrevision,string"`	// returned to client for revision
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
	} else if function == "mark_revision_needed" {								
		return t.mark_revision_needed(stub, args)
	} else if function == "mark_revised" {								
		return t.mark_revised(stub, args)
	} else if function == "enrich_and_settle" {								
		return t.enrich_and_settle(stub, args)
	}

	fmt.Println("run did not find func: " + function)						// error

	return nil, errors.New("Received unknown function invocation")

}

// ============================================================================================================================
// Invoke - Our entry point for Invokcations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													// initialize the chaincode state, used as reset
		return t.init(stub, args)
	} else if function == "write" {											// writes a value to the chaincode state
		return t.write(stub, args)
	} else if function == "create_and_submit_trade" {								// create and submit a new trade
		return t.create_and_submit_trade(stub, args)
	} else if function == "mark_revision_needed" {								
		return t.mark_revision_needed(stub, args)
	} else if function == "mark_revised" {								
		return t.mark_revised(stub, args)
	} else if function == "enrich_and_settle" {								
		return t.enrich_and_settle(stub, args)
	}

	fmt.Println("run did not find func: " + function)						// error

	return nil, errors.New("Received unknown function invocation")

}

// ============================================================================================================================
// Init - Our entry point for Invokcations
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	
	fmt.Println("run is running " + function)

	// Handle different functions
	if function == "init" {													// initialize the chaincode state, used as reset
		return t.init(stub, args)
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
	
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting key of the value to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)									//get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
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

// create_and_submit_trade - create a new trade, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) create_and_submit_trade(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
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
	
	// use jquery timestamp string for now with time zone
	timestamp := strings.ToLower(args[8])

	// timestamp := makeTimestamp()
	// timestampAsString := strconv.FormatInt(timestamp, 10)

	settled, err := strconv.Atoi(args[9])
	if err != nil {
		return nil, errors.New("9th argument must be a numeric string, either 0 or 1")
	}

	needsrevision, err := strconv.Atoi(args[10])
	if err != nil {
		return nil, errors.New("10th argument must be a numeric string, either 0 or 1")
	}

	str := `{"tradedate": "` + tradedate + `", "valuedate": "` + valuedate + `", "operation": "` + operation + `", "quantity": ` + strconv.Itoa(quantity) + `, "security": "` + security + `", "price": "` + price + `", "counterparty": "` + counterparty + `", "user": "` + user + `", "timestamp": "` + timestamp + `", "settled": "` + strconv.Itoa(settled) + `", "needsrevision": "` + strconv.Itoa(needsrevision) + `"}`

	fmt.Println("str: ", str)

	err = stub.PutState(timestamp, []byte(str))							// store trade with timestamp as key
	if err != nil {
		return nil, err
	}

	fmt.Println("put state for timestamp key: ", timestamp)
		
	tradesAsBytes, err := stub.GetState(tradeIndexStr)					//get the trade index
	if err != nil {
		return nil, errors.New("Failed to get value for key _tradeIndex")
	}

	var tradeIndex []string
	json.Unmarshal(tradesAsBytes, &tradeIndex)							// un stringify it aka JSON.parse()

	fmt.Println("stored timestamp at end of tradeIndexStr array: ", tradeIndex)

	//append
	tradeIndex = append(tradeIndex, timestamp)					// add trade timestamp to index list
	fmt.Println("! trade index: ", tradeIndex)
	jsonAsBytes, _ := json.Marshal(tradeIndex)
	err = stub.PutState(tradeIndexStr, jsonAsBytes)						// store name of trade

	fmt.Println("- end create_and_submit_trade")
	return nil, nil

}


// mark_revision_needed - Mark a Trade in need of revision
// ============================================================================================================================
func (t *SimpleChaincode) mark_revision_needed(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}

	fmt.Println("- start mark_revision_needed")

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	
	timestamp := strings.ToLower(args[0])
	newUser := strings.ToLower(args[1])

	tradeAsBytes, err := stub.GetState(timestamp)					// get the trade
	if err != nil {
		return nil, errors.New("Failed to get value for key timestamp")
	}

	var trade Trade
	json.Unmarshal(tradeAsBytes, &trade)

	trade.User = newUser
	trade.NeedsRevision = 1

	bytes, _ := json.Marshal(trade)

	err = stub.PutState(timestamp, bytes)							// store trade with timestamp as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end mark_revision_needed")

	return nil, nil

}

// mark_revised - Mark a Trade in as revised
// ============================================================================================================================
func (t *SimpleChaincode) mark_revised(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 11")
	}

	fmt.Println("- start mark_revised")

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	
	timestamp := strings.ToLower(args[0])
	newUser := strings.ToLower(args[1])

	tradeAsBytes, err := stub.GetState(timestamp)					// get the trade
	if err != nil {
		return nil, errors.New("Failed to get value for key timestamp")
	}

	var trade Trade
	json.Unmarshal(tradeAsBytes, &trade)

	trade.User = newUser
	trade.NeedsRevision = 0

	bytes, _ := json.Marshal(trade)

	err = stub.PutState(timestamp, bytes)							// store trade with timestamp as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end mark_revised")

	return nil, nil

}

// enrich_and_settle - Enrich a Trade and mark it as settled
// ============================================================================================================================
func (t *SimpleChaincode) enrich_and_settle(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start enrich_and_settle")

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	
	timestamp := strings.ToLower(args[0])
	newUser := strings.ToLower(args[1])

	tradeAsBytes, err := stub.GetState(timestamp)					// get the trade
	if err != nil {
		return nil, errors.New("Failed to get value for key timestamp")
	}

	var trade Trade
	json.Unmarshal(tradeAsBytes, &trade)

	trade.User = newUser
	trade.NeedsRevision = 0
	trade.Settled = 1

	bytes, _ := json.Marshal(trade)

	err = stub.PutState(timestamp, bytes)							// store trade with timestamp as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end enrich_and_settle")

	return nil, nil

}

// ============================================================================================================================
// Make Timestamp - create a timestamp in ms
// ============================================================================================================================
func makeTimestamp() int64 {
    return time.Now().UnixNano() / (int64(time.Millisecond)/int64(time.Nanosecond))
}
