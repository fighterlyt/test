package main

import (
	"orderbook/matchers"
	"os"
	"github.com/json-iterator/go"
	"time"
	"orderbook/common"
	"github.com/globalsign/mgo/bson"
	"orderbook/mysql"
	"orderbook/dealers"
	"github.com/fighterlyt/decimal"
	"fmt"
)

func main() {
	results := make([]matchers.MatchedResult, 0, 100000)
	if file, err := os.Open("./bench.json"); err != nil {
		panic(err.Error())
	} else {
		dao, _ := newDao()
		jsoniter.NewDecoder(file).Decode(&results)

		dao.Clear(results[0].CounterParts[0].Order)
		for _, result := range results {

			for i, counterPart := range result.CounterParts {
				if time.Time(counterPart.Order.CreateTime).IsZero() {
					counterPart.Order.CreateTime = common.TimeN(time.Now())
				}
				counterPart.Order.OrderId = bson.NewObjectId().Hex()
				dao.SaveOrder(counterPart.Order)
				result.CounterParts[i].Order.SetAvailableQuantity("0")

			}
		}

		start := time.Now()
		for i, result := range results {
			dealed := calculate(&result)
			dao.Deal(dealed)
			fmt.Printf("%d 花费 %s\n",i+1,time.Since(start).String())
		}
		println(time.Since(start).String())

	}
}

func newDao() (*mysql.MysqlDao, error) {
	if dao, err := mysql.NewMysqlDao("root", "123", "localhost", "3306", "test"); err != nil {
		return nil, err
	} else {
		if err = dao.Init([]string{"test", "ETH-BTC"}); err != nil {
			return nil, err
		}
		return dao, nil
	}
}
func calculate(matchResult *matchers.MatchedResult) *dealers.DealResult {
	result := &dealers.DealResult{
		Id: bson.NewObjectId(),

		Index:        matchResult.Index,
		CounterParts: make([]*dealers.DealElement, 0, len(matchResult.CounterParts)),
	}
	value := decimal.Decimal{}

	price := matchResult.Order.Order.Price()
	if price.Equal(decimal.Zero) { //价格为零，市价单，总价值按照限价单计算
		for _, element := range matchResult.CounterParts {
			elementValue := element.Order.Price().Mul(element.Quantity)
			value = value.Add(elementValue)
			//d.logger.Log(log.Debug, "市价单交易", value.String(), elementValue.String(), element.Order.Price().String())

			result.CounterParts = append(result.CounterParts, &dealers.DealElement{
				Order:    element.Order,
				Quantity: element.Quantity,
				Value:    elementValue,
			})
		}
	} else {
		//限价单，总价值按照限价单计算
		for _, element := range matchResult.CounterParts {

			elementValue := price.Mul(element.Quantity)
			value = value.Add(elementValue)
			result.CounterParts = append(result.CounterParts, &dealers.DealElement{
				Order:    element.Order,
				Quantity: element.Quantity,
				Value:    elementValue,
			})

		}
	}
	result.Provider = &dealers.DealElement{
		Order:    matchResult.Order.Order,
		Quantity: matchResult.Order.Quantity,
		Value:    value,
	}
	return result

}
