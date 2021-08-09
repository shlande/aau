package bolt

import (
	"encoding/json"
	"errors"
	"github.com/shlande/dmhy-rss/pkg/controller/mission"
	"github.com/shlande/dmhy-rss/pkg/controller/store"
	"github.com/shlande/dmhy-rss/pkg/data"
	"os"
	"reflect"
	"testing"
	"time"
)

var testJsonCollection = ` {
    "Id": "",
    "Name": "無職転生～異世界行ったら本気だす～",
    "Translated": "无职转生～到了异世界就拿出真本事～",
    "Summary": "",
    "AirDate": "0001-01-01T00:00:00Z",
    "AirWeekday": 0,
    "AirTime": 0,
    "TotalEpisodes": 0,
    "Category": "",
    "Fansub": [
      "桜都字幕组"
    ],
    "Quality": 2,
    "Type": 2,
    "Language": 2,
    "SubType": 1,
    "Latest": 11,
    "LastUpdate": "2021-03-22T16:19:20Z",
    "Items": [
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [01][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-01-11T07:28:31Z",
        "MagnetUrl": "magnet:?xt=urn:btih:WU55KUGTJOYHIHWXMH7A36KZZADJE3JI\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.xyz%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce",
        "TorrentUrl": "",
        "Episode": 1,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [02][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-01-18T07:27:00Z",
        "MagnetUrl": "magnet:?xt=urn:btih:5VK2SVQMNINSZJ2T52KYDLBJTSJGR4QT\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.xyz%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce",
        "TorrentUrl": "",
        "Episode": 2,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [03][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-01-26T07:45:04Z",
        "MagnetUrl": "magnet:?xt=urn:btih:JF4LE22ZTA6R76PYIZS4KANNOQD7KCZH\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.xyz%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce",
        "TorrentUrl": "",
        "Episode": 3,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [04][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-02-02T00:32:25Z",
        "MagnetUrl": "magnet:?xt=urn:btih:LCIX3YCS5SACRFDYIX2EOE7KWMFGR7J5\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.xyz%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce",
        "TorrentUrl": "",
        "Episode": 4,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [05][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-02-09T00:46:42Z",
        "MagnetUrl": "magnet:?xt=urn:btih:LECJTWPJZCTYYPIZAJZIZ563AQOBZV56\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.xyz%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce",
        "TorrentUrl": "",
        "Episode": 5,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [06][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-02-16T03:21:57Z",
        "MagnetUrl": "magnet:?xt=urn:btih:D3KGIFMRWJ6UEO65JP4NCYWGIXSA7X3M\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 6,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [07][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-02-22T23:17:17Z",
        "MagnetUrl": "magnet:?xt=urn:btih:4L5ZOX2MZLCBBZIT2JPF5JZHAA44W6S2\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 7,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [08][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-03-02T00:24:44Z",
        "MagnetUrl": "magnet:?xt=urn:btih:XMECL2CUEHUZ6V44EJQSDGD3PDYFZ52H\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 8,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [09][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-03-09T05:07:08Z",
        "MagnetUrl": "magnet:?xt=urn:btih:VCOKEIQDJYOBE6LM7ORWE4BPS2LYSBAN\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 9,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [10][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-03-16T05:04:44Z",
        "MagnetUrl": "magnet:?xt=urn:btih:UGOGQNB6ORNBJA22TEI3EOC2MIBU6LGR\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 10,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      },
      {
        "Name": "[桜都字幕组] 无职转生～到了异世界就拿出真本事～ / Mushoku Tensei Isekai Ittara Honki Dasu [11][1080p@60FPS][简体内嵌]",
        "CreateTime": "2021-03-22T16:19:20Z",
        "MagnetUrl": "magnet:?xt=urn:btih:6KSX3EGQ5JCBDVPGGGTT44GV2EE3SBED\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftracker.sakurato.art%3A23334%2Fannounce\u0026tr=http%3A%2F%2Ftracker.sakurato.art%3A23333%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Fsukebei.tracker.wf%3A8888%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=https%3A%2F%2Fopen.acgnxtracker.com%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 11,
        "Fansub": [
          "桜都字幕组"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 2,
        "SubType": 1
      }
    ]
  }`
