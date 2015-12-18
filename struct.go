package goyate

import (
	"fmt"
	"net"
    "strings"
    "time"
    "strconv"
    "bufio"
    "os"
)

type Yate struct {
	yatestatus bool
	host string
	port int
	Con net.Conn
	Handlers map[string]func(*Message)
	Watchers map[string]func(*Message)
	Quit chan bool
	Daemon bool
}


func (yate *Yate) start_connection() {
	conn := yate.setup_conn()
	yate.Con = conn
	yate.yatestatus = true
	go yate.yate_out_handler()
}

func (yate *Yate) RetryHandlers() {
	for n, f:= range yate.Handlers {
		yate.Install(n, f)
	}
	
	for n, f := range yate.Watchers {
		yate.Installwatch(n, f)
	}
	
}

func (yate *Yate) yate_out_handler() {
	br := bufio.NewReader(yate.Con)
	for {
        msg, err := br.ReadString('\n')
        if err != nil {
            yate.Close()
            return
        }
        msg = strings.TrimSuffix(msg, "\n")
        t := &msg
        Log("GET: ", *t)
        yate.rawmsgparser(t)
	}
}

func (yate *Yate) setup_conn() net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", yate.host, yate.port))
	if err != nil {
		Log("ERROR: ", fmt.Sprint(err))
		if (yate.Daemon == true) {
			Log("STOP")
			os.Exit(1)
		}
		time.Sleep(5 *time.Second)
		return yate.setup_conn()
	}
	return conn
}

func (yate *Yate) Install(event string, handler func(*Message)) {
	msg := "%%" + fmt.Sprintf(">install::%s", event)
	go yate.Send(&msg)
	yate.Handlers[event] = handler
}

func (yate *Yate) Installwatch(event string, handler func(*Message)) {
	msg := "%%" +fmt.Sprintf(">watch:%s", event)
	go yate.Send(&msg)
	yate.Watchers[event] = handler
}

func (yate *Yate) Send(msg *string) {
	Log("SEND:", *msg)
	yate.Con.Write([]byte(*msg + "\n"))
}

func (yate *Yate) Close() {
    yate.Con.Close()
    Log("Disconnect")
    if (yate.Daemon == true) {
		Log("STOP")
		os.Exit(1)
		//yate.Quit <- true
		return
	}
	yate.start_connection()
	yate.RetryHandlers()
}

func (yate *Yate) Run() {
	<-yate.Quit
	Log("STOP")
}

func (yate *Yate) messageReceived(values *string) {
	/*
        Message parser.
	*/
	message := strings.SplitN(*values, ":", 5)
	newmsg := &Message{Mid: message[0], 
					   TimeStamp: message[1], 
					   Name: message[2],
					   Attrs: make(map[string]string),
					   Yate: yate,
                       Raw: *values}
	newmsg.Parse_attrs(&message[4])
	go yate.Handlers[newmsg.Name](newmsg)

}

func (yate *Yate) watchReceived(values []string) {
	/*
        Watch parser.
	*/
	newmsg := &Message{Mid: values[0], 
					   Returned: values[1], 
					   Name: values[2],
					   RetValue: values[3],
					   Attrs: make(map[string]string),
					   Yate: yate}
	newmsg.Parse_attrs(&values[4])
	go yate.Watchers[newmsg.Name](newmsg)
}

func (yate *Yate) watchOrResponseReceived(values *string) {
	message := strings.SplitN(*values, ":", 5)
	if message[0] == "" {
		yate.watchReceived(message)
	} else {
		yate.messageResponse(message)
	}
}

func (yate *Yate) messageResponse(values []string) {
	//TODO messageResponse
	Log("TODO", "messageResponse")
}

func (yate *Yate) installResponse(values *string) {
	message := strings.SplitN(*values, ":", 5)
	if (message[2] == "true" || message[2] == "ok") {
		Log(message[1], "Handler Succefully installed")
	} else {
		Log("Can't install handler for:", message[1])
	}
}

func (yate *Yate) watchinstallResponse(values *string) {
	message := strings.SplitN(*values, ":", 5)
	if (message[1] == "true" || message[1] == "ok") {
		Log(message[0], "watcher Succefully installed")
	} else {
		Log("Can't install watcher for:", message[1])
	}
}

func (yate *Yate) rawmsgparser(raw *string) {
	message := strings.SplitN(*raw, ":", 2)
	switch message[0] {
		case "%%>message":
			yate.messageReceived(&message[1])
		case "%%<message":
			yate.watchOrResponseReceived(&message[1])
		case "%%<install":
			yate.installResponse(&message[1])
		case "%%<watch":
			yate.watchinstallResponse(&message[1])
		//case "%%<setlocal":
			//Log("msgparser:", "setlocalResponse")
	}
}

func Parse_message(values *string) (newmsg *Message){
	message := strings.SplitN(*values, ":", 5)
	newmsg = &Message{Mid: message[0], 
					   TimeStamp: message[1], 
					   Name: message[2],
					   Attrs: make(map[string]string),
                       Raw: *values}
	newmsg.Parse_attrs(&message[4])
	return 
}


/************MESSAGE*************************************************/

type Message struct {
	Mid string
	TimeStamp string
	Name string
	RetValue string
	Attrs map[string]string
	Returned string
	Yate *Yate
    Raw string
}

func (msg *Message) Parse_attrs(attrs *string) {
	for _,v := range strings.Split(*attrs, ":") {
		raw := strings.SplitN(v, "=", 2)
        if len(raw) < 2 {
            Log(raw, "<<<<parse attrs error")
            continue
        }
		msg.Attrs[raw[0]] = strings.Replace(unescape(raw[1]), " ", "", -1)
	}
}

func (msg *Message) Format_attrs() (result string) {
	for k,v := range msg.Attrs {
		if k != "handlers" {
			result = result + ":" + k + "=" + escape(v)
		}
	}
	return
}
        
func (msg *Message) Format_message_response(returned bool, retvalue string) (*string) {
	result := fmt.Sprintf("%%%%<message:%s:%s:%s:", 
					     msg.Mid, 
					     strconv.FormatBool(returned), 
					     msg.Name)
	if retvalue != "" {
		result = result + escape(retvalue)
	}
	result = result + msg.Format_attrs()
	return &result
}

func (msg *Message) Ret(handled bool, retvalue string) {
	resp := msg.Format_message_response(handled, retvalue)
	msg.Yate.Send(resp)
}
