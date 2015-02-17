# go-yate
Golang modules comunicate with YATE core via external protocol


example:

```
import (
	"./go-yate"
)

func authHandler(msg *goyate.Message) {
    //TODO something
	msg.Ret(true, "")
}

func main() {
    yate := goyate.Start("localhost", 5039)
	yate.Install("user.auth", authHandler)
	goyate.Run()
}
```
