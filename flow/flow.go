package flow

import (
	"sync/atomic"
)

// 统计流量 不限制流量
type Flow struct {
	SendFlow  int64 `json:"sendFlow"`
	RecvFlow  int64 `json:"recvFlow"`
	FlowLimit int64 `json:"flowLimit"`
}

func (f *Flow) Add(send, recv int64) {
	atomic.AddInt64(&f.SendFlow, send)
	atomic.AddInt64(&f.RecvFlow, recv)
}

func (f *Flow) Reset() (int64, int64) {
	send := atomic.SwapInt64(&f.SendFlow, 0)
	recv := atomic.SwapInt64(&f.RecvFlow, 0)
	return send, recv
}

func (f *Flow) Get() (int64, int64) {
	send := atomic.LoadInt64(&f.SendFlow)
	recv := atomic.LoadInt64(&f.RecvFlow)
	return send, recv
}
