package rate

import (
	"sync/atomic"
	"time"
)

type Rate struct {
	Limit   int32     `json:"limit"`   // 限制速度 字节 0为不限制
	Surplus int32     `json:"surplus"` // 剩余流量 字节
	Now     int32     `json:"now"`     // 接收(下载/下行) 当前流量 字节
	Max     int32     `json:"max"`     // 接收(下载/下行) 最大流量 字节
	IsStop  chan bool `json:"-"`
}

func NewRate(bandwidth int) *Rate {
	// 带宽 应该是Mbps 要换算成Mbyte
	limit := bandwidth * 1024 * 1024 / 8
	r := new(Rate)
	r.Limit = int32(limit)
	r.Now = 0
	r.Surplus = int32(limit)
	return r
}

func (r *Rate) SetLimit(bandwidth int) {
	// 带宽 应该是Mbps 要换算成Mbyte
	limit := 0
	if bandwidth > 0 {
		limit = bandwidth * 1024 * 1024 / 8
	} else {
		limit = 0
	}

	atomic.StoreInt32(&r.Limit, int32(limit))
}

func (r *Rate) Start() {
	go r.proc()
}

func (r *Rate) Stop() {
	r.IsStop <- true
}

func (r *Rate) proc() {
	ticker := time.NewTicker(time.Second * 1)

	for {
		select {
		case <-ticker.C:
			r.reset()
		case <-r.IsStop:
			ticker.Stop()
			return
		}
	}
}

func (r *Rate) reset() {
	n := r.Limit - atomic.LoadInt32(&r.Surplus)
	atomic.StoreInt32(&r.Now, n)
	atomic.StoreInt32(&r.Surplus, r.Limit)

	if n > 0 {
		now := atomic.LoadInt32(&r.Now)
		if now > 0 && now > atomic.LoadInt32(&r.Max) {
			atomic.StoreInt32(&r.Max, now)
		}
	}
	//fmt.Printf("now:%d limit:%d Surplus:%d \n", n, r.Limit, r.Limit)
}

func (r *Rate) GetNow() int {
	n := atomic.LoadInt32(&r.Now)
	return int(n)
}

func (r *Rate) ResetMax() int {
	n := atomic.SwapInt32(&r.Now, 0)
	return int(n)
}

func (r *Rate) Add(size int) {
	if atomic.LoadInt32(&r.Surplus) > 0 || atomic.LoadInt32(&r.Limit) == 0 {
		atomic.AddInt32(&r.Surplus, -int32(size))
		return
	}
	for {
		//fmt.Println("等待")
		time.Sleep(time.Millisecond * 10)
		if atomic.LoadInt32(&r.Surplus) > 0 {
			atomic.AddInt32(&r.Surplus, -int32(size))
			return
		}
	}
}
