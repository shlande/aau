package classify

import "github.com/shlande/dmhy-rss/pkg/parser"

// Find 根据option筛选出指定的item
func Find(items []*parser.Detail, option *Option) []*parser.Detail {
	var matches []*parser.Detail
	// var fansub,episode,quality,subtype,lan bool
	for _, i := range items {
		// 暂时不进行name匹配
		for _, fs := range option.Fansub {
			for _, ifs := range i.Fansub {
				if fs == ifs {
					goto EPISODE
				}
			}
		}
		continue
	EPISODE:
		if option.Episode != 0 && i.Episode != option.Episode {
			continue
		}
		if option.Quality != parser.UnknownQuality && option.Quality != i.Quality {
			continue
		}
		if option.SubType != parser.UnknownSubType && i.SubType != option.SubType {
			continue
		}
		if option.Category != parser.UnknownCategory && i.Category != option.Category {
			continue
		}
		// 暂时只支持And匹配
		if option.Language != parser.UnknownLanguage && i.Language&option.Language != option.Language {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

// FindInCollection 在collection中查找
func FindInCollection(collection []*Collection, option *Option) {
	panic("impl me")
}