var testJsonMissingCollection = `{
    "Id": "",
    "Name": "無職転生～異世界行ったら本気だす～",
    "Translated": "无职转生～到了异世界就拿出真本事～",
    "Summary": "",
    "AirDate": "0001-01-01T00:00:00Z",
    "AirWeekday": 0,
    "AirTime": 0,
    "TotalEpisodes": 0,
    "Category": "",
    "Fansub": [
      "NC-Raws"
    ],
    "Quality": 2,
    "Type": 2,
    "Language": 6,
    "SubType": 1,
    "Latest": 11,
    "LastUpdate": "2021-03-21T16:32:00Z",
    "Items": [
      {
        "Name": "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 06 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-02-14T16:34:40Z",
        "MagnetUrl": "magnet:?xt=urn:btih:JX3YSY7L3N2HYGSV6ORZKP7WVY5D6G5F\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 6,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      },
      {
        "Name": "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 07 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-02-21T16:35:22Z",
        "MagnetUrl": "magnet:?xt=urn:btih:2XFJCOXUUQB5TVSU3SPYPUBCXVTG7PUP\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 7,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      },
      {
        "Name": "[NC-Raws] 无职转生 ～到了异世界就拿出真本事～ / Mushoku Tensei - 08 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-03-04T06:52:50Z",
        "MagnetUrl": "magnet:?xt=urn:btih:POGUYM3EZENERRPQXVF5B66U4EIFTM6S\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 8,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      },
      {
        "Name": "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 09 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-03-07T16:35:39Z",
        "MagnetUrl": "magnet:?xt=urn:btih:5IQYRG4MEQCAW6OAMZYKHYPWWT2OMH7L\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 9,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      },
      {
        "Name": "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 10 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-03-14T16:32:05Z",
        "MagnetUrl": "magnet:?xt=urn:btih:ROXAFBAWQFO64TSIREXICWOX4GRM7BLJ\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce",
        "TorrentUrl": "",
        "Episode": 10,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      },
      {
        "Name": "[NC-Raws] 無職轉生～到了異世界就拿出真本事～（僅限港澳台地區） / Mushoku Tensei - 11 [WEB-DL][1080p][AVC AAC][CHS_CHT_SRT][MKV]",
        "CreateTime": "2021-03-21T16:32:00Z",
        "MagnetUrl": "magnet:?xt=urn:btih:M6ITECXM6GR33HLTHJGLOR5JE4CPDRWM\u0026dn=\u0026tr=http%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=udp%3A%2F%2F104.238.198.186%3A8000%2Fannounce\u0026tr=http%3A%2F%2Ftracker.openbittorrent.com%3A80%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=http%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftracker.publicbt.com%3A80%2Fannounce\u0026tr=http%3A%2F%2Ftracker.prq.to%2Fannounce\u0026tr=http%3A%2F%2Fopen.acgtracker.com%3A1096%2Fannounce\u0026tr=https%3A%2F%2Ft-115.rhcloud.com%2Fonly_for_ylbud\u0026tr=http%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=http%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker1.itzmx.com%3A8080%2Fannounce\u0026tr=udp%3A%2F%2Ftracker2.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker3.itzmx.com%3A6961%2Fannounce\u0026tr=udp%3A%2F%2Ftracker4.itzmx.com%3A2710%2Fannounce\u0026tr=http%3A%2F%2Ftr.bangumi.moe%3A6969%2Fannounce\u0026tr=http%3A%2F%2Ft.nyaatracker.com%2Fannounce\u0026tr=http%3A%2F%2Fopen.nyaatorrents.info%3A6544%2Fannounce\u0026tr=http%3A%2F%2Ft2.popgo.org%3A7456%2Fannonce\u0026tr=http%3A%2F%2Fshare.camoe.cn%3A8080%2Fannounce\u0026tr=http%3A%2F%2Fopentracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Ftracker.acgnx.se%2Fannounce\u0026tr=http%3A%2F%2Fnyaa.tracker.wf%3A7777%2Fannounce\u0026tr=https%3A%2F%2Ftr.bangumi.moe%3A9696%2Fannounce\u0026tr=http%3A%2F%2Ft.acg.rip%3A6699%2Fannounce",
        "TorrentUrl": "",
        "Episode": 11,
        "Fansub": [
          "NC-Raws"
        ],
        "Quality": 2,
        "Type": 2,
        "Language": 6,
        "SubType": 1
      }
    ]
  }`

var cl = &data.Collection{}
var clFull = &data.Collection{}

var testAnm = &data.Animation{
	Id:         "sfesgal",
	Name:       "無職転生～異世界行ったら本気だす～",
	Translated: "无职转生～到了异世界就拿出真本事～",
}

var ms1 *mission.Mission
var ms2 *mission.Mission

var log11 *mission.Log
var log12 *mission.Log

var log21 *mission.Log
var log22 *mission.Log

