package bolt

import (
	"context"
	"encoding/json"
	worker3 "github.com/shlande/dmhy-rss/pkg/controller/manager"
	"reflect"
	"testing"
	"time"
)

var testJsonCollection = `{
    "Name": "无职转生 ～到了异世界就拿出真本事～",
    "Fansub": [
      "九十九朔夜"
    ],
    "Quality": 2,
    "Type": 2,
    "SubType": 1,
    "Language": 6,
    "Latest": 11,
    "LastUpdate": "2021-04-26T22:47:30.025383+08:00",
    "Items": [
	  {
        "Fansub": [
          "九十九朔夜"
        ],
        "Title": "[NC-Raws] 无职转生 ～到了异世界就拿出真本事～ / Mushoku Tensei - 08 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "Type": 2,
        "Name": "无职转生 ～到了异世界就拿出真本事～",
        "Language": 6,
        "Quality": 2,
        "Episode": 8,
        "SubType": 1,
        "MagnetUrl": "magnet:?xt=urn:btih:POGUYM3EZENERRPQXVF5B66U4EIFTM6S\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "Link": "http://share.dmhy.org/topics/view/561704_NC-Raws_Mushoku_Tensei_-_08_WEB-DL_1080p_AVC_AAC_CHS_CHT_SRT_MKV.html",
        "CreateTime": "2021-03-04T06:52:50Z"
      },
	  {
        "Fansub": [
          "九十九朔夜"
        ],
        "Title": "[NC-Raws] 无职转生 ～到了异世界就拿出真本事～ / Mushoku Tensei - 09 [WEB-DL][1080p][AVC AAC][CHS_CHT_TH_SRT][MKV]",
        "Type": 2,
        "Name": "无职转生 ～到了异世界就拿出真本事～",
        "Language": 6,
        "Quality": 2,
        "Episode": 9,
        "SubType": 1,
        "MagnetUrl": "magnet:?xt=urn:btih:DJOW4PJXPJIPZ2D2JHD2FVGHQBVVX25I\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "Link": "http://share.dmhy.org/topics/view/562055_NC-Raws_Mushoku_Tensei_-_09_WEB-DL_1080p_AVC_AAC_CHS_CHT_TH_SRT_MKV.html",
        "CreateTime": "2021-03-07T16:33:35Z"
      },
	  {
        "Fansub": [
          "九十九朔夜"
        ],
        "Title": "[NC-Raws] 无职转生 ～到了异世界就拿出真本事～ / Mushoku Tensei - 10 [WEB-DL][1080p][AVC AAC][CHS_CHT_TH_SRT][MKV]",
        "Type": 2,
        "Name": "无职转生 ～到了异世界就拿出真本事～",
        "Language": 6,
        "Quality": 2,
        "Episode": 10,
        "SubType": 1,
        "MagnetUrl": "magnet:?xt=urn:btih:GE7ETVTKRWOIKN3FW6XKRYGMDDFPXMN3\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce",
        "Link": "http://share.dmhy.org/topics/view/562693_NC-Raws_Mushoku_Tensei_-_10_WEB-DL_1080p_AVC_AAC_CHS_CHT_TH_SRT_MKV.html",
        "CreateTime": "2021-03-14T16:31:56Z"
      },
      {
        "Fansub": [
          "九十九朔夜"
        ],
        "Title": "[NC-Raws] 无职转生 ～到了异世界就拿出真本事～ / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_TH_SRT][MKV]",
        "Type": 2,
        "Name": "无职转生 ～到了异世界就拿出真本事～",
        "Language": 6,
        "Quality": 2,
        "Episode": 11,
        "SubType": 1,
        "MagnetUrl": "magnet:?xt=urn:btih:SPASVCVM4LL7RQ2WRP2RWGGR2UJC3KDQ\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "Link": "http://share.dmhy.org/topics/view/563334_NC-Raws_Mushoku_Tensei_-_11_WEB-DL_1080p_AVC_AAC_CHS_CHT_TH_SRT_MKV.html",
        "CreateTime": "2021-03-21T16:31:33Z"
      }
    ]
  }`
