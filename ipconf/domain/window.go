package domain

// 默认采样间隔 单位 秒
const windowSize = 5

// 状态观测窗口
type stateWindow struct {
	// 状态采样队列,每秒一次
	stateQueue []*Stat
	// 统计通道，用于传输状态统计后的结果
	statChan chan *Stat
	// 统计结果
	sumStat *Stat
	// 统计采样的当前递增序号
	idx int64
}

// 创建一个新的观测窗口
func newStateWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*Stat, windowSize),
		statChan:   make(chan *Stat),
		sumStat:    &Stat{},
	}
}

func (s *stateWindow) getStat() *Stat {
	res := s.sumStat.Clone()
	res.Avg(windowSize)
	return res
}
func (s *stateWindow) updateSumStat(newStat *Stat) {
	// 窗口里始终保持5个数值
	// 减去最后一个状态state
	// 这里是直接让idx 模 窗口大小，也就是窗口里的5个值分别替换一次
	s.sumStat.sub(s.stateQueue[s.idx%windowSize])
	// 更新最新的stat
	s.stateQueue[s.idx%windowSize] = newStat
	// 计算新的值
	s.sumStat.Add(newStat)
	s.idx++
}
