package domain

import (
	"cocoIM/ipconf/source"
	"sort"
	"sync"
)

type Dispatcher struct {
	// 候选人名单
	candidateTable map[string]*Endport
	sync.RWMutex
}

// 调度器对象
var dp *Dispatcher

func Init() {
	// 初始化调度器
	dp = &Dispatcher{}
	dp.candidateTable = make(map[string]*Endport)
	// 监听事件，并通过事件修改调度器里的候选人名单
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dp.addNode(event)
			case source.DelNodeEvent:
				dp.delNode(event)
			}
		}
	}()
}

// 调度请求
func Dispatch(ctx *IpConfContext) []*Endport {
	// 获得所有候选者列表
	candidateList := dp.GetCandidateList(ctx)
	// 每个endport都要算一下分数
	for _, endport := range candidateList {
		endport.CalculateScore(ctx)
	}
	// 按分数排序
	// 这里其实处理的不太好，应该使用融合分而不是简单的动态和静态分来暴力排序
	sort.Slice(candidateList, func(i, j int) bool {
		//先按照动态分排序
		if candidateList[i].ActiveScore > candidateList[j].ActiveScore {
			return true
		}
		//如果动态分相同，使用静态分排序
		if candidateList[i].ActiveScore == candidateList[j].ActiveScore {
			if candidateList[i].StaticScore > candidateList[j].StaticScore {
				return true
			}
			return false
		}
		return false
	})
	return candidateList
}

// 获取候选者列表
func (dp *Dispatcher) GetCandidateList(ctx *IpConfContext) []*Endport {
	dp.Lock()
	defer dp.Unlock()
	// 复制一份，避免影响源数据
	newCandidateList := make([]*Endport, 0, len(dp.candidateTable))
	for _, endport := range dp.candidateTable {
		newCandidateList = append(newCandidateList, endport)
	}
	return newCandidateList
}
func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()

	var (
		ok bool
		ed *Endport
	)
	// 给候选人列表添加 新的  endport
	if ed, ok = dp.candidateTable[event.Key()]; !ok {
		ed = NewEndport(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	// 添加成功后，更新状态
	ed.UpdateStat(&Stat{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})

}
func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}
