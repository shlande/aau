package parser

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"regexp"
	"strconv"
	"strings"
)

func New() *Parse {
	return &Parse{
		tRex:     regexp.MustCompile(`[ |\[]([0-9]{1,3})[ |\]]`),
		gbRex:    regexp.MustCompile(`[\[|【]GB[\]|】]|简体|简中|简繁|簡繁|簡日|简日|[_|\[]CHS[_|\]]`),
		bigRex:   regexp.MustCompile(`[\[|【]BIG5[\]|】]|繁體|繁中|简繁|簡繁|繁日|[_|\[]CHT[_|\]]`),
		jpRex:    regexp.MustCompile(`簡日|繁日|简日`),
		p720Rex:  regexp.MustCompile(`720[p|P]`),
		p1080Rex: regexp.MustCompile(`1080[p|P]`),
		exRex:    regexp.MustCompile(`外挂|外掛`),
	}
}

type Parse struct {
	tRex     *regexp.Regexp
	gbRex    *regexp.Regexp
	bigRex   *regexp.Regexp
	jpRex    *regexp.Regexp
	p720Rex  *regexp.Regexp
	p1080Rex *regexp.Regexp
	exRex    *regexp.Regexp
}

func (p Parse) Parse(title string) (*Result, error) {
	name := parseName(title, false)
	// 分割所有的书名号
	eps := p.tRex.FindString(title)
	if len(eps) == 4 {
		eps = eps[1:3]
	}
	episode, _ := strconv.ParseInt(eps, 10, 8)

	// 这里不直接使用单个词，因为可能会出现很高记录的误匹配
	var lan data.Language
	gb := p.gbRex.MatchString(title)
	big := p.bigRex.MatchString(title)
	jp := p.jpRex.MatchString(title)
	if gb {
		lan = lan | data.GB
	}
	if big {
		lan = lan | data.BIG5
	}
	if jp {
		lan = lan | data.JP
	}

	var quality data.Quality
	p720 := p.p720Rex.MatchString(title)
	if p720 {
		quality = data.P720
	}
	// TODO：还有1080x60fps
	p1080 := p.p1080Rex.MatchString(title)
	if p1080 {
		quality = data.P1080
	}

	var sub = data.Internal
	external := p.exRex.MatchString(title)
	if external {
		sub = data.External
	}

	var tp = data.Episode
	if episode == 0 {
		tp = data.Full
	}
	return &Result{
		Name: name,
		// 当前不能区分分类
		Metadata: data.Metadata{
			Fansub:   parseFansub(title),
			Quality:  quality,
			Language: lan,
			SubType:  sub,
			// TODO:完成检测类型
			Type: tp, // ParseCategory(),
		},
		Episode: int(episode),
	}, nil
}

func parseName(title string, keep bool) string {
	var cname, left string
	if !keep {
		// 查找第一个右括号书名号位置,这个位置一般是字幕组名称截止的位置
		begin := regexp.MustCompile(`\]|】`).FindStringIndex(title)
		// 删除掉字幕组
		title = title[begin[1]:]
	}

	bw := regexp.MustCompile(`[\[|【]`).FindStringIndex(title)
	// 如果内容是以括号开始的，那么name应该就是这个括号内的内容
	if bw == nil {
		return title
	} else if bw[0] == 0 {
		ew := regexp.MustCompile(`\]|】`).FindStringIndex(title)
		cname = title[bw[1]:ew[0]]
		left = title[ew[1]:]
	} else {
		// 获取没被书名号括起来的内容
		cname = title[:bw[0]]
		left = title[bw[0]:]
	}
	// 判断是否是有效的名称,有些字幕组会加一月新番等多余的名字
	if regexp.MustCompile(`新番`).MatchString(cname) {
		// 这类型的暂时无法匹配
		return parseName(left, true)
	}
	if regexp.MustCompile(`^\s*$`).MatchString(cname) {
		return parseName(left, true)
	}
	// 对有效名称进行最后的处理
	// 通过反斜杠来分割名称中中英文
	{
		if raw := strings.Split(cname, "\\"); len(raw) != 1 {
			cname = raw[0]
		}
		if raw := strings.Split(cname, "/"); len(raw) != 1 {
			cname = raw[0]
		}
		if raw := strings.Split(cname, "_"); len(raw) != 1 {
			cname = raw[0]
		}
		// 删除仅限港澳台地区的字样
		{
			cname = regexp.MustCompile(`[\[|【|(|（](仅|僅)限港澳台(地区|地區)*[\)|\]|】|）]`).ReplaceAllString(cname, "")
		}
		// 删除两侧的空白符号
		cname = regexp.MustCompile(`^ *`).ReplaceAllString(cname, "")
		cname = regexp.MustCompile(` *$`).ReplaceAllString(cname, "")
	}
	return cname
}

func parseFansub(title string) []string {
	// 查找第一个右括号书名号位置,这个位置一般是字幕组名称截止的位置
	begin := regexp.MustCompile(`\]|】`).FindStringIndex(title)
	start := regexp.MustCompile(`\[|【`).FindStringIndex(title)
	// 删除掉字幕组
	title = title[start[1]:begin[0]]
	return []string{title}
}

func ParseCategory(title string) data.Type {
	//if title == "動畫" {
	//	return data.Episode
	//}
	//if title == "季度全集" {
	//	return data.Full
	//}
	return data.UnknownType
}

type Result struct {
	Name string
	// 分类
	data.Metadata
	// 集数
	Episode int
}
