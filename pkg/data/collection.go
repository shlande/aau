package data

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

type Collection struct {
	*Animation
	Metadata

	// Collection 的信息
	Latest     int
	LastUpdate time.Time

	Items []*Source
}

// NewCollection create a collection. only name is required
// it will automatically fill other info when adding first item.
func NewCollection(animation *Animation, metadata Metadata) *Collection {
	return &Collection{
		Animation: animation,
		Metadata:  metadata,
	}
}

func (c *Collection) Add(item *Source) error {
	// compare info
	if !c.compare(item) {
		return errors.New("item与collection信息不匹配")
	}
	if c.Has(item) {
		return nil
	}
	// add info
	c.Fansub = append(c.Fansub, diff(c.Fansub, item.Fansub)...)
	// 排序插入
	c.Items = append(c.Items, item)
	for i := len(c.Items) - 1; i > 0; i-- {
		if c.Items[i-1].Episode > c.Items[i].Episode {
			c.Items[i], c.Items[i-1] = c.Items[i-1], c.Items[i]
			continue
		}
		break
	}
	if c.LastUpdate.Before(item.CreateTime) {
		c.LastUpdate = item.CreateTime
	}
	if c.Latest < item.Episode {
		c.Latest = item.Episode
	}
	return nil
}

// TODO：同时还要检测source的名字是否一致
func (c *Collection) compare(item *Source) bool {
	if item == nil {
		return false
	}
	return c.Language == item.Language &&
		c.SubType == item.SubType &&
		c.Quality == item.Quality &&
		c.Type == item.Type
}

func (c *Collection) Has(item *Source) bool {
	for _, i := range c.Items {
		if i.Episode == item.Episode {
			return true
		}
	}
	return false
}

func (c *Collection) String() string {
	return fmt.Sprintf("%v", c.Metadata)
}

func (c *Collection) Id() string {
	data := md5.Sum([]byte(c.String()))
	return hex.EncodeToString(data[:])
}

func (c *Collection) GetLatest() *Source {
	for _, item := range c.Items {
		if item.Episode == c.Latest {
			return item
		}
	}
	return nil
}

func GetCollectionId(detail *Source) string {
	data := md5.Sum([]byte(fmt.Sprintf("%v-%v-%v-%v-%v", detail.Type, detail.Name, detail.Fansub, detail.Quality, detail.Language)))
	return hex.EncodeToString(data[:])
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
