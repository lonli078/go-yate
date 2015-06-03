package goyate

import (
  "testing"
)

func TestMessageParseAttrs(t *testing.T) {
    msg := Message{Attrs: make(map[string]string)} 
    text := "caller=125:called:::called=125"
    msg.Parse_attrs(&text)
    if msg.Attrs["called"] != "125" {
        t.Error("Error in parse!.")
    }
}
