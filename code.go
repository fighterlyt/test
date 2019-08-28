package main

import (
	"encoding/json"
	"os"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	file, err := os.Open("/Users/mac/bankCode.json")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	data := &Data{}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		panic(err.Error())
	}
	bankNames := make(map[string]struct{}, 100)
	for _, record := range data.U {
		spew.Dump(record)
		bankNames[record["bankName"].(string)] = struct{}{}
	}
	target, err := os.Create("/Users/mac/bank.json")
	if err != nil {
		panic(err.Error())
	}
	defer target.Close()
	banks := make([]string, 0, len(bankNames))
	for name := range bankNames {
		banks = append(banks, name)
	}
	json.NewEncoder(target).Encode(banks)
}

type Data struct {
	U []map[string]interface{} `json:"allBankList"`
}