func load() {
	cl = &data.Collection{}
	clFull = &data.Collection{}
	err := json.Unmarshal([]byte(testJsonCollection), clFull)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(testJsonMissingCollection), cl)
	if err != nil {
		panic(err)
	}

	cl.Animation = testAnm
	clFull.Animation = testAnm

	ms1 = mission.NewMission(testAnm, cl.Metadata)
	ms2 = mission.NewMission(testAnm, clFull.Metadata)

	// ms1 模拟进行一次更新，产生如日志
	log11 = ms1.Next(errors.New("测试内容，产生日志用"))
	// 因为id是时间戳，必须要有一秒才能产生不同的id
	time.Sleep(time.Second * 2)
	log12 = ms1.Next(errors.New("测试内容，产生日志用"))
	time.Sleep(time.Second * 2)
	log21 = ms2.Next([]*data.Source{})
}

func seed(p store.Interface) {
	err := p.Collection().Save(cl)
	if err != nil {
		panic(err)
	}
	err = p.Collection().Save(clFull)
	if err != nil {
		panic(err)
	}
	err = p.Mission().Save(ms1)
	if err != nil {
		panic(err)
	}
	err = p.Mission().Save(ms2)
	if err != nil {
		panic(err)
	}
	err = p.Log().Save(ms1.Id(), log11)
	if err != nil {
		panic(err)
	}
	err = p.Log().Save(ms1.Id(), log12)
	if err != nil {
		panic(err)
	}
	err = p.Log().Save(ms2.Id(), log21)
	if err != nil {
		panic(err)
	}
}

func TestBolt_Save(t *testing.T) {
	load()

	kv := New("./test.db")
	defer func() {
		kv.Close()
		os.Remove("./test.db")
	}()

	tests := []struct {
		name string
		*data.Collection
		wantErr bool
	}{
		{name: "ep", Collection: cl, wantErr: false},
		{name: "full", Collection: clFull},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := kv.Collection().Save(tt.Collection); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(1)
			c, err := kv.Collection().Get(tt.Collection.Id())
			if err != nil {
				t.Errorf("无法获取内容：%v", err)
			}
			if !reflect.DeepEqual(c, tt.Collection) {
				t.Errorf("not equal, %v %v", c, tt.Collection)
			}
		})
	}
}

func TestBolt_SaveMission(t *testing.T) {
	load()

	kv := New("./test.db")
	defer func() {
		kv.Close()
		os.Remove("./test.db")
	}()
	// 先保存collection的内容
	err := kv.Collection().Save(cl)
	if err != nil {
		panic(err)
	}
	err = kv.Mission().Save(ms1)
	if err != nil {
		panic(err)
	}
	ms, err := kv.Mission().Get(cl.Id())
	if err != nil {
		panic(err)
	}
	// 这里不能通过是因为结构体中有chan，其余是正常的
	if !reflect.DeepEqual(ms1, ms) {
		panic("not equal")
	}
}

func TestBolt_ListLog(t *testing.T) {
	load()

	kv := New("./test.db")
	defer func() {
		kv.Close()
		os.Remove("./test.db")
	}()

	seed(kv)

	wl, err := kv.Log().GetAll(ms1.Id())
	if err != nil {
		panic(err)
	}
	for _, v := range wl {
		// TODO:解决精度丢失的问题，除此之外没有问题
		if reflect.DeepEqual(v, log11) || reflect.DeepEqual(v, log12) {
			continue
		}
		panic("not equal")
	}
}

func TestBolt_Pin(t *testing.T) {
	load()

	kv := New("./test.db")
	defer func() {
		kv.Close()
		os.Remove("./test.db")
	}()

	seed(kv)

	res, err := kv.Pin().IsPin(testAnm)
	if err != nil && res {
		panic(err)
	}
	res, err = kv.Pin().IsFinish(testAnm)
	if err != nil && res {
		panic(err)
	}

	kv.Pin().Pin(testAnm)

	res, err = kv.Pin().IsPin(testAnm)
	if err != nil && !res {
		panic(err)
	}
	res, err = kv.Pin().IsFinish(testAnm)
	if err != nil && res {
		panic(err)
	}

	kv.Pin().Unpin(testAnm)
	res, err = kv.Pin().IsPin(testAnm)
	if err != nil && res {
		panic(err)
	}
	res, err = kv.Pin().IsFinish(testAnm)
	if err != nil && res {
		panic(err)
	}

	kv.Pin().Pin(testAnm)
	kv.Pin().Finish(testAnm)

	res, err = kv.Pin().IsPin(testAnm)
	if err != nil && !res {
		panic(err)
	}
	res, err = kv.Pin().IsFinish(testAnm)
	if err != nil && !res {
		panic(err)
	}

	err = kv.Pin().Unpin(testAnm)
	if err == nil {
		panic(err)
	}

}
