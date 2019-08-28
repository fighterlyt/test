package main

import (
	"github.com/fighterlyt/test/mongodb/line/charts"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"time"
)

func main() {
	if session, err := mgo.Dial("my_user:password123@localhost:27018/orderbook"); err != nil {
		panic(err.Error())
	} else {
		t := time.Now().Add(-24 * time.Hour)
		if data, err := getData(session, t.Unix()); err != nil {
			panic(err.Error())
		} else {
			names, values := processData(data)
			engine := gin.Default()

			now := time.Now()
			now = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())
			if data, err = getDataHour(session, now.Unix()); err != nil {
				panic(err.Error())
			} else {
				hourNames, hourValues := processDataMinute(data)
				engine.GET("/line", gin.WrapF(charts.Handle([]charts.LineArguments{
					{
						Title:  "user_assets_deal_log24小时",
						Name:   t.String(),
						Names:  names,
						Values: values,
					},
					{
						Title:  "user_assets_deal_log最近1小时每分钟",
						Name:   now.String(),
						Names:  hourNames,
						Values: hourValues,
					},
				})))
			}

			engine.Run(":8001")
		}
	}
}

type DealRecord struct {
	Count int `bson:"count"`
	ID    Id  `bson:"_id"`
}
type Id struct {
	Year   int
	Month  int
	Day    int
	Hour   int
	Minute int
}

func getData(session *mgo.Session, start int64) ([]DealRecord, error) {
	var records []DealRecord
	collection := session.DB("orderbook").C("user_assets_deal_log")

	if err := collection.Pipe([]bson.M{
		{
			"$match": bson.M{"timestamp": bson.M{"$gte": start}},
		},
		{
			"$project": bson.M{
				"date": bson.M{
					"$dateFromParts": bson.M{
						"year":   1970.0,
						"second": "$timestamp",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year": bson.M{
						"$year": "$date",
					},
					"month": bson.M{
						"$month": "$date",
					},
					"day": bson.M{
						"$dayOfMonth": "$date",
					},
					"hour": bson.M{
						"$hour": "$date",
					},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}).All(&records); err != nil {
		return nil, err
	}
	return records, nil
}

func getDataHour(session *mgo.Session, start int64) ([]DealRecord, error) {
	println("start", start)
	var records []DealRecord
	collection := session.DB("orderbook").C("user_assets_deal_log")

	if err := collection.Pipe([]bson.M{
		{
			"$match": bson.M{"timestamp": bson.M{"$gte": start}},
		},
		{
			"$project": bson.M{
				"date": bson.M{
					"$dateFromParts": bson.M{
						"year":   1970.0,
						"second": "$timestamp",
					},
				},
			},
		},
		{
			"$group": bson.M{
				"_id": bson.M{
					"year": bson.M{
						"$year": "$date",
					},
					"month": bson.M{
						"$month": "$date",
					},
					"day": bson.M{
						"$dayOfMonth": "$date",
					},
					"hour": bson.M{
						"$hour": "$date",
					},
					"minute": bson.M{
						"$minute": "$date",
					},
				},
				"count": bson.M{
					"$sum": 1,
				},
			},
		},
		{
			"$sort": bson.M{"_id": 1},
		},
	}).All(&records); err != nil {
		return nil, err
	}
	return records, nil
}
func processData(records []DealRecord) ([]string, []int) {
	names := make([]string, 0, len(records))
	values := make([]int, 0, len(records))
	for _, record := range records {
		t := time.Date(record.ID.Year, time.Month(record.ID.Month), record.ID.Day, record.ID.Hour, 0, 0, 0, time.FixedZone("here", 0))
		t = t.Add(time.Hour * 8)
		names = append(names, t.Format("2006-01-02-15"))
		values = append(values, record.Count)
	}
	return names, values
}
func processDataMinute(records []DealRecord) ([]string, []int) {
	println(len(records))
	names := make([]string, 0, len(records))
	values := make([]int, 0, len(records))
	for _, record := range records {
		t := time.Date(record.ID.Year, time.Month(record.ID.Month), record.ID.Day, record.ID.Hour, record.ID.Minute, 0, 0, time.FixedZone("here", 0))
		t = t.Add(time.Hour * 8)
		names = append(names, t.Format("2006-01-02-15-04"))
		values = append(values, record.Count)
	}
	return names, values
}
