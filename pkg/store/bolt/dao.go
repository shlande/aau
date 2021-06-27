package bolt

import (
	"encoding/json"
	"github.com/shlande/dmhy-rss/pkg/classify"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"time"
	"unsafe"
)

func NewInfo(collection *classify.Collection) *Info {
	return (*Info)(unsafe.Pointer(collection))
}

type Info struct {
	Name   string
	Fansub []string
	parser.Quality
	parser.Category
	parser.SubType
	parser.Language
	// Collection 的信息
	Latest     int
	LastUpdate time.Time
}

func (i *Info) Encode() []byte {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

func decodeInfo(data []byte, info *Info) error {
	err := json.Unmarshal(data, info)
	if err != nil {
		panic(err)
	}
	return nil
}

func NewEpisode(detail *parser.Detail) *Episode {
	return (*Episode)(unsafe.Pointer(detail))
}

type Episode struct {
	parser.Detail
}

func (e *Episode) Encode() []byte {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return data
}

func decodeEpisode(data []byte, ep *Episode) error {
	err := json.Unmarshal(data, ep)
	if err != nil {
		panic(err)
	}
	return nil
}
