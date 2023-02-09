package shortener

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPut_InValidURL(t *testing.T) {
	mr := &MockRepository{}
	serv := New(mr)
	_, err := serv.Put("123456")
	assert.Equal(t, ErrInValidURL, err)
	_, err = serv.Put("")
	assert.Equal(t, ErrInValidURL, err)
}

func TestPut_GetWithHashURLError(t *testing.T) {
	mr := &MockRepository{}
	someBadErr := errors.New("some bad error")
	mr.On("GetWithHashURL", mock.Anything).Return(nil, someBadErr)
	serv := New(mr)
	_, err := serv.Put("https://www.google.com")
	assert.Equal(t, someBadErr, err)
}

func TestPut_GetWithHashURLHit(t *testing.T) {
	s := ShortURL{
		ID:      1,
		Code:    "abc",
		HashURL: "HashURL",
		LongURL: "https://www.google.com",
	}
	mr := &MockRepository{}
	mr.On("GetWithHashURL", mock.Anything).Return(&s, nil)
	serv := New(mr)
	res, err := serv.Put("https://www.google.com")
	assert.Nil(t, err)
	assert.Equal(t, res.Code, s.Code)
}

func TestPut_RepoPutError(t *testing.T) {
	someBadErr := errors.New("some bad error")
	mr := &MockRepository{}
	mr.On("GetWithHashURL", mock.Anything).Return(nil, ErrCodeNotFound)
	mr.On("Put", mock.Anything).Return(someBadErr)
	serv := New(mr)
	res, err := serv.Put("https://www.google.com")
	assert.Nil(t, res)
	assert.Equal(t, someBadErr, err)
}

func TestPut_Success(t *testing.T) {
	mr := &MockRepository{}
	mr.On("GetWithHashURL", mock.Anything).Return(nil, ErrCodeNotFound)
	mr.On("Put", mock.Anything).Return(nil)
	serv := New(mr)
	res, err := serv.Put("https://www.google.com")
	assert.NotNil(t, res)
	assert.Equal(t, nil, err)
}
