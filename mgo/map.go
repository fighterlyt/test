package main

import "github.com/globalsign/mgo"

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err.Error())
	}

	c := session.DB("test").C("map")

	if err = c.Insert(map[int]map[int]string{
		1: {
			1: "1",
		},
	}); err != nil {
		panic(err.Error())
	}
	if err = c.Insert(map[int]string{

		1: "1",
	}); err != nil {
		panic(err.Error())
	}
}
