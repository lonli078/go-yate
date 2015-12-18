package goyate

import (
	"fmt"
    "flag"
    "log"
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

func escape(str string) (string) {
    str = str + ""
    s := ""
    n := len(str)
    i := 0
    for i < n {
        c := string(str[i])
        if([]rune(c)[0] < 32  || c == ":") {
            c = string([]rune(c)[0] + 64)
            s = s + "%"
        } else if ( c == "%" ){
            s = s + c
        }
        s = s + c
        i = i + 1
		
	}
	return s
}
    
func unescape(str string) (string) {
    s := ""
    n := len(str)
    i := 0
    for i < n {
        c := string(str[i])
        if (c == "%") {
            i = i + 1
            c = string(str[i])
            if (c != "%") {
                c = string([]rune(c)[0] - 64)
	    }
	}
        s = s + c
        i = i + 1
    }
    return s
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
