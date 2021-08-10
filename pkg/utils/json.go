package utils

import (
	"bytes"
	"encoding/json"
)

func UnescapeJsonMarshal(data interface{}) (bt []byte, err error) {
	byteBuf := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(byteBuf)
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(data)
	bt = byteBuf.Bytes()
	return
}
