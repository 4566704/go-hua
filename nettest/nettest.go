package nettest

import (
	"fmt"
	"sync"
	"time"
)

type TestResult struct {
	SuccessCount int // 成功计次
	FailCount    int // 失败计次
	TotalCount   int // 总数
	TotalTime    int // 总耗时
	Avg          int // 平均耗时 平均值
	Min          int // 最小耗时 最小值
	Max          int // 最大耗时 最大值
}

type NetTest struct {
	Id      int
	Name    string
	Addr    string
	Port    int
	Timeout int64
	TestResult
	Ping     *Ping
	IsStop   bool
	IsEcho   bool
	IsUpdate bool
	mux      sync.Mutex
}

func NewNetTest(id int, name string, addr string, port int, timeout int64, echo bool) *NetTest {
	t := new(NetTest)
	t.Id = id
	t.Name = name
	t.Addr = addr
	t.Port = port
	t.Timeout = timeout
	t.IsEcho = echo
	t.IsUpdate = false
	t.Ping = NewPing(addr, port, timeout)
	return t
}

func (t *NetTest) Set(name string, addr string, port int, timeout int64) {
	t.mux.Lock()
	t.Name = name
	t.Addr = addr
	t.Port = port
	t.Timeout = timeout
	t.IsUpdate = true
	t.mux.Unlock()
	t.Ping = NewPing(addr, port, timeout)
}

func (t *NetTest) Start() {
	go t.process()
}

func (t *NetTest) Stop() {
	t.mux.Lock()
	t.IsStop = true
	t.mux.Unlock()
}

func (t *NetTest) process() {
	isLoop := true
	for isLoop {

		t.mux.Lock()
		timeout := int(t.Timeout)
		if t.IsUpdate {
			t.Ping.Set(t.Addr, t.Port, t.Timeout)
		}
		if t.IsStop {
			isLoop = false
		}
		t.mux.Unlock()
		if !isLoop {
			break
		}
		et := t.Ping.Test()
		if et > 0 {
			timeout = timeout - et
		}
		if timeout > int(timeout) {
			timeout = int(timeout)
		}
		t.mux.Lock()
		t.TotalCount++
		if et >= 0 {
			t.SuccessCount++
		} else {
			t.FailCount++
		}
		if et == 0 {
			t.TotalTime++
		} else if et > 0 {
			t.TotalTime += et
		}
		if t.SuccessCount > 0 {
			t.Avg = t.TotalTime / t.SuccessCount
		}
		if (et < t.Min || t.Min == 0) && et > 0 {
			t.Min = et
		}
		if et > t.Max {
			t.Max = et
		}
		t.mux.Unlock()

		if timeout > 0 {
			time.Sleep(time.Duration(timeout) * time.Millisecond)
		}
		if t.IsEcho {
			fmt.Printf("来自 %s 的回复: 时间=%dms\n", t.Addr, et)
		}
	}
}

func (t *NetTest) GetResult(isReset bool) TestResult {
	t.mux.Lock()
	r := TestResult{}
	r = t.TestResult
	if isReset {
		t.SuccessCount = 0
		t.FailCount = 0
		t.TotalCount = 0
		t.TotalTime = 0
		t.Avg = 0
		t.Min = 0
		t.Max = 0
	}
	t.mux.Unlock()
	return r
}
