package shortener

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/bwmarrin/snowflake"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInValidURL   = errors.New("invalid URL")
	ErrCodeNotFound = errors.New("code not found")
)

type shortener struct {
	repo     Repository
	validate *validator.Validate
	sfn      *snowflake.Node
}

func New(repo Repository) Shortener {
	node, _ := snowflake.NewNode(0)
	return &shortener{
		repo:     repo,
		validate: validator.New(),
		sfn:      node,
	}
}

func (s *shortener) Get(code string) (*ShortURL, error) {
	return s.repo.GetWithCode(code)
}

func (s *shortener) Put(reqURL string) (res *ShortURL, err error) {
	// check reqURL is valid or not
	err = s.validate.Var(reqURL, "required,url")
	if err != nil {
		err = ErrInValidURL
		return
	}

	// get data with hashURL first
	hashURL := fmt.Sprintf("%x", sha256.Sum256([]byte(reqURL)))
	hr, err := s.repo.GetWithHashURL(hashURL)
	if err != nil && err != ErrCodeNotFound {
		return
	}
	if hr != nil {
		res = hr
		return
	}

	// no data exist then create a new one
	id := s.sfn.Generate()
	var tmp ShortURL
	tmp.ID = id.Int64()
	tmp.Code = id.Base58()
	tmp.LongURL = reqURL
	tmp.HashURL = hashURL
	err = s.repo.Put(&tmp)
	if err != nil {
		return
	}
	res = &tmp

	return
}
