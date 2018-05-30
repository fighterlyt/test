package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"strings"
	"github.com/fighterlyt/decimal"
	"fmt"
	"orderbook/common"
)

type Order struct {
	OrderId                string          `gorm:"type:varchar(100);unique_index"`
	ProductKey             string          `json:"productKey"`        //产品id
	QuantityValue          decimal.Decimal `json:"quantity"`          //数量
	AvailableQuantityValue decimal.Decimal `json:"availableQuantity"` //可匹配数量
	PriceValue             decimal.Decimal `json:"price"`             //价格
	OrderType              OrderType       `json:"orderType"`         //委托类型
	CreateTime             common.Time
}

//OrderType 委托类型
type OrderType int

const (
	SELLLIMIT        OrderType = iota //限价卖单
	SELLMARKET                        //市价卖单
	BUYLIMIT                          //限价买单
	BUYMARKET                         //市价买单
	InvalidOrderType                  //非法
)

var (
	orderTypeStringMap = map[OrderType]string{
		SELLLIMIT:  "限价卖单",
		SELLMARKET: "市价卖单",
		BUYLIMIT:   "限价买单",
		BUYMARKET:  "市价买单",
	}
	OrderTypes = []OrderType{SELLLIMIT, SELLMARKET, BUYLIMIT, BUYMARKET}
)

func (o OrderType) String() string {
	return orderTypeStringMap[o]
}

func main() {
	db, err := sql.Open("mysql", "root@/test")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if _, err := db.Exec("CREATE  TABLE IF NOT EXISTS Orders (" + fileds + ")"); err != nil {
		panic(err.Error())
	}

	//save:=Order{
	//	OrderId:"3",
	//	ProductKey:"test",
	//	QuantityValue:decimal.RequireFromString("0.7"),
	//	AvailableQuantityValue:decimal.RequireFromString("0.6"),
	//	PriceValue:decimal.RequireFromString("11.1"),
	//	CreateTime:time.Now().Local(),
	//}
	//fmt.Print(save.CreateTime.Location().String(),save.CreateTime.Format("2006-01-02 15:04:05"))
	//data:=save.trans()
	//
	//
	//if _,err=db.Exec("insert  into Orders (orderId,productKey,quantity,availableQuantity,price,priceExp,orderType,createTime) values(?,?,?,?,?,?,?,?)",data.OrderId,data.ProductKey,data.Quantity,data.AvailableQuantity,data.Price,data.PriceExp,data.OrderType,data.CreateTime);err!=nil{
	//	panic(err.Error())
	//}

	if rows, err := db.Query("select * from Orders where orderId=?", 3); err != nil {
		panic(err.Error())
	} else {
		o := order{}
		for rows.Next() {
			types, _ := rows.ColumnTypes()
			for _, t := range types {
				fmt.Println(t.DatabaseTypeName())
			}
			if err := rows.Scan(&o.OrderId, &o.ProductKey, &o.Quantity, &o.AvailableQuantity, &o.Price, &o.PriceExp, &o.OrderType, &o.CreateTime); err != nil {
				panic(err.Error())
			} else {
				fmt.Println(o.Trans())
			}
		}
	}
	if tx,err:=db.Begin();err!=nil{
		panic(err.Error())
	}
}
func (o Order) String() string {
	return fmt.Sprintf("委托 orderId:%s 价格:%s 数量:%s 可用数量:%s 类型:%s 创建时间:%s", o.OrderId, o.Price().String(), o.Quantity().String(), o.AvailableQuantity().String(), o.OrderType.String(),o.CreateTime.String())
}

func (o Order) trans() *order {
	exp := 0

	order := &order{
		ProductKey:        o.ProductKey,
		Quantity:          o.QuantityValue.String(),
		AvailableQuantity: o.AvailableQuantityValue.String(),
		Price:             o.PriceValue.String(),
		OrderType:         o.OrderType,
		OrderId:           o.OrderId,
		CreateTime:        o.CreateTime,
	}
	if index := strings.Index(order.Price, "."); index != -1 {
		if index >= 1 {

			if order.Price[0] != '0' {
				order.Price = "0." + order.Price[:index] + order.Price[index+1:]
				exp = index
			}

		}
	} else {
		exp = len(order.Price)
		order.Price = "0." + order.Price
	}
	order.PriceExp = exp
	return order
}
func (o Order) Quantity() decimal.Decimal {
	return o.QuantityValue
}

//SetQuantity 设置数量
func (o *Order) SetQuantity(value string) {
	if value == "" {
		o.QuantityValue = decimal.Zero
		return
	}
	o.QuantityValue = decimal.RequireFromString(value)
}
func (o *Order) AddQuantity(value decimal.Decimal) {
	o.QuantityValue = o.QuantityValue.Add(value)
}
func (o *Order) SubQuantity(value decimal.Decimal) {
	o.QuantityValue = o.QuantityValue.Sub(value)
}

/* ----------------------获取和设置可用数量 ----------------*/
//AvailableQuantity 获取可用数量
func (o Order) AvailableQuantity() decimal.Decimal {
	return o.AvailableQuantityValue
}

//SetAvailableQuantity 设置可用数量
func (o *Order) SetAvailableQuantity(value string) {
	if value == "" {
		o.AvailableQuantityValue = decimal.Zero
		return
	}
	o.AvailableQuantityValue = decimal.RequireFromString(value)
}
func (o *Order) AddAvailableQuantiyu(value decimal.Decimal) {
	o.AvailableQuantityValue = o.AvailableQuantityValue.Add(value)
}
func (o *Order) SubAvailableQuantiyu(value decimal.Decimal) {
	o.AvailableQuantityValue = o.AvailableQuantityValue.Sub(value)
}

/* ----------------------获取和设置价格 ---------------------*/
//Price 获取价格
func (o Order) Price() decimal.Decimal {
	return o.PriceValue
}

//SetPrice 设置价格
func (o *Order) SetPrice(value string) {
	if value == "" {
		o.PriceValue = decimal.Zero
		return
	}
	o.PriceValue = decimal.RequireFromString(value)
}

type order struct {
	OrderId           string    `bson:"orderId" json:"orderId"`
	ProductKey        string    `bson:"productKey" json:"productKey"`               //产品id
	Quantity          string    `bson:"quantity" json:"quantity"`                   //数量
	AvailableQuantity string    `bson:"availableQuantity" json:"availableQuantity"` //可匹配数量
	Price             string    `bson:"price" json:"price"`                         //价格
	PriceExp          int       `bson:"priceExp" json:"priceExp"`
	OrderType         OrderType `bson:"orderType" json:"orderType"` //委托类型
	CreateTime        common.Time
}

func (o order) Trans() *Order {
	order := &Order{}
	order.SetQuantity(o.Quantity)
	order.SetAvailableQuantity(o.AvailableQuantity)
	if o.PriceExp != 0 {
		o.Price = o.Price[2:o.PriceExp+2] + "." + o.Price[o.PriceExp+2:]
	}
	order.SetPrice(o.Price)
	order.OrderId = o.OrderId
	order.ProductKey = o.ProductKey
	order.OrderType = o.OrderType
	order.CreateTime = o.CreateTime
	return order
}

var fileds = `
orderId varchar(255) not null,
productKey varchar(255) not null,
quantity varchar(255) not null,
availableQuantity varchar(255) not null,
price varchar(255) not null,
priceExp varchar(30) not null,
orderType int(8) not null,
createTime timestamp not null,
PRIMARY KEY (orderId)`
