package domain

import "math"

// 状态机
// 服务器如果有多个的话，尽量保证机器配置是相同的，但是也可能由于各种原因，不同时间，不同地区的配置不同
// 因此要考虑不同机器在一起组群的时候，负载均衡依赖什么来分配

// 比如，如果单纯的按负载率，假设3台机器，ABC，AB是老机器，最大能承担1000个链接，而C是新机器，能承担5000个
// 倘若此时，C机器上已经存在4000个，剩余1000个，AB都是剩余200个，单纯的按负载率来说，ABC的负载率都是80%，但是明显C的性能还有1000的剩余，而AB显然要炸了
// 此时应该考虑分配给C，那么状态不能再是单纯的负载率，而是按剩余资源量

// 这里认为，每个机器有一个状态，这个状态将会表明当前机器的剩余资源量，负载均衡应该优先把拥有更多剩余资源的机器分配出去

// 这里数值是每一个网关机（endpoint）的资源指标

// Stat 资源状态
type Stat struct {
	// 当前机器总体持有的长连接数量的
	ConnectNum   float64 // 总体持有长连接数量的剩余值，也就是还有多少个可以分配的长连接
	MessageBytes float64 // 美妙收发消息的总字节数 的剩余值，比如带宽假定为10MB，现在收发9MB/s，那么剩余值就是1mb
}

// Avg 求平均值
// float64 num 要平均成num份
func (s *Stat) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}

func (s *Stat) Clone() *Stat {
	// 避免浅拷贝，直接重新赋值
	newStat := &Stat{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

// 删除最后一项数据
func (s *Stat) sub(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

func (s *Stat) Add(st *Stat) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

// 计算动态分
func (s *Stat) CalculateActiveScore() float64 {
	return getGB(s.MessageBytes)
}

// 把字节数转为GB，除以1<<30,也就是 1024*1024*1024（2^10 * 2^10 * 2^10）
func getGB(score float64) float64 {
	return decimal(score / (1 << 30))
}

// 四舍五入拿到一个整数，加 0.5 的目的是确保小数点后第三位大于等于 5 时，上舍入到整数部分。
func decimal(score float64) float64 {
	return math.Trunc(score*1e2+0.5) * 1e-2
}

// 计算静态分
func (s *Stat) CalculateStaticScore() float64 {
	// 目前无需计算,直接返回就好
	return s.ConnectNum
}
