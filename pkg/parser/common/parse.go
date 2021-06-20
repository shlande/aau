package common

import (
	"github.com/shlande/dmhy-rss/pkg/parser"
	"github.com/shlande/dmhy-rss/pkg/provider"
	"regexp"
	"strconv"
	"strings"
)

type Parse struct{}

func (p Parse) ParseTitle(title string) (*parser.TitleInfo, error) {
	return parseTitle(title), nil
}

func (p Parse) Parse(infos ...*provider.Info) (res []*parser.Detail, err error) {
	for _, info := range infos {
		ti, err := p.ParseTitle(info.Title)
		if err != nil {
			return nil, err
		}
		ti.Category = parser.ParseCategory(info.RawCategory)
		res = append(res, &parser.Detail{
			TitleInfo: ti,
			Info:      info,
		})
	}
	return
}

// 尝试解析头部，获取到一些信息
func parseTitle(title string) *parser.TitleInfo {
	name := parseName(title, false)
	// 分割所有的书名号
	eps := regexp.MustCompile(`[ |\[]([0-9]{1,3})[ |\]]`).FindString(title)
	if len(eps) == 4 {
		eps = eps[1:3]
	}
	episode, _ := strconv.ParseInt(eps, 10, 8)

	// 这里不直接使用单个词，因为可能会出现很高记录的误匹配
	var lan parser.Language
	gb := regexp.MustCompile(`[\[|【]GB[\]|】]|简体|简中|简繁|簡繁|簡日|简日|[_|\[]CHS[_|\]]`).MatchString(title)
	big := regexp.MustCompile(`[\[|【]BIG5[\]|】]|繁體|繁中|简繁|簡繁|繁日|[_|\[]CHT[_|\]]`).MatchString(title)
	jp := regexp.MustCompile(`簡日|繁日|简日`).MatchString(title)
	if gb {
		lan = lan | parser.GB
	}
	if big {
		lan = lan | parser.BIG5
	}
	if jp {
		lan = lan | parser.JP
	}

	var quality parser.Quality
	p720 := regexp.MustCompile(`720[p|P]`).MatchString(title)
	if p720 {
		quality = parser.P720
	}
	// TODO：还有1080x60fps
	p1080 := regexp.MustCompile(`1080[p|P]`).MatchString(title)
	if p1080 {
		quality = parser.P1080
	}
	// 1440p匹配
	//k2 := regexp.MustCompile(``)

	var sub = parser.Internal
	external := regexp.MustCompile(`外挂|外掛`).MatchString(title)
	if external {
		sub = parser.External
	}
	return &parser.TitleInfo{
		Name: name,
		// 当前不能区分分类
		Category: parser.UnknownCategory,
		Language: lan,
		Quality:  quality,
		Episode:  int(episode),
		SubType:  sub,
	}
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
	if bw[0] == 0 {
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
