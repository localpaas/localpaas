package nginxservice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tiendc/gofn"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/config"
	"github.com/localpaas/localpaas/localpaas_app/entity"
	"github.com/localpaas/localpaas/services/nginx"
)

const (
	defaultConfBuf      = 10 * 1024 // 10KB
	defaultConfFileMode = 0o755
)

const (
	nginxShareDir = "/usr/share/nginx" // dir path within container of nginx

	certsBaseDir = nginxShareDir + "/certs"
	certsMainDir = certsBaseDir + "/main"
	certsFakeDir = certsBaseDir + "/fake" // self-signed

	basicAuthDir = nginxShareDir + "/basic-auth"
)

var (
	defaultConfTemplate = func() []byte {
		data, err := os.ReadFile("config/nginx/app.conf.template")
		if err != nil {
			panic(err)
		}
		return data
	}()

	forceHttpsConfTemplate = func() []byte {
		data, err := os.ReadFile("config/nginx/force_https.conf.template")
		if err != nil {
			panic(err)
		}
		return data
	}()

	redirectConfTemplate = func() []byte {
		data, err := os.ReadFile("config/nginx/redirect.conf.template")
		if err != nil {
			panic(err)
		}
		return data
	}()
)

func (s *nginxService) GetDefaultNginxConfig() (*entity.NginxSettings, error) {
	conf, err := nginxParseConfDefault(string(defaultConfTemplate))
	if err != nil {
		return nil, apperrors.Wrap(err)
	}

	nginxSettings := &entity.NginxSettings{}
	for _, block := range conf.BlocksByName("server", -1) {
		directives := gofn.MapSlice(block.AllDirectives(), func(d *nginx.Directive) *entity.NginxDirective {
			return &entity.NginxDirective{
				Directive: d,
			}
		})
		if len(directives) > 0 {
			nginxSettings.ServerBlock = &entity.NginxServerBlock{
				Directives: directives,
			}
		}
	}
	return nginxSettings, nil
}

type AppConfigData struct {
	HttpSettings  *entity.AppHttpSettings
	RefSettingMap map[string]*entity.Setting

	confBuf *bytes.Buffer
}

