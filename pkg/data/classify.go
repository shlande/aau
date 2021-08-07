package data

// Classify 把item归类成collection
func Classify(items []*Source) (cls []*Collection) {
	for _, val := range classify(items) {
		cls = append(cls, val)
	}
	return
}

func classify(items []*Source) map[string]*Collection {
	var res = make(map[string]*Collection)
	for _, i := range items {
		// 先尝试创建cl
		//cl := NewCollection(i)
		cl := NewCollection(nil, Metadata{})
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
