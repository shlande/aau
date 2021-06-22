package classify

import (
	"github.com/shlande/dmhy-rss/pkg/parser"
)

type Option struct {
	Name    string
	Episode int
	Fansub  []string
	parser.Category
	parser.Quality
	parser.SubType
	parser.Language
}

// Classify 把item归类成collection
func Classify(items []*parser.Detail) (cls []*Collection) {
	for _, val := range classify(items) {
		cls = append(cls, val)
	}
	return
}

func classify(items []*parser.Detail) map[string]*Collection {
	var res = make(map[string]*Collection)
	for _, i := range items {
		// 先尝试创建cl
		cl := NewCollection(i)
		id := cl.Id()
		ocl, has := res[id]
		if !has {
			res[id] = cl
			continue
		}
		// 这里不应该出现错误，因为cl相同肯定是可以添加的
		_ = ocl.Add(i)
	}
	return res
}
