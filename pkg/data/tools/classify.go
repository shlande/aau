package tools

import "github.com/shlande/dmhy-rss/pkg/data"

// Classify 把item归类成collection
func Classify(animation *data.Animation, items []*data.Source) (cls []*data.Collection) {
	for _, val := range classify(animation, items) {
		cls = append(cls, val)
	}
	return
}

func classify(animation *data.Animation, items []*data.Source) map[string]*data.Collection {
	var res = make(map[string]*data.Collection)
	for _, i := range items {
		// 先尝试创建cl
		//cl := NewCollection(i)
		cl := data.NewCollection(animation, i.Metadata)
		id := cl.Id()
		ocl, has := res[id]
		if !has {
			res[id] = cl
		} else {
			cl = ocl
		}
		// 这里不应该出现错误，因为cl相同肯定是可以添加的
		err := cl.Add(i)
		if err != nil && err != data.ErrEpisodeExist {
			panic(err)
		}
	}
	return res
}
