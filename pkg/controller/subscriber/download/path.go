package download

import (
	"fmt"
	"github.com/shlande/dmhy-rss/pkg/data"
	pt "path"
	"strconv"
)

// PathInterface 用来控制下载后文件保存到那个位置
type PathInterface interface {
	// Path 获取下载路径
	Path(resource *data.Collection) string
	// Tag 获取下载标签
	Tag(resource *data.Collection) string
	// Name 获取文件保存时名称
	Name(resource *data.Source, collection *data.Collection) string
}

func NewPath(base string) *path {
	return &path{base: base}
}

type path struct {
	base string
}

func (p *path) Name(resource *data.Source, collection *data.Collection) string {
	return ""
}

func (p *path) Path(collection *data.Collection) string {
	var m string
	y, s := collection.GetSession()

	s = s*3 + 1
	if s < 10 {
		m = "0" + strconv.Itoa(s)
	} else {
		m = strconv.Itoa(s)
	}
	return pt.Join(p.base, fmt.Sprintf("%v-%v-%v-%v", y, m, collection.Translated, collection.Fansub))
}

func (p *path) Tag(resource *data.Collection) string {
	return ""
}
