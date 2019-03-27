package main

import (
	"flag"
	"github.com/deckarep/gosx-notifier"
	"github.com/fighterlyt/test/monitor/monitor"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	method    monitor.MonitorMethod
	arguments               = ""
	interval  time.Duration = time.Second * 10
)

const (
	argumentsHelp = "arguments参数必须为ip:port:port"
)

func main() {
	flag.Var(&method, "method", "method")
	flag.StringVar(&arguments, "arguments", "", argumentsHelp)
	flag.DurationVar(&interval, "interval", interval, "interval")
	flag.Parse()

	println(method, arguments, interval.String())
	switch method {
	case monitor.Portmonitor:
		if arguments == "" {
			panic(argumentsHelp)
		}
		if fields := strings.Split(arguments, ":"); len(fields) < 2 {
			panic(argumentsHelp)
		} else {
			monitors := make([]monitor.Monitor, 0, len(fields)-1)
			for i := 1; i < len(fields); i++ {
				println(fields[0], fields[i])
				m := monitor.NewPortMonitor(fields[0], fields[i])
				monitors = append(monitors, m)
			}
			wg := &sync.WaitGroup{}
			for _, m := range monitors {
				wg.Add(1)
				go func(m monitor.Monitor) {
					for {
						if err := m.Monitor(); err != nil {
							notify(arguments, err.Error())
						}
						time.Sleep(interval)


					}
					wg.Done()
				}(m)
			}
			wg.Wait()

		}
	default:

	}
}

func notify(title, msg string) {
	println("notify")
	note := gosxnotifier.NewNotification("监控故障")

	//Optionally, set a title
	note.Title = title

	//Optionally, set a subtitle
	note.Subtitle = msg

	//Optionally, set a sound from a predefined set.
	note.Sound = gosxnotifier.Basso

	//Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	note.Group = "com.unique.yourapp.identifier"

	//Optionally, set a sender (Notification will now use the Safari icon)
	note.Sender = "com.apple.Safari"

	//Optionally, an app icon (10.9+ ONLY)
	note.AppIcon = "gopher.png"

	//Optionally, a content image (10.9+ ONLY)
	note.ContentImage = "gopher.png"

	//Then, push the notification
	err := note.Push()

	//If necessary, check error
	if err != nil {
		log.Println("Uh oh!")
	}
}
