package parser

import (
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

//func TestParseTitle(t *testing.T) {
//	ps := New()
//	input := []struct {
//		Title string
//		Want  *Result
//	}{
//		//{
//		//	// FIXME：修复这个bug
//		//	Title: "[時雨初空] 剃掉鬍子。然後撿了個女高中生。 12 繁體 MP4 720p",
//		//	Want:  &Result{Name: "剃掉鬍子。然後撿了個女高中生。", Language: GB, Quality: P720, Episode: 12, SubType: Internal},
//		//},
//		{
//			Title: "[桜都字幕組] 無職轉生～到了異世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][繁體內嵌]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 11, SubType: data.Internal},
//		}, {
//			Title: "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.GB | data.BIG5, Quality: data.P1080, Episode: 11, SubType: data.Internal},
//		}, {
//			Title: "[爱恋\\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][10][1080p][MP4][GB][简中]",
//			Want:  &Result{Name: "无职转生～到了异世界就拿出真本事～", Language: data.GB, Quality: data.P1080, Episode: 10, SubType: data.Internal},
//		}, {
//			Title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][08][x264 1080p][CHT]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 8, SubType: data.Internal},
//		}, {
//			Title: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 10, SubType: data.Internal},
//		}, {
//			Title: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 10 [ViuTV][WEB-DL][1080p][AVC AAC][繁體外掛][MP4+ASSx2](正式版本)",
//			Want:  &Result{Name: "無職轉生，到了異世界就拿出真本事", Language: data.BIG5, Quality: data.P1080, Episode: 10, SubType: data.External},
//		}, {
//			Title: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
//			Want:  &Result{Name: "無職轉生", Language: data.BIG5, Quality: data.P720, Episode: 9, SubType: data.Internal},
//		},
//	}
//	for _, i := range input {
//		got, err := ps.ParseName(i.Title)
//		if err != nil {
//			t.Error(err)
//		}
//		if !reflect.DeepEqual(got, i.Want) {
//			t.Errorf("want: %v , got : %v , name: %v", i.Want, got, i.Title)
//		}
//	}
//}

// 22519	     52476 ns/op df7f0998
// 32079	     35803 ns/op
//func BenchmarkParse_ParseTitle(b *testing.B) {
//	input := []struct {
//		Title string
//		Want  *Result
//	}{
//		//{
//		//	// FIXME：修复这个bug
//		//	Title: "[時雨初空] 剃掉鬍子。然後撿了個女高中生。 12 繁體 MP4 720p",
//		//	Want:  &Result{Name: "剃掉鬍子。然後撿了個女高中生。", Language: GB, Quality: P720, Episode: 12, SubType: Internal},
//		//},
//		{
//			Title: "[桜都字幕組] 無職轉生～到了異世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][繁體內嵌]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 11, SubType: data.Internal},
//		}, {
//			Title: "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.GB | data.BIG5, Quality: data.P1080, Episode: 11, SubType: data.Internal},
//		}, {
//			Title: "[爱恋\\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][10][1080p][MP4][GB][简中]",
//			Want:  &Result{Name: "无职转生～到了异世界就拿出真本事～", Language: data.GB, Quality: data.P1080, Episode: 10, SubType: data.Internal},
//		}, {
//			Title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][08][x264 1080p][CHT]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 8, SubType: data.Internal},
//		}, {
//			Title: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
//			Want:  &Result{Name: "無職轉生～到了異世界就拿出真本事～", Language: data.BIG5, Quality: data.P1080, Episode: 10, SubType: data.Internal},
//		}, {
//			Title: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 10 [ViuTV][WEB-DL][1080p][AVC AAC][繁體外掛][MP4+ASSx2](正式版本)",
//			Want:  &Result{Name: "無職轉生，到了異世界就拿出真本事", Language: data.BIG5, Quality: data.P1080, Episode: 10, SubType: data.External},
//		}, {
//			Title: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
//			Want:  &Result{Name: "無職轉生", Language: data.BIG5, Quality: data.P720, Episode: 9, SubType: data.Internal},
//		},
//	}
//	ps := New()
//	b.StartTimer()
//	for i := 0; i < b.N; i++ {
//		got, err := ps.ParseName(input[0].Title)
//		if err != nil {
//			b.Error(err)
//		}
//		if !reflect.DeepEqual(got, input[0].Want) {
//			b.Errorf("want: %v , got : %v , name: %v", input[0].Want, got, input[0].Title)
//		}
//	}
//}

func ParseFunsub(t *testing.T) {

}

//func TestSortSingle(t *testing.T) {
//	items := []*Item{
//		&Item{
//			Fansub:   []string{"nekomoekissaten"},
//			Title:    "【喵萌奶茶屋】★01月新番★[无职转生/Mushoku Tensei][11][1080p][简体][招募翻译校对]",
//			RawCategory: Animate,
//			Detail: &Detail{
//				Name:     "",
//				Language: GB,
//				Quality:  P1080,
//				Episode:  11,
//				SubType:  Internal,
//			},
//		},
//	}
//	items = sort(items, &Option{
//		Episode:  11,
//		Fansub:   []string{"nekomoekissaten"},
//		Quality:  P1080,
//		Language: GB,
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
//		Quality:  P1080,
//		Language: GB,
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

func Test_parseFansub(t *testing.T) {
	tests := []struct {
		name  string
		title string
		want  []string
	}{
		{
			name:  "[]",
			title: "[桜都字幕組] 無職轉生～到了異世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][繁體內嵌]",
			want:  []string{"桜都字幕組"},
		}, {
			name:  "【】",
			title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][09][x264 1080p][CHT]",
			want:  []string{"悠哈璃羽字幕社"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFansub(tt.title); reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFansub() = %v, want %v", got, tt.want)
			}
		})
	}
}
