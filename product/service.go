package product

type Service interface {
	GetProductByID(param *getProductByIDRequest) (*Product, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetProductByID(param *getProductByIDRequest) (*Product, error) {
	return s.repo.GetProductByID(param.ProductID)
}
