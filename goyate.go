package goyate

import (
	"fmt"
    "flag"
    "log"
    "time"
    "os"
    "bufio"
    "unicode/utf8"
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
		c, size := utf8.DecodeRuneInString(str[i:])
        if(c < 32  || c == ':') {
            c = rune(c + 64)
            s = s + "%"
        } else if ( c == '%' ){
            s = s + string(c)
        }
        s = s + string(c)
        i = i + size
		
	}
	return s
}

func unescape(str string) (string) {
    s := ""
    n := len(str)
    i := 0
    for i < n {
		c, size := utf8.DecodeRuneInString(str[i:])
        if (c == '%') {
            i = i + size
            c, size = utf8.DecodeRuneInString(str[i:])
            if (c != '%') {
                c = rune(c - 64)
	        }
	    }
        s = s + string(c)
        i = i + size
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

func StartScript() *Myate {
	yate := &Myate{status:true,
		           Handlers:make(map[string]func(*Myate, *Message)),
		           stdin:bufio.NewScanner(os.Stdin),
		           stdout:bufio.NewWriter(os.Stdout),
		           stderr:bufio.NewScanner(os.Stderr)}
	return yate
}

func Run() {
	for {
		time.Sleep(1*1e9)
	}
}
