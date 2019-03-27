package monitor

import (
	"fmt"
	"time"

	"github.com/reiver/go-telnet"
)

type MonitorMethod string

const (
	Portmonitor MonitorMethod = "port"
)

var (
	allMethods = map[string]struct{}{
		string(Portmonitor): struct{}{},
	}
)

func (m *MonitorMethod) String() string {
	switch *m {
	case Portmonitor:
		return "端口检测"
	default:
		return "未知"
	}
}
func (m *MonitorMethod) Set(value string) error {
	if _, exist := allMethods[value]; exist {
		*m = MonitorMethod(value)
	} else {
		return fmt.Errorf("方法[%s]不存在", value)
	}
	return nil
}

type Monitor interface {
	SetInterval(duration time.Duration)
	Monitor() error
}

type PortMonitor struct {
	port     string
	ip       string
	interval time.Duration
	caller   telnet.Caller
}

func NewPortMonitor(ip, port string) Monitor {
	return &PortMonitor{
		port:   port,
		ip:     ip,
		caller: telnet.StandardCaller,
	}
}

func (p *PortMonitor) SetInterval(duration time.Duration) {
	p.interval = duration
}

func (p *PortMonitor) Monitor() error {
	//@TOOD: replace "example.net:5555" with address you want to connect to.
	if conn, err := telnet.DialTo(p.ip + ":" + p.port); err != nil {
		println(time.Now().String()+err.Error())
		return err
	} else {
		if _, err = conn.Write([]byte("hello")); err != nil {
			return err
		} else {

			return conn.Close()
		}

	}
}
