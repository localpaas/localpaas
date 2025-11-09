package entity

type OAuth struct {
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret"`
	Organization string   `json:"org,omitempty"`
	CallbackURL  string   `json:"callbackURL,omitempty"`
	AuthURL      string   `json:"authURL,omitempty"`
	TokenURL     string   `json:"tokenURL,omitempty"`
	ProfileURL   string   `json:"profileURL,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`

	// Salt used to encrypt the secret
	Salt string `json:"salt,omitempty"`
}

func (s *Setting) ParseOAuth() (*OAuth, error) {
	if s != nil && s.Data != "" {
		res := &OAuth{}
		err := s.parseData(res)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
	return nil, nil
}
