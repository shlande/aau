package classify

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/parser"
	"time"
)

// NewCollection create a collection. only name is required
// it will automatically fill other info when adding first item.
func NewCollection(item *parser.Detail) *Collection {
	return &Collection{
		Fansub:     item.Fansub,
		Quality:    item.Quality,
		Category:   item.TitleInfo.Category,
		SubType:    item.SubType,
		Language:   item.Language,
		Latest:     item.Episode,
		LastUpdate: time.Now(),
		Items:      []*parser.Detail{item},
	}
}

// Collection 包含了单个种类的所有record
// 自内容的所有信息必须相似，否则无法添加
type Collection struct {
	// 重复保存了的内容
	Name   string
	Fansub []string
	parser.Quality
	parser.Category
	parser.SubType
	parser.Language
	// Collection 的信息
	Latest     int
	LastUpdate time.Time

	Items []*parser.Detail
}

func (c *Collection) Add(item *parser.Detail) error {
	// compare info
	if !c.compare(item) {
		return errors.New("item与collection信息不匹配")
	}
	if c.Has(item) {
		return nil
	}
	// add info
	c.Fansub = append(c.Fansub, diff(c.Fansub, item.Fansub)...)
	c.Items = append(c.Items, item)
	c.LastUpdate = time.Now()
	if c.Latest < item.Episode {
		c.Latest = item.Episode
	}
	return nil
}

func (c *Collection) compare(item *parser.Detail) bool {
	return c.Language == item.Language &&
		c.SubType != item.SubType &&
		c.Quality != item.Quality &&
		c.Category != item.TitleInfo.Category
}

func (c *Collection) Has(item *parser.Detail) bool {
	for _, i := range c.Items {
		if i.Episode == item.Episode {
			return true
		}
	}
	return false
}

func (c *Collection) String() string {
	return fmt.Sprintf("%v-%v-%v-%v-%v", c.Category, c.Name, c.Fansub, c.Quality, c.Language)
}

func (c *Collection) Id() string {
	return string(md5.New().Sum([]byte(c.String())))
}

// diff find the // of new and old.
func diff(old []string, new []string) []string {
	var res []string
	var has = make(map[string]struct{})
	for _, v := range old {
		has[v] = struct{}{}
	}
	for _, v := range new {
		if _, ok := has[v]; !ok {
			res = append(res, v)
		}
	}
	return res
}
