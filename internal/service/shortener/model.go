package shortener

type ShortURL struct {
	ID      int64  `json:"id"`
	Code    string `json:"code"`
	HashURL string `json:"hash_URL"`
	LongURL string `json:"long_URL"`
}