func (s *nginxService) ApplyAppConfig(
	ctx context.Context,
	app *entity.App,
	data *AppConfigData,
) (err error) {
	httpSettings := data.HttpSettings
	confPath := filepath.Join(config.Current.DataPathNginxEtcConf(), app.Key+".conf")

	// Not enabled, delete the config file, then return
	if httpSettings == nil || !httpSettings.Enabled {
		return s.RemoveAppConfig(ctx, app)
	}

	data.confBuf = bytes.NewBuffer(make([]byte, 0, defaultConfBuf))
	for _, domain := range httpSettings.Domains {
		err = s.buildConfigForDomain(app, domain, data)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	err = os.WriteFile(confPath, data.confBuf.Bytes(), defaultConfFileMode)
	if err != nil {
		return apperrors.Wrap(err)
	}

	// Requests nginx to reload the config files
	err = s.ReloadNginxConfig(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (s *nginxService) buildConfigForDomain(
	app *entity.App,
	domain *entity.AppDomain,
	data *AppConfigData,
) (err error) {
	if !domain.Enabled {
		return nil
	}

	if domain.ForceHttps {
		err = s.buildForceHttpsConfig(domain, data)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	isRedirect := domain.DomainRedirect != "" && domain.Domain != domain.DomainRedirect
	if isRedirect {
		err = s.buildRedirectionConfig(domain, data)
	} else {
		err = s.buildMainConfig(app, domain, data)
	}
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}

func (s *nginxService) buildForceHttpsConfig(
	domain *entity.AppDomain,
	data *AppConfigData,
) error {
	conf, err := nginxParseConfDefault(string(forceHttpsConfTemplate))
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = conf.IterBlocksByName("server", func(block *nginx.Block, _ int) (bool, error) {
		block.SetDirectiveArgs("server_name", []string{domain.Domain}, 1)
		return true, nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = nginxBuildConfDefault(conf, data.confBuf)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *nginxService) buildRedirectionConfig(
	domain *entity.AppDomain,
	data *AppConfigData,
) error {
	conf, err := nginxParseConfDefault(string(redirectConfTemplate))
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = conf.IterBlocksByName("server", func(block *nginx.Block, _ int) (bool, error) {
		block.SetDirectiveArgs("server_name", []string{domain.Domain}, 1)

		if domain.SslCert.ID != "" {
			certFile, keyFile := getSslFilePaths(domain.SslCert.ID)
			block.SetDirectiveArgs("ssl_certificate", []string{certFile}, 1)
			block.SetDirectiveArgs("ssl_certificate_key", []string{keyFile}, 1)
		}

		err = conf.IterBlocksByName("location", func(block *nginx.Block, _ int) (bool, error) {
			err = block.IterDirectivesByName("return", func(dir *nginx.Directive, _ int) (bool, error) {
				target := fmt.Sprintf("https://%s$request_uri", domain.DomainRedirect)
				dir.Args = []string{dir.Args[0], target}
				return true, nil
			})
			if err != nil {
				return false, apperrors.Wrap(err)
			}
			return true, nil
		})
		if err != nil {
			return false, apperrors.Wrap(err)
		}
		return true, nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = nginxBuildConfDefault(conf, data.confBuf)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *nginxService) buildMainConfig(
	app *entity.App,
	domain *entity.AppDomain,
	data *AppConfigData,
) error {
	if domain.NginxSettings.ServerBlock == nil {
		return nil
	}
	directives := gofn.MapSlice(domain.NginxSettings.ServerBlock.Directives,
		func(directive *entity.NginxDirective) *nginx.Directive {
			return directive.Directive
		})

	conf := nginx.NewConfig()
	conf.AddBlock(nginx.NewServerBlock(directives...))

	err := conf.IterBlocksByName("server", func(serverBlock *nginx.Block, _ int) (bool, error) {
		s.buildSslConfig(domain, serverBlock)
		err := s.buildUpstreamConfig(app, domain, serverBlock)
		if err != nil {
			return false, apperrors.Wrap(err)
		}

		rootLocationBlock, _ := serverBlock.AddLocationBlock("/", nil)
		s.buildWebsocketConfig(domain, rootLocationBlock)
		s.buildBasicAuthConfig(domain, rootLocationBlock)

		return true, nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}

	err = nginxBuildConfDefault(conf, data.confBuf)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *nginxService) buildSslConfig(
	domain *entity.AppDomain,
	serverBlock *nginx.Block,
) {
	serverBlock.AddDirectiveNamesIfNotExist("server_name")
	if domain.Domain != "" {
		serverBlock.SetDirectiveArgs("server_name", []string{domain.Domain}, 1)
	}

	certFile, keyFile := getSslFilePaths(domain.SslCert.ID)
	serverBlock.AddDirectiveNamesIfNotExist("ssl_certificate", "ssl_certificate_key")
	serverBlock.SetDirectiveArgs("ssl_certificate", []string{certFile}, 1)
	serverBlock.SetDirectiveArgs("ssl_certificate_key", []string{keyFile}, 1)
}

func (s *nginxService) buildUpstreamConfig(
	app *entity.App,
	domain *entity.AppDomain,
	serverBlock *nginx.Block,
) error {
	upstream := fmt.Sprintf("http://%s:%d", app.Key, domain.ContainerPort)
	err := serverBlock.IterDirectivesByName("set", func(dir *nginx.Directive, _ int) (bool, error) {
		if dir.Args[0] == "$upstream" {
			dir.Args[1] = upstream
		}
		return true, nil
	})
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *nginxService) buildWebsocketConfig(
	domain *entity.AppDomain,
	locationBlock *nginx.Block,
) {
	fn := gofn.If(domain.WebsocketEnabled, locationBlock.AddDirectivesIfNotExist, locationBlock.RemoveDirectives)
	fn(
		nginx.NewDirective("proxy_set_header", []string{"Upgrade", "$http_upgrade"}),
		nginx.NewDirective("proxy_set_header", []string{"Connection", "\"upgrade\""}),
		nginx.NewDirective("proxy_http_version", []string{"1.1"}),
	)
}

func (s *nginxService) buildBasicAuthConfig(
	domain *entity.AppDomain,
	locationBlock *nginx.Block,
) {
	if domain.BasicAuth.ID == "" {
		locationBlock.RemoveDirectiveNames("auth_basic", "auth_basic_user_file")
		return
	}

	locationBlock.AddDirectiveNamesIfNotExist("auth_basic", "auth_basic_user_file")
	locationBlock.SetDirectiveArgs("auth_basic", []string{"\"Restricted Access\""}, 1)
	authFilepathInContainer := basicAuthDir + "/" + domain.BasicAuth.ID
	locationBlock.SetDirectiveArgs("auth_basic_user_file", []string{authFilepathInContainer}, 1)
}

func getSslFilePaths(sslID string) (certFile, keyFile string) {
	if sslID == "" {
		return filepath.Join(certsFakeDir, "local.crt"), filepath.Join(certsFakeDir, "local.key")
	}
	return filepath.Join(certsMainDir, sslID+".crt"), filepath.Join(certsMainDir, sslID+".key")
}

func nginxParseConfDefault(data string, options ...nginx.ParseOption) (*nginx.Config, error) {
	options = append(options, func(opts *nginx.ParseOptions) {
		opts.SingleFile = true
		opts.SkipDirectiveContextCheck = true
		opts.ParseComments = true
	})
	conf, err := nginx.ParseString(data, options...)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return conf, nil
}

func nginxBuildConfDefault(conf *nginx.Config, buf *bytes.Buffer, options ...nginx.BuildOption) error {
	options = append(options, func(opts *nginx.BuildOptions) {
		opts.Tabs = true
	})
	if buf.Len() > 0 {
		_, _ = buf.Write([]byte("\n\n"))
	}
	err := nginx.Build(conf, buf, options...)
	if err != nil {
		return apperrors.Wrap(err)
	}
	return nil
}

func (s *nginxService) RemoveAppConfig(ctx context.Context, app *entity.App) error {
	confPath := filepath.Join(config.Current.DataPathNginxEtcConf(), app.Key+".conf")
	err := os.Remove(confPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return apperrors.Wrap(err)
	}

	// Requests nginx to reload the config files
	err = s.ReloadNginxConfig(ctx)
	if err != nil {
		return apperrors.Wrap(err)
	}

	return nil
}
