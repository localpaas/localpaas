package main

import (
	"context"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-github/v79/github"
	"github.com/palantir/go-githubapp/githubapp"
	"go.uber.org/fx"
	"golang.org/x/oauth2"

	"github.com/localpaas/localpaas/localpaas_app/cmd/internal"
	"github.com/localpaas/localpaas/localpaas_app/registry"
)

func main() {
	provides := []any{
		context.Background,
	}
	provides = append(provides, registry.Provides...)

	app := fx.New(
		fx.Provide(provides...),
		fx.Invoke(internal.InitLogger),
		fx.Invoke(internal.InitConfig),
		fx.Invoke(internal.InitDBConnection),
		fx.Invoke(internal.InitCache),
		fx.Invoke(internal.InitDockerManager),
		fx.Invoke(internal.InitJWTSession),
		fx.Invoke(internal.InitHTTPServer),
		fx.Invoke(abc),
	)

	app.Run()
}

func abc() {
	ctx := context.Background()

	//client, err := NewClient(ctx, ClientOptions{PrivateKey: pk, RepoURL: "", GithubAppID: "2258661"})
	//if err != nil {
	//	fmt.Println("Error getting client:", err)
	//	return
	//}
	//
	//token, err := client.GetToken2(ctx, "localpaas-test")
	//if err != nil {
	//	fmt.Println("Error getting token:", err)
	//	return
	//}
	//
	//token = "Token " + token
	//cc2 := client.client.WithAuthToken(token)
	////cc2 := client.client
	//
	//repos, _, err := cc2.Apps.ListRepos(ctx, nil)
	//print(repos, err)

	cc := githubapp.NewClientCreator(
		"https://api.github.com/",
		"",
		2258661,
		[]byte(pk),

		//githubapp.Config{
		//	V3APIURL: "https://api.github.com/",
		//	App: struct {
		//		IntegrationID int64  `yaml:"integration_id" json:"integrationId"`
		//		WebhookSecret string `yaml:"webhook_secret" json:"webhookSecret"`
		//		PrivateKey    string `yaml:"private_key" json:"privateKey"`
		//	}{
		//		IntegrationID: 2258661,
		//		WebhookSecret: "abc123",
		//		PrivateKey:    pk,
		//	},
		//	OAuth: struct {
		//		ClientID     string `yaml:"client_id" json:"clientId"`
		//		ClientSecret string `yaml:"client_secret" json:"clientSecret"`
		//	}{
		//		ClientID:     "Iv23liObQsEr3GigALXt",
		//		ClientSecret: "a469f4bbc4612cb4075e540af6a8f3abcf1b00d2",
		//	},
		//},
	)
	//if err != nil {
	//	panic(err)
	//}

	//ghClient, err := cc.NewAppClient()
	//if err != nil {
	//	panic(err)
	//}
	//
	//rr, _, err := ghClient.Apps.ListInstallations(ctx, nil)
	//if err != nil {
	//	panic(err)
	//}
	//print(rr)
	//
	//ii, _, err := ghClient.Apps.CreateInstallationTokenListRepos(ctx, *rr[0].ID, nil)
	//if err != nil {
	//	panic(err)
	//}
	//print(ii)
	//
	//ghClient = ghClient.WithAuthToken(ii.GetToken())
	//
	//repos, _, err := ghClient.Repositories.ListAll(ctx, nil)
	//if err != nil {
	//	panic(err)
	//}
	//print(repos, err)

	//prCommentHandler := &PRCommentHandler{
	//	ClientCreator: cc,
	//	preamble:      config.AppConfig.PullRequestPreamble,
	//}
	//
	//webhookHandler := githubapp.NewDefaultEventDispatcher(config.Github, prCommentHandler)
	//
	//http.Handle(githubapp.DefaultWebhookRoute, webhookHandler)
	//
	//addr := fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port)
	//logger.Info().Msgf("Starting server on %s...", addr)
	//err = http.ListenAndServe(addr, nil)
	//if err != nil {
	//	panic(err)
	//}

	// Replace with your GitHub App ID and private key path
	appID := int64(2258661)
	installationID := int64(0) // int64(96871321) // The ID of the specific installation
	_ = ctx
	_ = cc

	// Create a new installation transport
	itr, err := ghinstallation.New(http.DefaultTransport, appID, installationID, []byte(pk))
	if err != nil {
		print("Failed to create installation transport: %v", err)
	}

	// Create a GitHub client using the installation transport
	client := github.NewClient(&http.Client{Transport: itr})

	rr, _, err := client.Apps.ListInstallations(context.Background(), nil)
	print(rr, err)

	// Example: List repositories for the installation
	repos, _, err := client.Apps.ListRepos(context.Background(), nil)
	if err != nil {
		print("Failed to list repositories: %v", err)
	}
	print(repos)

	//fmt.Println("Repositories accessible by the installation:")
	//for _, repo := range repos {
	//	fmt.Printf("- %s\n", *repo.FullName)
	//}
}

