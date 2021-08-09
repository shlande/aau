package bgm

import (
	"github.com/shlande/dmhy-rss/pkg/data"
	"strconv"
	"strings"
	"time"
)

type subjectSmall struct {
	Id       int
	Url      string
	Name     string
	NameCN   string `json:"name_cn"`
	Summary  string `json:"summary"`
	EpsCount int    `json:"eps_count"`
}

func (s *subjectSmall) GetBgmId() string {
	return strconv.Itoa(s.Id)
}

func (s *subjectSmall) generate(airDate time.Time, airBreak time.Duration, tp string) *data.Animation {
	return &data.Animation{
		Id:            generateUniId(s.GetBgmId()),
		Name:          s.Name,
		Translated:    s.NameCN,
		Summary:       s.Summary,
		AirDate:       airDate,
		AirBreak:      airBreak,
		TotalEpisodes: s.EpsCount,
		Category:      tp,
	}
}

type SearchResult struct {
	Results int `json:"results"`
	List    []*subjectSmall
}

type bgmDataResponse struct {
	Items []*bgmDataItem
}

type bgmDataItem struct {
	Title           string
	TitleTranslated struct {
		ZHS []string `json:"zh-Hans"`
	} `json:"titleTranslate"`
	Begin     string
	Broadcast string `json:"broadcast"`
	Type      string `json:"type"`
	Sites     []struct {
		Site      string
		Id        string
		Begin     string
		Broadcast string `json:"broadcast"`
	}
}

func (b *bgmDataItem) GetUniId() string {
	return generateUniId(b.GetBgmId())
}

func (b *bgmDataItem) GetBgmId() string {
	// 查找bgmId
	for _, v := range b.Sites {
		if v.Site == "bangumi" {
			return v.Id
		}
	}
	return ""
}

func (b *bgmDataItem) generate(summary string, eps int) (*data.Animation, error) {
	dt, bk, err := b.findAir()
	var anm = &data.Animation{
		Id:            b.GetUniId(),
		Name:          b.Title,
		Summary:       summary,
		AirDate:       dt,
		AirBreak:      bk,
		TotalEpisodes: eps,
		Category:      b.Type,
	}
	if len(b.TitleTranslated.ZHS) > 0 {
		anm.Translated = b.TitleTranslated.ZHS[0]
	}
	return anm, err
}

func (b *bgmDataItem) findAir() (date time.Time, bt time.Duration, err error) {
	if len(b.Broadcast) != 0 {
		return ParseBroadcast(b.Broadcast)
	}
	for _, v := range b.Sites {
		if len(v.Broadcast) != 0 {
			return ParseBroadcast(v.Broadcast)
		}
	}
	return
}

func generateUniId(bgmId string) string {
	return "bgm-" + bgmId
}

func getRawBangumiId(uinId string) string {
	return strings.Split(uinId, "-")[1]
}
