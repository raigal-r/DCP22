package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type RaquelChaincode struct {
}

// Asset struct and properties must be exported (start with capitals) to work with contract api metadata
//type Asset struct {
//	ObjectType        string `json:"objectType"` // ObjectType is used to distinguish different object types in the same chaincode namespace
//	ID                string `json:"assetID"`
//	OwnerOrg          string `json:"ownerOrg"`
//	PublicDescription string `json:"publicDescription"`
//}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *RaquelChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Init method gets called")
	//_, args := stub.GetFunctionAndParameters()
	var args [3]string
	args[0] = "a"          // Assign a value to the first element
	args[1] = "1a2b3c4d5e" // Assign a value to the second element The NFT
	args[2] = "b"          // Assign a value to the third element
	//args1 := [3]string{"EntityA","1a2b3c4d5e","EntityB"}
	var A, B string       // Entities
	var Aval, Bval string // Asset holdings
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
	A = args[0]
	Aval = strconv.Quote(args[1])
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}
	//B = args[2]

	//Bval, err = strconv.Atoi(args[3])
	Bval = strconv.Quote("")
	//if err != nil {
	//	return shim.Error("Expecting integer value for asset holding")
	//}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger Thats the mapping
	err = stub.PutState(A, []byte(strconv.QuoteToASCII(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.QuoteToASCII(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (t *RaquelChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("Invoke method gets called")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
		//} else if function == "delete" { //not necessary
		//	// Deletes an entity from its state
		//	return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke // When you call query, it give you the entity holdings
		fmt.Printf("Query Response:%s\n", t.query(stub, args))
		return t.query(stub, args)
	}
	//else if function == "NFThold" {
	//	// the old "Query" is now implemtned in invoke
	//	return t.holds(stub, args)
	//}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Transaction makes payment of X units from A to B
func (t *RaquelChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("invoke method gets called")
	var A, B string       // Entities
	var Aval, Bval string // Asset holdings
	//var X string          // Transaction value XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = (string(Avalbytes))

	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	Bval, _ = (string(Bvalbytes))

	// Perform the execution of transfer
	//X, err = strconv.Atoi(args[2])
	//if err != nil {
	//	return shim.Error("Invalid transaction amount, expecting a integer value")
	//}
	//Aval = Aval - X
	Aval = delete(string(Avalbytes), string(args[2]))
	Bval = Bval + (string(args[2]))
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = stub.PutState(A, []byte(strconv.Quote(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Quote(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Deletes an entity from state
/*func (t *RaquelChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("delete method gets called")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}*/

// query callback representing the query of a chaincode
func (t *RaquelChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("query method gets called")
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

// Deletes an entity from state
/*func (t *RaquelChaincode) queryNFTowner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("queryNFTowner method gets called")
	var Aval string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	Aval = args[0]

	// Get the state from the ledger
	Avalentitybytes, err := stub.GetQueryResult(Aval)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + Aval + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalentitybytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + Aval + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + Aval + "\",\"Amount\":\"" + string(Avalentitybytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalentitybytes)
}*/

func main() {
	err := shim.Start(new(RaquelChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
