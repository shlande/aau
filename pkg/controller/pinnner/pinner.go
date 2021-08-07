package pinnner

import (
	"context"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/data"
	"github.com/shlande/dmhy-rss/pkg/data/tools"
	"log"
	"time"
)

type Pinner interface {
	Pin(ctx context.Context, animation *data.Animation, strategy Strategy) error
	Unpin(animation data.Animation) error
}

type Strategy struct {
	// 倾向的语言
	Language []data.Language
	// 倾向的翻译组
	Fansub []string
	// 倾向的画质
	Quality []data.Quality
	// 允许不选择的期限
	Tolerate time.Duration
	Type     []data.Type
}

var defaultStrategy = Strategy{
	Language: []data.Language{data.GB, data.BIG5},
	Fansub:   []string{},
	Quality:  []data.Quality{data.P1080, data.P720, data.K2},
	// 允许在一周内作出选择
	Tolerate: time.Hour * 24 * 7,
}

type pinner struct {
	tools.CollectionProvider
	finished map[string]struct{}
	pined    map[string]*data.Animation
	stg      map[string]Strategy
}

func (p pinner) Pin(ctx context.Context, animation *data.Animation, strategy Strategy) error {
	if _, has := p.finished[animation.Name]; has {
		return errors.New("番剧已经pin过且完成了")
	}
	p.pined[animation.Name] = animation
	p.stg[animation.Name] = strategy
	return nil
}

func (p pinner) Unpin(animation data.Animation) error {
	if _, has := p.pined[animation.Name]; !has {
		return errors.New("番剧还没有被pin")
	}
	delete(p.pined, animation.Name)
	delete(p.stg, animation.Name)
	return nil
}

// find best
//func (p pinner) Update() {
//	p.Provider.Keywords()
//}

func (p *pinner) tryFindBest(animation *data.Animation, strategy Strategy) *data.Collection {
	ctx, cf := context.WithTimeout(context.Background(), time.Second*10)
	defer cf()
	cls, err := p.CollectionProvider.Search(ctx, animation)
	if err != nil {
		log.Println(err)
	}
	var points = make([]int, len(cls))
	for _, v := range cls {
		points = append(points, mark(strategy, v))
	}
	// 找出做好的分数
	var maxPoint, maxIndex int
	for index, point := range points {
		if point > maxPoint {
			maxIndex = index
			maxPoint = point
		}
	}

	// TODO:如果没有找到好的分数，而且coolection数量小于2，那么就暂时跳过
	// 如果还允许等待
	if maxPoint < 10 && animation.AirDate.Add(strategy.Tolerate).After(time.Now()) && len(cls) < 3 {
		return nil
	}
	return cls[maxIndex]
}

func mark(strategy Strategy, collection *data.Collection) int {
	var point, lanPoint, qualityPoint, fansubPoint, typePoint, subPoint int

	// 计算画质分
	//for i, v := range strategy.Quality {
	//	if v == collection.Quality && qualityPoint == 0 {
	//		qualityPoint = len(strategy.Quality) - i
	//		break
	//	}
	//}
	switch collection.Quality {
	case data.P1080:
		qualityPoint = 3
	case data.P720:
		qualityPoint = 2
	case data.K2:
		qualityPoint = 1
	}

	// 计算汉化组分
	for i, v := range strategy.Fansub {
		for _, f := range collection.Fansub {
			if f == v {
				fansubPoint += len(strategy.Fansub) - i
			}
		}
	}

	// 计算类型分
	//for i, v := range strategy.Type {
	//	// 类型匹配一旦匹配中，一定会被优先选择
	//	if v == collection.Type {
	//		typePoint = len(strategy.Type) - i
	//		break
	//	}
	//}
	switch collection.Type {
	case data.Full:
		typePoint = 2
	case data.Episode:
		typePoint = 1
	}

	// 计算语言分
	for i, v := range strategy.Language {
		// 尝试匹配语言
		if collection.Language&v != 0 {
			// 找到最符合的语言后就会，如果还有其他语言，则作为扣分项
			if lanPoint == 0 {
				lanPoint = len(strategy.Language) - i
			} else {
				lanPoint--
			}
		}
	}

	// 包含中文
	//if collection.Language&data.GB != 0 {
	//	lanPoint += 2
	//	if
	//}

	// 内嵌字幕部分
	switch collection.SubType {
	case data.Internal:
		subPoint = 1
	case data.External:
		subPoint = 0
	}

	// 最后计算的时候，如果语言不匹配，直接pass掉
	if lanPoint == 0 {
		return 0
	}
	// 如果有内嵌字幕，一定优先内嵌
	point += subPoint * 8
	// 如果类型中有集合类型，那么基本上只会从中选择了
	point += typePoint * 8
	// 接下来的都谁加分项目，画质>汉化组>
	point += qualityPoint*4 + fansubPoint*3 + lanPoint*3
	return point
}
