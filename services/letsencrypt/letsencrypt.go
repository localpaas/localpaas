package letsencrypt

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"os"

	"github.com/go-acme/lego/providers/http/webroot"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
)

type Client struct {
	client *lego.Client
	user   *User
}

type User struct {
	Email        string
	Registration *registration.Resource
	PrivateKey   crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.PrivateKey
}

func NewClient(email, http01NginxRoot string) (client *Client, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to generate private key for user")
	}

	user := User{
		Email:      email,
		PrivateKey: privateKey,
	}
	cfg := lego.NewConfig(&user)

	// Default is RSA2048
	// cfg.Certificate.KeyType = certcrypto.RSA2048

	c, err := lego.NewClient(cfg)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to create lego client")
	}

	webrootProvider, err := webroot.NewHTTPProvider(http01NginxRoot)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to create http provider for webroot")
	}

	err = c.Challenge.SetHTTP01Provider(webrootProvider)
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to set http-01 challenge")
	}

	return &Client{
		client: c,
		user:   &user,
	}, nil
}

func NewClientFromConfig(cfg *config.Config) (client *Client, err error) {
	nginxWebroot := cfg.DataPathNginxShareDomains()
	err = os.MkdirAll(nginxWebroot, 0755) //nolint
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to create nginx webroot directory")
	}
	return NewClient(cfg.SSL.LeUserEmail, nginxWebroot)
}

func (client *Client) registerUser(_ context.Context) error {
	if client.user.Registration != nil {
		return nil
	}
	reg, err := client.client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return apperrors.New(err).WithMsgLog("failed to register user")
	}
	client.user.Registration = reg
	return nil
}

func (client *Client) ObtainCertificate(ctx context.Context, domains []string,
	savePrivateKeyPath, saveCertPath string) (*certificate.Resource, error) {
	// New users will need to register
	err := client.registerUser(ctx)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	certificates, err := client.client.Certificate.Obtain(certificate.ObtainRequest{
		Domains: domains,
		Bundle:  true,
	})
	if err != nil {
		return nil, apperrors.New(err).WithMsgLog("failed to obtain certificate")
	}

	if savePrivateKeyPath != "" {
		privKeyFile, err := os.Create(savePrivateKeyPath)
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to create private key file")
		}
		_, err = privKeyFile.Write(certificates.PrivateKey)
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to write data to private key file")
		}
	}
	if saveCertPath != "" {
		certFile, err := os.Create(saveCertPath)
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to create certificate file")
		}
		_, err = certFile.Write(certificates.Certificate)
		if err != nil {
			return nil, apperrors.New(err).WithMsgLog("failed to write data to certificate file")
		}
	}

	return certificates, nil
}
