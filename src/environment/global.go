package environment

import "sync"

type GlobalVarsStore interface {
	SaveGlobalVars(string) error
	GlobalVars() (string, error)
}

type GlovalVars struct {
	s     GlobalVarsStore
	cache string
	mx    sync.RWMutex
}

func NewGlovalVars(s GlobalVarsStore) (*GlovalVars, error) {
	cache, err := s.GlobalVars()
	if err != nil {
		return nil, err
	}

	return &GlovalVars{
		s:     s,
		cache: cache,
	}, nil
}

func (e *GlovalVars) Save(vars string) error {
	e.mx.Lock()
	defer e.mx.Unlock()

	e.cache = vars
	return e.s.SaveGlobalVars(e.cache)
}

func (e *GlovalVars) Get() string {
	e.mx.RLock()
	defer e.mx.RUnlock()
	return e.cache
}
