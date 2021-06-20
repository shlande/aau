package common

import (
	"github.com/shlande/dmhy-rss/pkg/parser"
	"reflect"
	"testing"
)

//func TestFindByKeywords(t *testing.T) {
//	cls, err := FindCollectionsByKeywords(context.Background(), "无职转生")
//	if err != nil {
//		panic(err)
//	}
//	data, _ := json.Marshal(cls)
//	w, _ := os.Create("../test.json")
//	defer w.Close()
//	io.Copy(w, bytes.NewReader(data))
//}

//func TestParse(t *testing.T) {
//	f, err := os.Open("../test.xml")
//	defer f.Close()
//	if err != nil {
//		panic(err)
//	}
//	rcds, err := parse(f)
//	cls := classify(rcds)
//	if err != nil {
//		panic(err)
//	}
//	data, _ := json.Marshal(cls)
//	w, _ := os.Create("../test.json")
//	defer w.Close()
//	io.Copy(w, bytes.NewReader(data))
//}

func TestParseTitle(t *testing.T) {
	input := []struct {
		Title string
		Want  *parser.TitleInfo
	}{
		{
			Title: "[桜都字幕組] 無職轉生～到了異世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][繁體內嵌]",
			Want:  &parser.TitleInfo{Name: "無職轉生～到了異世界就拿出真本事～", Language: parser.BIG5, Quality: parser.P1080, Episode: 11, SubType: parser.Internal},
		}, {
			Title: "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
			Want:  &parser.TitleInfo{Name: "無職轉生～到了異世界就拿出真本事～", Language: parser.GB | parser.BIG5, Quality: parser.P1080, Episode: 11, SubType: parser.Internal},
		}, {
			Title: "[爱恋\\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][10][1080p][MP4][GB][简中]",
			Want:  &parser.TitleInfo{Name: "无职转生～到了异世界就拿出真本事～", Language: parser.GB, Quality: parser.P1080, Episode: 10, SubType: parser.Internal},
		}, {
			Title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][08][x264 1080p][CHT]",
			Want:  &parser.TitleInfo{Name: "無職轉生～到了異世界就拿出真本事～", Language: parser.BIG5, Quality: parser.P1080, Episode: 8, SubType: parser.Internal},
		}, {
			Title: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
			Want:  &parser.TitleInfo{Name: "無職轉生～到了異世界就拿出真本事～", Language: parser.BIG5, Quality: parser.P1080, Episode: 10, SubType: parser.Internal},
		}, {
			Title: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 10 [ViuTV][WEB-DL][1080p][AVC AAC][繁體外掛][MP4+ASSx2](正式版本)",
			Want:  &parser.TitleInfo{Name: "無職轉生，到了異世界就拿出真本事", Language: parser.BIG5, Quality: parser.P1080, Episode: 10, SubType: parser.External},
		}, {
			Title: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
			Want:  &parser.TitleInfo{Name: "無職轉生", Language: parser.BIG5, Quality: parser.P720, Episode: 9, SubType: parser.Internal},
		},
	}
	for _, i := range input {
		got := parseTitle(i.Title)
		if !reflect.DeepEqual(got, i.Want) {
			t.Errorf("want: %v , got : %v , name: %v", i.Want, got, i.Title)
		}
	}
}

//func TestSortSingle(t *testing.T) {
//	items := []*parser.Item{
//		&parser.Item{
//			Fansub:   []string{"nekomoekissaten"},
//			Title:    "【喵萌奶茶屋】★01月新番★[无职转生/Mushoku Tensei][11][1080p][简体][招募翻译校对]",
//			RawCategory: parser.Animate,
//			Detail: &parser.Detail{
//				Name:     "",
//				Language: parser.GB,
//				Quality:  parser.P1080,
//				Episode:  11,
//				SubType:  parser.Internal,
//			},
//		},
//	}
//	items = sort(items, &Option{
//		Episode:  11,
//		Fansub:   []string{"nekomoekissaten"},
//		Quality:  parser.P1080,
//		Language: parser.GB,
//	})
//	if len(items) == 0 {
//		panic("sort error")
//	}
//}
//
//func TestSort(t *testing.T) {
//	f, err := os.Open("../test.xml")
//	defer f.Close()
//	if err != nil {
//		panic(err)
//	}
//	items, err := parse(f)
//	items = sort(items, &Option{
//		Episode:  11,
//		Fansub:   []string{"nekomoekissaten"},
//		Quality:  parser.P1080,
//		Language: parser.GB,
//	})
//	if len(items) != 1 {
//		panic("sort error")
//	}
//}
//
//func TestClassify(t *testing.T) {
//	f, err := os.Open("../test.xml")
//	defer f.Close()
//	if err != nil {
//		panic(err)
//	}
//	items, err := parse(f)
//	fmt.Println(items)
//}

func Test_parseName(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "aaaa",
			args: `[神枫字幕组] [1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][02][1080p][MP4][GB][简体中文]`,
			want: "无职转生～到了异世界就拿出真本事～",
		},
		{
			name: "en",
			args: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 11 [ViuTV][WEB-RIP][CHT][SRTx2][1080p][AVC AAC][MKV](先行版本)",
			want: `無職轉生，到了異世界就拿出真本事`,
		},
		{
			name: "",
			args: `[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]`,
			want: "無職轉生～到了異世界就拿出真本事～",
		},
		{
			name: "en-only",
			args: "[爱恋\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][11][1080p][MP4][GB][简中]",
			want: `无职转生～到了异世界就拿出真本事～`,
		},
		{
			name: "neko",
			args: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
			want: "無職轉生",
		},
		{
			name: "cc",
			args: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
			want: "無職轉生～到了異世界就拿出真本事～",
		},
		{
			name: "UHA-wings",
			args: `【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][09][x264 1080p][CHT]`,
			want: "無職轉生～到了異世界就拿出真本事～",
		},
		{
			name: "cn",
			args: "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][720p][简体内嵌]",
			want: `无职转生～到了异世界就拿出真本事～`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseName(tt.args, false); got != tt.want {
				t.Errorf("parseFansub() = %v , want %v", got, tt.want)
			}
		})
	}
}
