package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	// Step 1: Read the JSON file
	jsonData, err := ioutil.ReadFile("./contracts/contract.data")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	var wasmBinary []byte
	// Step 2: Unmarshal the JSON content
	err = json.Unmarshal(jsonData, &wasmBinary)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Step 3: Write the byte array back to a .wasm file
	err = ioutil.WriteFile("reconstructed.wasm", wasmBinary, 0644)
	if err != nil {
		fmt.Println("Error writing to WASM file:", err)
		return
	}

	fmt.Println("WASM file reconstructed successfully!")
}