package main

import (
	"fmt"
	"net/http"
	"os"
)

var gitlab_domain = "http://172.105.195.24";

var private_token = "vyPuzfkBVziUG-wcgyg9";

var default_page = "/activity";

var order_by = "last_activity_at";

var nb_project = "30";

var font_size = "15";

var path = "/api/v4/projects/%s/events&target_type=project";

func main() {
	checkProject("18")
}

func checkProject(project string) string {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf(gitlab_domain+path,project), nil)
	println(fmt.Sprintf(gitlab_domain+path,project))
	req.Header.Add("Private-Token", private_token)

	if resp, err := http.DefaultClient.Do(req); err != nil {
		panic(err.Error())
	} else {
		resp.Write(os.Stdout)
		return ""
	}
}
