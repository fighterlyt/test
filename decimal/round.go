package main

import (
	"fmt"
	"github.com/fighterlyt/decimal"
)

func main(){
	fmt.Println(decimal.NewFromFloat(5.55).RoundBank(1).String())
	fmt.Println(decimal.NewFromFloat(5.54).RoundBank(1).String())
	fmt.Println(decimal.NewFromFloat(5.56).RoundBank(1).String())
	fmt.Println(decimal.NewFromFloat(5.555).RoundBank(1).String())
	fmt.Println(decimal.NewFromFloat(5.554).RoundBank(1).String())

	fmt.Println(decimal.NewFromFloat(0.27462770316074864).RoundBank(8).String())
}