package toggle

type Toggle struct {
	EnableSomething bool `yaml:"enable_something"`
}

func GetToggle(toggle *Toggle) *Toggle {
	return toggle
}
