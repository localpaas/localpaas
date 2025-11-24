package nginxservice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"

	crossplane "github.com/nginxinc/nginx-go-crossplane"
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

func (s *nginxService) InitAppConfig(_ context.Context, _ *entity.App, httpSettings *entity.AppHttpSettings) error {
	conf, err := nginxParseConfDefault(string(defaultConfTemplate))
	if err != nil {
		return apperrors.Wrap(err)
	}

	httpSettings.NginxSettings = &entity.NginxSettings{
		Enabled: true,
	}
	for _, block := range conf.BlocksByName("server", -1) {
		directives := gofn.MapSlice(block.AllDirectives(), func(d *crossplane.Directive) *entity.NginxDirective {
			return &entity.NginxDirective{
				Directive: d,
			}
		})
		if len(directives) > 0 {
			httpSettings.NginxSettings.ServerBlock = &entity.NginxServerBlock{
				Directives: directives,
			}
		}
	}
	return nil
}

func (s *nginxService) ApplyAppConfig(ctx context.Context, app *entity.App,
	httpSettings *entity.AppHttpSettings) (err error) {
	confPath := filepath.Join(config.Current.DataPathNginxEtcConf(), app.Key+".conf")

	// Not enabled, delete the config file, then return
	if httpSettings == nil || !httpSettings.Enabled {
		return s.RemoveAppConfig(ctx, app)
	}

	// Http settings is enabled
	buf := bytes.NewBuffer(make([]byte, 0, defaultConfBuf))
	for _, domain := range httpSettings.Domains {
		err = s.buildConfigForDomain(app, httpSettings, domain, buf)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	err = os.WriteFile(confPath, buf.Bytes(), defaultConfFileMode)
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

func (s *nginxService) buildConfigForDomain(app *entity.App, httpSettings *entity.AppHttpSettings,
	domain *entity.AppDomain, buf *bytes.Buffer) (err error) {
	if httpSettings.ForceHttps {
		conf, err := nginxParseConfDefault(string(forceHttpsConfTemplate))
		if err != nil {
			return apperrors.Wrap(err)
		}

		conf.IterBlocksByName("server", func(block *nginx.Block, _ int) bool {
			block.IterDirectivesByName("server_name", func(dir *crossplane.Directive, _ int) bool {
				dir.Args = []string{domain.Domain}
				return true
			})
			return true
		})

		err = nginxBuildConfDefault(conf, buf)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	var sslCertPath, sslKeyPath string
	if domain.SslCert.ID != "" {
		dirPath := fmt.Sprintf("/var/lib/localpaas/certs/%s", domain.SslCert.ID)
		sslCertPath = dirPath + "/certificate.crt"
		sslKeyPath = dirPath + "/private.key"
	}

	isRedirect := httpSettings.DomainRedirect != "" && domain.Domain != httpSettings.DomainRedirect
	if isRedirect {
		conf, err := nginxParseConfDefault(string(redirectConfTemplate))
		if err != nil {
			return apperrors.Wrap(err)
		}

		conf.IterBlocksByName("server", func(block *nginx.Block, _ int) bool {
			block.IterDirectivesByName("server_name", func(dir *crossplane.Directive, _ int) bool {
				dir.Args = []string{domain.Domain}
				return true
			})
			if sslCertPath != "" && sslKeyPath != "" {
				block.IterDirectivesByName("ssl_certificate", func(dir *crossplane.Directive, _ int) bool {
					dir.Args = []string{sslCertPath}
					return true
				})
				block.IterDirectivesByName("ssl_certificate_key", func(dir *crossplane.Directive, _ int) bool {
					dir.Args = []string{sslKeyPath}
					return true
				})
			}
			conf.IterBlocksByName("location", func(block *nginx.Block, _ int) bool {
				block.IterDirectivesByName("return", func(directive *crossplane.Directive, _ int) bool {
					target := fmt.Sprintf("https://%s$request_uri", httpSettings.DomainRedirect)
					directive.Args = []string{directive.Args[0], target}
					return true
				})
				return true
			})
			return true
		})

		err = nginxBuildConfDefault(conf, buf)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	if !isRedirect && httpSettings.NginxSettings.ServerBlock != nil {
		directives := gofn.MapSlice(httpSettings.NginxSettings.ServerBlock.Directives,
			func(directive *entity.NginxDirective) *crossplane.Directive {
				return directive.Directive
			})

		conf := nginx.NewConfig()
		conf.AddBlock(nginx.NewServerBlock(directives...))

		conf.IterBlocksByName("server", func(block *nginx.Block, _ int) bool {
			block.IterDirectivesByName("server_name", func(dir *crossplane.Directive, _ int) bool {
				dir.Args = []string{domain.Domain}
				return true
			})
			if sslCertPath != "" && sslKeyPath != "" {
				block.IterDirectivesByName("ssl_certificate", func(dir *crossplane.Directive, _ int) bool {
					dir.Args = []string{sslCertPath}
					return true
				})
				block.IterDirectivesByName("ssl_certificate_key", func(dir *crossplane.Directive, _ int) bool {
					dir.Args = []string{sslKeyPath}
					return true
				})
			}
			upstream := fmt.Sprintf("http://%s:%d", app.Key, httpSettings.ContainerPort)
			block.IterDirectivesByName("set", func(dir *crossplane.Directive, _ int) bool {
				if dir.Args[0] == "$upstream" {
					dir.Args[1] = upstream
				}
				return true
			})
			return true
		})

		err = nginxBuildConfDefault(conf, buf)
		if err != nil {
			return apperrors.Wrap(err)
		}
	}

	return nil
}

func nginxParseConfDefault(data string, options ...nginx.ParseOption) (*nginx.Config, error) {
	options = append(options, func(opts *crossplane.ParseOptions) {
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
	options = append(options, func(opts *crossplane.BuildOptions) {
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
