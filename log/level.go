package main

import "github.com/fighterlyt/log"

func main(){
	logger:=NewLogger("test")
	logger.SetMinLevel(log.Fatal)
	logtest(logger)
}
type Logger struct {
	*log.Logger
	name string
}

func NewLogger(name string) *Logger {
	return &Logger{
		Logger: log.NewLogger(3),
		name:   name,
	}
}

func (m Logger) Log(level log.Level, info ...interface{}) {
	m.Logger.NewEntry().AddField("模块", m.name).Output(level, info...)
}
func (m Logger) Logf(level log.Level, format string, value ...interface{}) {
	m.Logger.NewEntry().AddField("模块", m.name).Outputf(level, format, value...)
}

func logtest(logger *Logger){
	logger.Log(log.Info,"info")

}