var pk = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqV7zjfk4kpjryA/rUbkvUxUXBXUWkJUI+n/fMzmo7uzbOLPb
apmZSRoc3N8lN9OM16B5I6V9OSxcqi5ITbBHjP2S9YpxHnSp6+H+xAEezwqeguKq
qBMtwB+tfLcJseNPuERGnts0+DoDnWxlfIyxWjKwyNHxeE8mMTeToJGZvFnwHbch
kYF9o1yqaKxliNIUZH1U2pUVjzYEJvAr6rOMtAavVUAfHMDpQr5UarTraG1V/BYR
JEJMhgqnReKtgZes7b17am7MKFLehk1Zl01/vUiwTW38v9T8cfvoxNQkwCxcOFEe
nRu+w2R6qZCGaUREMJI48Kqz7GHIfpLz5XE89QIDAQABAoIBAAqIqUUvNuGHOULc
Gpqeld7o9OGWAq1DV6ecxFj/QQ57FLdWGFNs8ielxB3IkxwBwES46k/ZPwyLQcgj
0MNkL08JhlZXIenu/5J4H8w49Z2S34DscAi8CKBmV/flumih5pdjR7EhmQ3HLepg
j39LkMw3Ut9qj0YYRhvRhRD7fUBFjcvkqzQyYjN3aHTBdux+avhJntkSjzuWvUm+
bBAeW8sSCw847beQLR998pjmq32S+CtFo2WDWCB6ml2iSNu9Z6I9JvRcxNJKcwq2
/HC90jmkjesbMEN4JGPVDROBakMN6P7YE1OBQHKOsR/hV2Is+9g47tQmXL/v/YNN
biy+QRECgYEA1na5UAw+LH9YT4fH4TekqnT8IQYgfZOVyTglHvWlAu5Z4QTO3jB0
RDVNdeEEDT0te4rDltwPRBgPsWq0e7JhUPuoxesedMvErFncKAzJVuzhgQe6kgkl
DLsRLOnnDoz3niW+8+EVAOnAKEhi4PPuKsOpfNgbmxpbxAD0nFeVKnMCgYEAyix9
P0CUD7SfXeIHNM50i24V+6Mgad+dqbGDFhfTQRPxQcyMlAxWZkRERVykSCvZPkpE
vmrynatcAaYhegkagyjMiC5vq5TgQ29kbxvVnyzuX0cjgOfm9sAfLuHN0td/WooW
bVX985/p7nsmqIDeo8aMV8nZp5+7T+czOgYrmPcCgYEApZPIAvtQzT4MSmrEcSdq
pIfpWP3H++wECvSkBEEXAypOjiIOLREv0rFseoQdgrMm+GjsFP8Vcjc5dnqxmYbh
B4skhJnAS+A+mepOxUUJ9WixudLjwnf4+Nk8q3ZzA5LmYl02DkhK1QejafZpENCD
otSpeE/JEHjLPIqYUFR2P+0CgYBPqGPY7jroTUSVmx83wTjOTxu+QDHfRUo7uENs
CwsjkxX11QB2vL22IaO8qJnaDdzU3DJlzIG3efMQe0KpcLjPgN3FsnYxZsFOEK/D
z3XF99AcHk1w1u57vosKz9FnB52YMNxRTBZ5TULAwikjL1OJuAtH87Ich6UoAHKA
pDm9VQKBgAnO307QllmKju0GuCBrj2tXEmNYtdzzFZKLt56C5Jt7Qdq6H+aAuPSm
lRJcHK25NbiMjYuORF8PYR+Hh6nXL+iXPd7AFBsX/9NPnNnemFa94rZp/+9eHejc
Y+8mW3YQWLnBXOkn55nYX/Z5tIbyYATM82Ma7IEUF9fanf1hybLE
-----END RSA PRIVATE KEY-----`

type ClientOptions struct {
	PrivateKey  string
	RepoURL     string
	GithubAppID string
}

func NewClient(ctx context.Context, options ClientOptions) (*github.Client, error) {
	signedToken, err := getGitHubAppToken(options.GithubAppID, options.PrivateKey)
	if err != nil {
		return nil, err
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedToken})

	httpClient := oauth2.NewClient(ctx, src)

	return github.NewClient(httpClient), nil
}

//func (g *Client) GetToken2(ctx context.Context, org string) (string, error) {
//	//is, _, err := g.client.Apps.ListInstallations(ctx, nil)
//	//if err != nil {
//	//	return "", err
//	//}
//	//print(is)
//
//	installation, _, err := g.client.Apps.FindOrganizationInstallation(ctx, org)
//	if err != nil {
//		return "", err
//	}
//	installationToken, _, err := g.client.Apps.CreateInstallationToken(ctx, *installation.ID, nil)
//	if err != nil {
//		return "", err
//	}
//	return installationToken.GetToken(), err
//}

func getGitHubAppToken(appID, privateKey string) (string, error) {
	// Create the JWT token using the app's private key
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.MapClaims{
		"iss": appID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(), // Token expiration time (10 minutes)
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		return "", err
	}

	// Sign the token with the private key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
