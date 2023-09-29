package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type proposal struct {
	Code []byte
}

func main() {
	code, err := ioutil.ReadFile("./contracts/new_ics10_grandpa_cw.opt.wasm")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(code)
	proposal := proposal{Code: []byte(code)}

	data, err := json.MarshalIndent(proposal, "", "  ") // MarshalIndent for pretty printing
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}

	// Save JSON to a file
	err = ioutil.WriteFile("output.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}
