package clawer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestFindByKeywords(t *testing.T) {
	FindCollectionsByKeywords(context.Background(), "无职转生")
}

func TestParse(t *testing.T) {
	f, err := os.Open("./test.xml")
	defer f.Close()
	if err != nil {
		panic(err)
	}
	rcds, err := parse(f)
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(rcds)
	w, _ := os.Create("./test.json")
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
			Want:  &Detail{Language: GB | BIG, Quality: P1080, Episode: 11, SubType: Internal},
		}, {
			Title: "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
			Want:  &Detail{Language: GB | BIG, Quality: P1080, Episode: 11, SubType: Internal},
		}, {
			Title: "[爱恋\\u0026漫猫字幕组][1月新番][无职转生～到了异世界就拿出真本事～][Mushoku Tensei Isekai Ittara Honki Dasu][10][1080p][MP4][GB][简中]",
			Want:  &Detail{Language: GB, Quality: P1080, Episode: 10, SubType: Internal},
		}, {
			Title: "【悠哈璃羽字幕社】[無職轉生～到了異世界就拿出真本事～_Mushoku Tensei][08][x264 1080p][CHT]",
			Want:  &Detail{Language: GB, Quality: P1080, Episode: 10, SubType: Internal},
		}, {
			Title: "​[c.c動漫][1月新番][無職轉生～到了異世界就拿出真本事～][10][BIG5][1080P][MP4]",
			Want:  nil,
		}, {
			Title: "[Skymoon-Raws] 無職轉生，到了異世界就拿出真本事 / Mushoku Tensei - 10 [ViuTV][WEB-DL][1080p][AVC AAC][繁體外掛][MP4+ASSx2](正式版本)",
			Want:  nil,
		}, {
			Title: "【喵萌奶茶屋】★01月新番★[無職轉生/Mushoku Tensei][09][720p][繁體][招募翻譯校對]",
			Want:  nil,
		},
	}
	for _, i := range input {
		got := parseTitle(i.Title)
		if !reflect.DeepEqual(got, i.Want) {
			t.Errorf("want: %v , got : %v", i.Want, got)
		}
	}
}
