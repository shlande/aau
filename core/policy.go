package core

import "time"

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
	FinishUpdate Action = iota
	SkipUpdate
	FinishDownload
	CancelDownload
)

// Policy 负责控制更新策略
type Policy struct {
	Status
	// 更新的日期
	UpdateDay  time.Weekday
	LastUpdate *time.Time
	Logs       []Log
}

// FinishUpdate 设置为完成更新
func (p *Policy) FinishUpdate() {}

// SkipUpdate 跳过本次更新
func (p *Policy) SkipUpdate() {

}

func (p *Policy) CancelDownload() {

}

func (p *Policy) FinishDownload() {

}

func (p *Policy) Terminate() {

}

func (p *Policy) Log() []*Log {
	return nil
}

// NextUpdate 下一次更新时间
func (p *Policy) NextUpdate() *time.Time {
	return &time.Time{}
}

// CheckUpdate 判断是否需要更新
func (p *Policy) CheckUpdate() bool {
	// 如果当前确定需要更新
	if p.Status == Update {
		return true
	}
	// 如果更新时间还是属于一周的范围内
	//if
	//p.LastUpdate.Weekday()
	//if p.LastUpdate.Add(time.Hour * 24 * 7).Before(time.Now()) {
	//
	//}
	return false
}

func (p *Policy) TryUpdate() {

}

type Log struct {
	Action
	EmitTime *time.Time
}
