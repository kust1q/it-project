package showcase

func New(repo Repository) *Service {
	return &Service{repo: repo}
}
