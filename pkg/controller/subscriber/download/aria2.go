package download

import (
	"bytes"
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/utils"
	"net/http"
)

type Aria2 struct {
	rpc    string
	secret string
}

func NewAria2(rpc string, secret string) *Aria2 {
	return &Aria2{rpc: rpc, secret: secret}
}

func (a *Aria2) Add(path string, url string) error {
	return a.doRpc(context.Background(), "aria2.addUri", &aria2Options{Dir: path}, []string{url})
}

func (a *Aria2) doRpc(ctx context.Context, methods string, options *aria2Options, args ...interface{}) error {
	var params []interface{}
	if len(a.secret) != 0 {
		params = append(params, "token:"+a.secret)
	}
	for _, v := range args {
		params = append(params, v)
	}
	if options != nil {
		params = append(params, options)
	}

	data, err := utils.UnescapeJsonMarshal(&aria2Request{
		Version: "2.0",
		Id:      "from_aau",
		Method:  methods,
		Params:  params,
	})
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(a.rpc, "application/json", bytes.NewReader(data))
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// TODO:添加错误信息
		return errors.New("error")
	}
	if err != nil {
		return err
	}
	return nil
}

type aria2Options struct {
	Dir string `json:"dir"`
}

type aria2Request struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}
