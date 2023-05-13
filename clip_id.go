package svgJoin

import "sync"

type clipId struct {
	id uint64
	rw sync.Mutex
}

func (s *clipId) get() (id uint64) {
	s.rw.Lock()
	defer s.rw.Unlock()
	id = s.id
	s.id++
	return
}
