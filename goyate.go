package goyate

import (
	"fmt"
    "flag"
    "log"
    "strings"
    "time"
    "os"
)
var DEBUG = flag.Bool("d", false, "set the debug modus( print informations )")
var LOGFILE = flag.String("w", "", "set the log file")

func Log(v ...interface{}) {
	if *DEBUG == true {
		log.Printf("GOYATE: %s", fmt.Sprint(v))
	}
}

func logger_start() {
	if len(*LOGFILE) > 0 {
		f, err := os.OpenFile(*LOGFILE, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if err != nil {
			Log("error opening file: ", fmt.Sprint(err))
		}
		log.SetOutput(f)
	}
}

func unescape(s string) (string) {
	return strings.Replace(s, "%z", ":", -1)
}
func escape(s string) (string) {
	return strings.Replace(s, ":", "%z", -1)
}

func Start(host string, port int, daemon bool) *Yate {
	flag.Parse()
	logger_start()
	yate := &Yate{yatestatus:false,
		          host:host,
		          port:port,
		          Daemon:daemon,
		          Handlers:make(map[string]func(*Message)),
		          Watchers:make(map[string]func(*Message))}
	yate.start_connection()
	return yate
}

func Run() {
	for {
		time.Sleep(1*1e9)
	}
}
