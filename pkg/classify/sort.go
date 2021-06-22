package classify

import "github.com/shlande/dmhy-rss/pkg/parser"

type Options func(opt *Condition)

func After(episode int) Options {
	return func(opt *Condition) {
		opt.episode = func(i int) bool {
			if i > episode {
				return true
			}
			return false
		}
	}
}

func Equal(episode int) Options {
	return func(opt *Condition) {
		opt.episode = func(i int) bool {
			if i == episode {
				return true
			}
			return false
		}
	}
}

type Condition struct {
	Name    string
	episode func(int) bool
	Fansub  []string
	parser.Category
	parser.Quality
	parser.SubType
	parser.Language
}

// Find 根据option筛选出指定的item
func Find(items []*parser.Detail, condition *Condition, opts ...Options) []*parser.Detail {
	for _, v := range opts {
		v(condition)
	}
	var matches []*parser.Detail
	// var fansub,episode,quality,subtype,lan bool
	for _, i := range items {
		// 暂时不进行name匹配
		for _, fs := range condition.Fansub {
			for _, ifs := range i.Fansub {
				if fs == ifs {
					goto EPISODE
				}
			}
		}
		continue
	EPISODE:
		if condition.episode != nil && !condition.episode(i.Episode) {
			continue
		}
		if condition.Quality != parser.UnknownQuality && condition.Quality != i.Quality {
			continue
		}
		if condition.SubType != parser.UnknownSubType && i.SubType != condition.SubType {
			continue
		}
		if condition.Category != parser.UnknownCategory && i.Category != condition.Category {
			continue
		}
		// 暂时只支持And匹配
		if condition.Language != parser.UnknownLanguage && i.Language&condition.Language != condition.Language {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

// FindInCollection 在collection中查找
func FindInCollection(collection []*Collection, condition *Collection) {
	panic("impl me")
}
