package product

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(input ProductCreate) (Product, error) {
	return s.repository.Create(input)
}

func (s *Service) GetByID(id int) (Product, error) {
	return s.repository.GetByID(id)
}

func (s *Service) List() ([]Product, error) {
	return s.repository.List()
}
