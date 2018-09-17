package main

import (
	"net/http"
	"fmt"
	"os"
)

func main(){
	http.HandleFunc("/",func(resp http.ResponseWriter,req *http.Request){
		if req.Method!=http.MethodPost{
			fmt.Fprintln(resp,"只支持POST")
			resp.WriteHeader(400)
		}else{
			req.Write(os.Stdout)
			resp.WriteHeader(200)
		}
	})
 	http.ListenAndServe(":12345",nil)
 }
