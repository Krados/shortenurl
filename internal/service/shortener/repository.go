package shortener

//go:generate mockery --name=Repository --filename=mock_repository.go --inpackage
type Repository interface {
	GetWithCode(code string) (*ShortURL, error)
	Put(shortURL *ShortURL) error
	GetWithHashURL(hashURL string) (*ShortURL, error)
}
