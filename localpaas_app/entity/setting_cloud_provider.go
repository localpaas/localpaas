package entity

type CloudProviderAWS struct {
	AccessKeyID string         `json:"accessKeyID"`
	SecretKey   EncryptedField `json:"secretKey"`
	Region      string         `json:"region,omitempty"`
}
