package record

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"sync"
)

func NewJsonKVFromFile(path string) *JsonKV {
	fil, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return NewJsonKV(fil)
}

func NewJsonKV(file io.WriteSeeker) *JsonKV {
	// TODO: 检测json是否正常，否则报错
	// 尝试获取内容
	jkv := &JsonKV{WriteSeeker: file, empty: true}
	if reader, ok := file.(io.Reader); ok {
		_, err := file.Seek(0, io.SeekStart)
		if err != nil {
			panic(err)
		}
		temp := make([]byte, 10)
		n, err := reader.Read(temp)
		if err != nil && err != io.EOF {
			panic(err)
		}
		jkv.empty = !regexp.MustCompile("{").Match(temp[:n])
	}
	return jkv
}

type JsonKV struct {
	empty bool
	chunk *bytes.Buffer
	sync.Mutex
	io.WriteSeeker
}

func (r *JsonKV) Set(key string, data interface{}) (err error) {
	var dt []byte
	switch temp := data.(type) {
	case []byte:
		dt = temp
	default:
		dt, err = json.Marshal(data)
		if err != nil {
			return err
		}
	}
	return r.append(key, dt)
}

func (r *JsonKV) append(key string, data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if r.chunk == nil {
		r.chunk = bytes.NewBuffer(make([]byte, len(key)))
	}
	r.chunk.Reset()
	json.HTMLEscape(r.chunk, []byte(key))
	if r.empty {
		r.WriteSeeker.Seek(0, io.SeekStart)
		_, err = r.Write([]byte("[{\"" + string(r.chunk.Bytes()) + "\":"))
	} else {
		r.WriteSeeker.Seek(-1, io.SeekEnd)
		_, err = r.Write([]byte(",{\"" + string(r.chunk.Bytes()) + "\":"))
	}
	if err != nil {
		return err
	}
	_, err = r.Write(data)
	if err != nil {
		return err
	}
	_, err = r.Write([]byte("}]"))
	if err != nil {
		return err
	}
	return nil
}
