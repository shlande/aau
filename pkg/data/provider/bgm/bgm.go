package bgm

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/provider"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func New(animationInterface store.AnimationInterface) *Provider {
	p := &Provider{
		data:               make(map[string]*bgmDataItem),
		AnimationInterface: animationInterface,
	}
	p.updateData()
	go func() {
		updateTicker := time.NewTicker(time.Hour * 24 * 7)
		for {
			<-updateTicker.C
			p.updateData()
		}
	}()
	return p
}

type Provider struct {
	store.AnimationInterface
	lastUpdate time.Time
	ld         []*bgmDataItem
	data       map[string]*bgmDataItem
	cache      map[string]*time.Time
}

func (p *Provider) Session(_ context.Context, year int, session provider.Session) (anms []*data.Animation, err error) {
	for _, v := range p.sortByTime(getSessionTime(year, int(session))) {
		anm, _ := p.get(v.GetUniId())
		anms = append(anms, anm)
	}
	return anms, nil
}

func (p *Provider) sortByTime(begin, end time.Time) (anms []*bgmDataItem) {
	for _, v := range p.data {
		dt, _, err := v.findAir()
		if err != nil {
			logrus.Error(err)
			continue
		}
		if dt.Before(end) && dt.After(begin) {
			anms = append(anms, v)
		}
	}
	return anms
}

func (p *Provider) Search(ctx context.Context, keywords string) (anms []*data.Animation, err error) {
	sr, err := p.search(ctx, keywords)
	if err != nil {
		return nil, err
	}
	for _, v := range sr.List {
		uniId := generateUniId(v.GetBgmId())
		anm, err := p.get(uniId)
		if err != nil {
			logrus.Error(err)
		}
		anms = append(anms, anm)
	}
	return anms, nil
}

func (p *Provider) getBgmData(id string) *bgmDataItem {
	return p.data[id]
}

func (p *Provider) get(uniId string) (*data.Animation, error) {
	anm, err := p.AnimationInterface.Get(uniId)
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}
	// 如果这个词条以前没有查过
	bgmId := getRawBangumiId(uniId)
	if err == store.ErrNotFound {
		bdi := p.getBgmData(bgmId)
		if bdi == nil {
			logrus.Errorf("bgmdata缺少信息 bgmId:%v", bgmId)
			bdi = &bgmDataItem{}
		}
		s, err := p.getSpec(bgmId)
		// 如果获取信息出现错误的话，就使用本地的残缺版本
		if err != nil {
			anm, err := bdi.generate("无法获取数据", 0)
			if err != nil {
				logrus.Error(err)
			}
			return anm, nil
		}
		// 整合悉信息，生成animation
		dt, bk, err := bdi.findAir()
		if err != nil {
			logrus.Error(err)
		}
		anm = s.generate(dt, bk, bdi.Type)
		if err := p.AnimationInterface.Save(anm); err != nil {
			logrus.Error(err)
		}
		return anm, nil
	}
	return anm, nil
}

func (p *Provider) Get(_ context.Context, uniId string) (*data.Animation, error) {
	return p.get(uniId)
}

func (p *Provider) updateData() {
	res, err := getBgmData()
	if err != nil {
		logrus.Error("更新bgmdata失败", err)
		return
	}
	p.ld = res.Items
	for _, v := range res.Items {
		p.data[v.GetBgmId()] = v
	}
	p.lastUpdate = time.Now()
	logrus.Info("bgmdata更新成功")
}

func (p *Provider) search(ctx context.Context, keywords string) (*SearchResult, error) {
	req, _ := http.NewRequestWithContext(
		ctx, http.MethodGet,
		"https://api.bgm.tv/search/subject/"+url.QueryEscape(keywords)+"?type=2",
		nil,
	)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var searchResult = &SearchResult{}
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return searchResult, json.Unmarshal(bt, searchResult)
}

func (p *Provider) getSpec(bgmId string) (*subjectSmall, error) {
	resp, err := http.Get("https://api.bgm.tv/subject/" + bgmId)
	if err != nil {
		return nil, err
	}
	var ss = &subjectSmall{}
	return ss, json.NewDecoder(resp.Body).Decode(ss)
}

func getBgmData() (*bgmDataResponse, error) {
	resp, err := http.Get("https://cdn.jsdelivr.net/npm/bangumi-data@0.3/dist/data.json")
	if err != nil {
		logrus.Error("无法更新番剧信息", err)
		return nil, err
	}
	var rsp = &bgmDataResponse{}
	err = json.NewDecoder(resp.Body).Decode(rsp)
	if err != nil {
		logrus.Error("无法更新番剧信息")
	}
	return rsp, nil
}

func ParseBroadcast(bs string) (airTime time.Time, airBreak time.Duration, err error) {
	raw := strings.Split(bs, "/")
	if len(raw) != 3 || len(raw[2]) != 3 {
		err = errors.New("错误的broadcast格式")
	}
	airTime, err = time.Parse(time.RFC3339, raw[1])
	if err != nil {
		return
	}

	airTime = airTime.Local()

	ab, err := strconv.ParseInt(raw[2][1:2], 10, 64)
	airBreak = time.Hour * 24 * time.Duration(ab)

	return
}

func getSessionTime(year int, session int) (start time.Time, end time.Time) {
	// 就一般理性而言，番剧都会在某一个月集中开始发布，因此在计算开始的时候，把范围往前面挪动会更加合适
	if session == 0 {
		year = year - 1
	}
	return time.Date(year, time.Month((11+session*3)%12+1), 0, 0, 0, 0, 0, time.Local),
		time.Date(year+session/3, time.Month((11+(session+1)*3)%12+1), 0, 0, 0, 0, 0, time.Local)
}
