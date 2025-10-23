package entity

type S3Storage struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Salt            string `json:"salt,omitempty"`
	Region          string `json:"region,omitempty"`
	Bucket          string `json:"bucket,omitempty"`
}

func (s *Setting) ParseS3Storage(decrypt bool) (*S3Storage, error) {
	if s != nil && s.Data != "" {
		res := &S3Storage{}
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
