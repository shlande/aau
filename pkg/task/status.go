package task

import (
	"log"
	"time"
)

type Status uint8

const (
	// Wait 已经完成当前阶段的更新，等待下一步更新
	Wait Status = iota
	// Update 正在等待获取到最新的更新
	Update
	// Download 正在下载
	Download
	// Finish 已经完成更新
	Finish
)

type Action uint8

const (
	UpdateFinish Action = iota
	UpdateFail
	UpdateSkip
	DownloadFinish
	DownloadCancel
	DownloadTimeout
	Terminate
)

func NewPolicy(updateDay time.Weekday) *Policy {
	plc := &Policy{
		Status:    Wait,
		UpdateDay: updateDay,
		Last:      time.Now(),
	}
	plc.Next = plc.nextUpdate()
	return plc
}

// Policy 负责控制更新策略
type Policy struct {
	Status
	// 更新的日期
	UpdateDay time.Weekday
	Last      time.Time
	Next      time.Time
	Logs      []*Log `json:"_"`
}

// 计算下一次应该更新的时间
func (p *Policy) nextUpdate() time.Time {
	if time.Now().Weekday() < p.UpdateDay {
		// 将时间快进到本周到这一天
	}
	// 否则就快进到下一周到这一条
	return time.Time{}
}

// Updated 设置为完成更新
func (p *Policy) Updated(msg string) {
	p.addLog(UpdateFinish, msg)
	p.Status = Download
}

// SkipUpdate 跳过本次更新
func (p *Policy) SkipUpdate(msg string) {
	p.addLog(UpdateSkip, msg)
	p.Next = p.nextUpdate()
	p.Status = Wait
}

func (p *Policy) CancelDownload(msg string) {
	p.addLog(DownloadCancel, msg)
	p.Next = p.nextUpdate()
	p.Status = Wait
}

// CheckTimeout 检查下载是否超时
func (p *Policy) CheckTimeout() bool {
	lastUpdate := p.Logs[len(p.Logs)-1].EmitTime
	now := time.Now()
	// 超过三天则标记超时
	if p.Status == Download && lastUpdate.Add(time.Hour*24*3).Before(now) {
		return true
	}
	return false
}

func (p *Policy) DownloadTimeout(msg string) {
	p.addLog(DownloadTimeout, "")
	p.Status = Wait
}

func (p *Policy) FinishDownload(msg string) {
	p.addLog(DownloadFinish, msg)
	p.Next = p.nextUpdate()
	p.Status = Wait
}

func (p *Policy) Terminate(msg string) {
	p.addLog(Terminate, msg)
	p.Status = Finish
}

// CheckUpdate 判断是否需要更新
func (p *Policy) CheckUpdate() bool {
	// 如果当前确定需要更新
	if p.Status == Update {
		return true
	}
	if p.Status == Wait && p.Next.Before(time.Now()) {
		p.Status = Update
		return true
	}
	return false
}

func (p *Policy) UpdateFail(err error) {
	p.addLog(UpdateFail, err.Error())
	p.Next = p.Next.Add(time.Hour)
	p.Status = Wait
}

func (p *Policy) Log() []*Log {
	return p.Logs
}

func (p *Policy) addLog(action Action, msg string) {
	log.Printf("Policy触发操作：%v %v", action, msg)
	p.Logs = append(p.Logs, &Log{
		Action:   action,
		EmitTime: time.Now(),
		Message:  msg,
	})
}

type Log struct {
	Action
	EmitTime time.Time
	Message  string
}
