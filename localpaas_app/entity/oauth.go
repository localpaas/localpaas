package entity

type OAuth struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Organization string `json:"org,omitempty"`
	BaseURL      string `json:"baseURL,omitempty"`
	RedirectURL  string `json:"redirectURL,omitempty"`
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
