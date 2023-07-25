package environment

type GlobalVarsStore interface {
	SaveGlobalVars(string) error
	GlobalVars() (string, error)
}

type GlovalVars struct {
	s     GlobalVarsStore
	cache string
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
	e.cache = vars
	return e.s.SaveGlobalVars(e.cache)
}

func (e *GlovalVars) Get() string {
	return e.cache
}
