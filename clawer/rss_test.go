package clawer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFindByKeywords(t *testing.T) {
	FindCollectionsByKeywords(context.Background(), "无职转生")
}

func TestParse(t *testing.T) {
	f, err := os.Open("../test.xml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	rcds, err := parse(f)
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(rcds)
	w, _ := os.Create("../test.json")
	defer w.Close()
	io.Copy(w, bytes.NewReader(data))
}

func TestParseTitle(t *testing.T) {
	input := []struct {
		Title string
		Want  *Detail
	}{
		{
			Title: "[桜都字幕組] 無職轉生～到了異世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][繁體內嵌]",
			Want:  &Detail{Language: BIG5, Quality: P1080, Episode: 11, SubType: Internal},
		}, {
			Title: "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
			Want:  &Detail{Language: GB | BIG5, Quality: P1080, Episode: 11, SubType: Internal},
		}, {
			Title: "[爱恋\\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][10][1080p][MP4][GB][简中]",
			Want:  &Detail{Language: GB, Quality: P1080, Episode: 10, SubType: Internal},
		}, {
			Title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][08][x264 1080p][CHT]",
			Want:  &Detail{Language: BIG5, Quality: P1080, Episode: 8, SubType: Internal},
		}, {
			Title: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
			Want:  &Detail{Language: BIG5, Quality: P1080, Episode: 10, SubType: Internal},
		}, {
			Title: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 10 [ViuTV][WEB-DL][1080p][AVC AAC][繁體外掛][MP4+ASSx2](正式版本)",
			Want:  &Detail{Language: BIG5, Quality: P1080, Episode: 10, SubType: External},
		}, {
			Title: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
			Want:  &Detail{Language: BIG5, Quality: P720, Episode: 9, SubType: Internal},
		},
	}
	for _, i := range input {
		got := parseTitle(i.Title)
		if !reflect.DeepEqual(got, i.Want) {
			t.Errorf("want: %v , got : %v , name: %v", i.Want, got, i.Title)
		}
	}
}

func TestSortSingle(t *testing.T) {
	items := []*Item{
		&Item{
			Fansub:   []string{"nekomoekissaten"},
			Title:    "【喵萌奶茶屋】★01月新番★[无职转生/Mushoku Tensei][11][1080p][简体][招募翻译校对]",
			Category: Animate,
			Detail: &Detail{
				Name:     "",
				Language: GB,
				Quality:  P1080,
				Episode:  11,
				SubType:  Internal,
			},
		},
	}
	items = sort(items, &Option{
		Episode:  11,
		Fansub:   []string{"nekomoekissaten"},
		Quality:  P1080,
		Language: GB,
	})
	if len(items) == 0 {
		panic("sort error")
	}
}

func TestSort(t *testing.T) {
	f, err := os.Open("../test.xml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	items, err := parse(f)
	items = sort(items, &Option{
		Episode:  11,
		Fansub:   []string{"nekomoekissaten"},
		Quality:  P1080,
		Language: GB,
	})
	if len(items) != 1 {
		panic("sort error")
	}
}

func TestClassify(t *testing.T) {
	f, err := os.Open("../test.xml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	items, err := parse(f)
	fmt.Println(items)
}
