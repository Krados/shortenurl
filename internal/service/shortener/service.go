package shortener

//go:generate mockery --name=Shortener --filename=mock_shortener.go --inpackage
type Shortener interface {
	Get(code string) (*ShortURL, error)
	Put(reqURL string) (*ShortURL, error)
}
