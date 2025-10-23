package entity

type SSHKey struct {
	PrivateKey string `json:"privateKey"`
	Salt       string `json:"salt,omitempty"`
}

func (s *Setting) ParseSSHKey(decrypt bool) (*SSHKey, error) {
	if s != nil && s.Data != "" {
		res := &SSHKey{}
		err := s.parseData(res)
		if err != nil {
			return nil, err
		}
		if decrypt { //nolint
			// TODO: implement encryption
		}
		return res, nil
	}
	return nil, nil
}
