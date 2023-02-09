package shortener

import "gorm.io/gorm"

type mySQLRepository struct {
	db *gorm.DB
}

func NewMySQLRepository(db *gorm.DB) Repository {
	return &mySQLRepository{
		db: db,
	}
}

func (m *mySQLRepository) GetWithCode(code string) (res *ShortURL, err error) {
	var tmp ShortURL
	err = m.db.First(&tmp, &ShortURL{Code: code}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrCodeNotFound
		}
		return
	}
	res = &tmp

	return
}

func (m *mySQLRepository) Put(shortURL *ShortURL) error {
	result := m.db.Create(&shortURL)
	return result.Error
}

func (m *mySQLRepository) GetWithHashURL(hashURL string) (res *ShortURL, err error) {
	var tmp ShortURL
	err = m.db.First(&tmp, &ShortURL{HashURL: hashURL}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrCodeNotFound
		}
		return
	}
	res = &tmp

	return
}
