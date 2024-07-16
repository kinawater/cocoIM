package domain

import (
	"sync/atomic"
	"unsafe"
)

// 每一个endport是一个独立服务
type Endport struct {
	IP          string       `json:"ip"`
	Port        string       `json:"port"`
	ActiveScore float64      `json:"-"`
	StaticScore float64      `json:"-"`
	Stats       *Stat        `json:"-"`
	Window      *stateWindow `json:"-"`
}

func NewEndport(ip, port string) *Endport {
	ed := &Endport{
		IP:   ip,
		Port: port,
	}
	// 每次创建一个Endport都会创建一个对应的观测窗口window
	ed.Window = newStateWindow()
	ed.Stats = ed.Window.getStat()
	go func() {
		// 查看窗口的统计数据是否有新值，如果不出错的情况下应该一秒一个
		for stat := range ed.Window.statChan {
			// 拿到新值后重新计算平均值
			ed.Window.updateSumStat(stat)
			newStat := ed.Window.getStat()
			// 把新算好的值直接替换掉旧值
			atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(ed.Stats)), unsafe.Pointer(newStat))
		}
	}()
	return ed
}
func (ed *Endport) UpdateStat(s *Stat) {
	ed.Window.statChan <- s
}
func (ed *Endport) CalculateScore(ctx *IpConfContext) {
	// 通常Stats不会为nil，特别是静态分一开始就在注册的时候就有了，如果一直为nil，应该报错
	// 报错后应该通知注册中心，调度器应该发出警报，检查服务的情况

	// 但是也有可能其他原因，在某一次计算的时候发现stat为nil，此时不需要处理，沿用上一次的计算结果即可
	if ed.Stats != nil {
		ed.ActiveScore = ed.Stats.CalculateActiveScore()
		ed.StaticScore = ed.Stats.CalculateStaticScore()
	}
}
