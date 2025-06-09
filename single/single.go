package single

import "sync"

type Single struct {
	mu     sync.Mutex
	values map[string]int64
}

// 创建变量
/* var counters = Single{
	values: make(map[string]int64),
} */
func (s *Single) Init() {
	s.values = make(map[string]int64)
}

func (s *Single) Get(key string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[key]
}

func (s *Single) Incr(key string, value int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] += value
	return s.values[key]
}

func (s *Single) Set(key string, val int64) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.values[key]
	s.values[key] = val
	return old
}