var testJsonFullSession = `{
    "Name": "无职转生～到了异世界就拿出真本事～",
    "Fansub": [
      "sakurato"
    ],
    "Quality": 2,
    "Type": 1,
    "SubType": 1,
    "Language": 2,
    "Latest": 11,
    "LastUpdate": "2021-04-26T22:47:30.025333+08:00",
    "Items": [
      {
        "Fansub": [
          "sakurato"
        ],
        "Title": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [01-11 Fin][1080p@60FPS][简体内嵌]",
        "Type": 1,
        "Name": "无职转生～到了异世界就拿出真本事～",
        "Language": 2,
        "Quality": 2,
        "Episode": 0,
        "SubType": 1,
        "MagnetUrl": "magnet:?xt=urn:btih:JZGBGXPG7AY4TLI4LP527AE5DBV7DXVI\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "Link": "http://share.dmhy.org/topics/view/564511_Mushoku_Tensei_Isekai_Ittara_Honki_Dasu_01-11_Fin_1080p_60FPS.html",
        "CreateTime": "2021-04-03T04:36:02Z"
      }
    ]
  }`

func TestBolt_Save(t *testing.T) {
	var cl = &classify.Collection{}
	var clFull = &classify.Collection{}
	err := json.Unmarshal([]byte(testJsonFullSession), clFull)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(testJsonCollection), cl)
	if err != nil {
		panic(err)
	}
	kv, err := New("./test.db")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name       string
		collection []*classify.Collection
		wantErr    bool
	}{
		{name: "ep", collection: []*classify.Collection{cl}},
		{name: "full", collection: []*classify.Collection{clFull}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kv.Save(tt.collection...); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			c, err := kv.Get(tt.collection[0].Id())
			if err != nil {
				t.Errorf("无法获取内容：%v", err)
			}
			if !reflect.DeepEqual(c, tt.collection[0]) {
				t.Errorf("not equal, %v %v", c, tt.collection[0])
			}
		})
	}
}

func TestBolt_SaveWorker(t *testing.T) {
	var cl = &classify.Collection{}
	err := json.Unmarshal([]byte(testJsonCollection), cl)
	if err != nil {
		panic(err)
	}
	worker := worker3.NewWorker(cl, time.Wednesday, nil, nil, nil)
	kv, err := New("./test.db")
	if err != nil {
		panic(err)
	}
	// 先保存collection的内容
	err = kv.Save(cl)
	if err != nil {
		panic(err)
	}
	err = kv.SaveWorker(worker)
	if err != nil {
		panic(err)
	}
	builder, err := kv.GetWorker(worker.Id)
	if err != nil {
		panic(err)
	}
	wok2 := builder.Recover(context.Background(), nil, nil, nil)

	// 这里不能通过是因为结构体中有chan，其余是正常的
	if !reflect.DeepEqual(worker, wok2) {
		panic("not equal")
	}
}

func TestBolt_ListWorker(t *testing.T) {
	var cl = &classify.Collection{}
	err := json.Unmarshal([]byte(testJsonCollection), cl)
	if err != nil {
		panic(err)
	}
	worker := worker3.NewWorker(cl, time.Wednesday, nil, nil, nil)
	kv, err := New("./test.db")
	if err != nil {
		panic(err)
	}
	// 先保存collection的内容
	err = kv.Save(cl)
	if err != nil {
		panic(err)
	}
	err = kv.SaveWorker(worker)
	if err != nil {
		panic(err)
	}
	wl, err := kv.ListWorker()
	if len(wl) == 0 {
		panic("len0")
	}
	err = kv.Save(cl)
	if err != nil {
		panic(err)
	}
	err = kv.SaveWorker(worker)
	if err != nil {
		panic(err)
	}
	wl, err = kv.ListWorker()
	if len(wl) == 0 {
		panic("len0")
	}
}
