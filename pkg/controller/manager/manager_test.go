package manager

//func TestWorker(t *testing.T) {
//	ctx := context.Background()
//
//	worker := (data.NewCollection(&data.Animation{
//		Name:          "無職転生～異世界行ったら本気だす～",
//		Translated:    "无职转生～到了异世界就拿出真本事～",
//		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
//		AirWeekday:    time.Sunday,
//		AirTime:       time.Hour * 3,
//		TotalEpisodes: 11,
//		Category:      "tv",
//	}, data.Metadata{
//		Fansub:   []string{"桜都字幕組"},
//		Quality:  data.P1080,
//		Type:     data.Episode,
//		Language: 4,
//		SubType:  data.Internal,
//	}), tools.CollectionProvider{
//		Parser:   parser.New(),
//		Provider: dmhy.NewProvider(),
//	}, nil)
//
//}
//
//func TestOnceCollect(t *testing.T) {
//	ctx := context.Background()
//
//	worker := NewWorker(data.NewCollection(&data.Animation{
//		Name:          "無職転生～異世界行ったら本気だす～",
//		Translated:    "无职转生～到了异世界就拿出真本事～",
//		AirDate:       time.Now().Truncate(time.Hour * 24 * 30 * 2),
//		AirWeekday:    time.Sunday,
//		AirTime:       time.Hour * 3,
//		TotalEpisodes: 11,
//		Category:      "tv",
//	}, data.Metadata{
//		Fansub:   []string{"桜都字幕組"},
//		Quality:  data.P1080,
//		Type:     data.Episode,
//		Language: 4,
//		SubType:  data.Internal,
//	}), tools.CollectionProvider{
//		Parser:   parser.New(),
//		Provider: dmhy.NewProvider(),
//	}, nil)
//
//	var m Machine = &waiting{Misson: worker, Timer: time.NewTimer(0)}
//	// 第一次，应该跳转到waiting状态
//	m = m.Do(ctx)
//	if m.Status() != Updating {
//		panic("0")
//	}
//	// 准备更新
//	m = m.Do(ctx)
//	if m.Status() != Finish {
//		panic("1")
//	}
//}
