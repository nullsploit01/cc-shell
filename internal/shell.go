package internal

type Shell struct{}

func NewShell() *Shell {
	return &Shell{}
}

func (s *Shell) Run() {
}
