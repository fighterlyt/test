package main

import (
	"orderbook/common"
	"github.com/fighterlyt/decimal"
	"github.com/globalsign/mgo/bson"
)

//定义一个复杂的监控数据，由多个值组成

type ChartData struct {
	Time     common.Time     //时间段起点
	Open     decimal.Decimal //开市价
	Close    decimal.Decimal //闭市价
	OpenTime common.Time     //开市单对应的时间

	High    decimal.Decimal //最高价
	Low     decimal.Decimal //最低价
	Volume  decimal.Decimal //交易量
	DealIds []bson.ObjectId //交易结果的id
	openSet bool
	saved   bool
}

