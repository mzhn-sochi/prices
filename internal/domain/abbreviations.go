package domain

type similar struct {
	words []string
}

func (s *similar) add(els ...string) {
	for _, item := range els {
		s.words = append(s.words, item)
	}
}
