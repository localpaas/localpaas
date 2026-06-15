package acme

import (
	"github.com/go-acme/lego/v5/challenge"
	"github.com/go-acme/lego/v5/providers/dns/acmedns"
	"github.com/go-acme/lego/v5/providers/dns/azuredns"
	"github.com/go-acme/lego/v5/providers/dns/baiducloud"
	"github.com/go-acme/lego/v5/providers/dns/cloudflare"
	"github.com/go-acme/lego/v5/providers/dns/digitalocean"
	"github.com/go-acme/lego/v5/providers/dns/dnsupdate"
	"github.com/go-acme/lego/v5/providers/dns/gcloud"
	"github.com/go-acme/lego/v5/providers/dns/godaddy"
	"github.com/go-acme/lego/v5/providers/dns/hetzner"
	"github.com/go-acme/lego/v5/providers/dns/huaweicloud"
	"github.com/go-acme/lego/v5/providers/dns/namecheap"
	"github.com/go-acme/lego/v5/providers/dns/route53"
	"github.com/go-acme/lego/v5/providers/dns/tencentcloud"

	"github.com/localpaas/localpaas/localpaas_app/apperrors"
	"github.com/localpaas/localpaas/localpaas_app/base"
	"github.com/localpaas/localpaas/localpaas_app/entity"
)

func NewDNS01Provider(
	providerKind base.AcmeDnsProvider,
	dnsConfig *entity.AcmeDnsProvider,
) (provider challenge.Provider, err error) {
	switch providerKind {
	case base.AcmeDnsProviderAcmeDNS:
		provider, err = dns01CreateProviderAcmeDNS(dnsConfig.AcmeDNS)
	case base.AcmeDnsProviderAzure:
		provider, err = dns01CreateProviderAzure(dnsConfig.Azure)
	case base.AcmeDnsProviderBaiduCloud:
		provider, err = dns01CreateProviderBaiduCloud(dnsConfig.BaiduCloud)
	case base.AcmeDnsProviderCloudflare:
		provider, err = dns01CreateProviderCloudflare(dnsConfig.Cloudflare)
	case base.AcmeDnsProviderDigitalOcean:
		provider, err = dns01CreateProviderDigitalOcean(dnsConfig.DigitalOcean)
	case base.AcmeDnsProviderGCloud:
		provider, err = dns01CreateProviderGCloud(dnsConfig.GCloud)
	case base.AcmeDnsProviderGoDaddy:
		provider, err = dns01CreateProviderGoDaddy(dnsConfig.GoDaddy)
	case base.AcmeDnsProviderHetzner:
		provider, err = dns01CreateProviderHetzner(dnsConfig.Hetzner)
	case base.AcmeDnsProviderHuaweiCloud:
		provider, err = dns01CreateProviderHuaweiCloud(dnsConfig.HuaweiCloud)
	case base.AcmeDnsProviderNamecheap:
		provider, err = dns01CreateProviderNamecheap(dnsConfig.Namecheap)
	case base.AcmeDnsProviderRFC2136:
		provider, err = dns01CreateProviderRFC2136(dnsConfig.RFC2136)
	case base.AcmeDnsProviderRoute53:
		provider, err = dns01CreateProviderRoute53(dnsConfig.Route53)
	case base.AcmeDnsProviderTencentCloud:
		provider, err = dns01CreateProviderTencentCloud(dnsConfig.TencentCloud)
	}
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderAcmeDNS(dns01 *entity.AcmeDnsProviderAcmeDNS) (challenge.Provider, error) {
	config := acmedns.NewDefaultConfig()
	config.APIBase = dns01.APIBase
	config.AllowList = dns01.AllowList
	config.StoragePath = dns01.StoragePath
	config.StorageBaseURL = dns01.StorageBaseURL

	provider, err := acmedns.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderAzure(dns01 *entity.AcmeDnsProviderAzure) (challenge.Provider, error) {
	config := azuredns.NewDefaultConfig()
	config.ClientID = dns01.ClientID
	config.ClientSecret = dns01.ClientSecret.MustGetPlain()
	config.TenantID = dns01.TenantID
	config.SubscriptionID = dns01.SubscriptionID
	config.ResourceGroup = dns01.ResourceGroupName

	provider, err := azuredns.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderBaiduCloud(dns01 *entity.AcmeDnsProviderBaiduCloud) (challenge.Provider, error) {
	config := baiducloud.NewDefaultConfig()
	config.AccessKeyID = dns01.AccessKey
	config.SecretAccessKey = dns01.SecretKey.MustGetPlain()

	provider, err := baiducloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderCloudflare(dns01 *entity.AcmeDnsProviderCloudflare) (challenge.Provider, error) {
	config := cloudflare.NewDefaultConfig()
	config.AuthToken = dns01.AuthToken.MustGetPlain()

	provider, err := cloudflare.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderDigitalOcean(dns01 *entity.AcmeDnsProviderDigitalOcean) (challenge.Provider, error) {
	config := digitalocean.NewDefaultConfig()
	config.AuthToken = dns01.AuthToken.MustGetPlain()

	provider, err := digitalocean.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderGCloud(dns01 *entity.AcmeDnsProviderGCloud) (challenge.Provider, error) {
	config := gcloud.NewDefaultConfig()
	config.Project = dns01.ProjectID
	config.ImpersonateServiceAccount = dns01.ServiceAccount.MustGetPlain()

	provider, err := gcloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderGoDaddy(dns01 *entity.AcmeDnsProviderGoDaddy) (challenge.Provider, error) {
	config := godaddy.NewDefaultConfig()
	config.APIKey = dns01.APIKey
	config.APISecret = dns01.APISecret.MustGetPlain()

	provider, err := godaddy.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderHetzner(dns01 *entity.AcmeDnsProviderHetzner) (challenge.Provider, error) {
	config := hetzner.NewDefaultConfig()
	config.APIToken = dns01.APIToken.MustGetPlain()

	provider, err := hetzner.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderHuaweiCloud(dns01 *entity.AcmeDnsProviderHuaweiCloud) (challenge.Provider, error) {
	config := huaweicloud.NewDefaultConfig()
	config.AccessKeyID = dns01.AccessKey
	config.SecretAccessKey = dns01.SecretKey.MustGetPlain()
	config.Region = dns01.Region

	provider, err := huaweicloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderNamecheap(dns01 *entity.AcmeDnsProviderNamecheap) (challenge.Provider, error) {
	config := namecheap.NewDefaultConfig()
	config.APIUser = dns01.APIUser
	config.APIKey = dns01.APIKey.MustGetPlain()

	provider, err := namecheap.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderRFC2136(dns01 *entity.AcmeDnsProviderRFC2136) (challenge.Provider, error) {
	config := dnsupdate.NewDefaultConfig()
	config.Nameserver = dns01.Nameserver
	config.TSIGKey = dns01.TSIGKeyName
	config.TSIGSecret = dns01.TSIGSecret.MustGetPlain()
	config.TSIGAlgorithm = dns01.TSIGAlgorithm

	provider, err := dnsupdate.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderRoute53(dns01 *entity.AcmeDnsProviderRoute53) (challenge.Provider, error) {
	config := route53.NewDefaultConfig()
	config.AccessKeyID = dns01.AccessKeyID
	config.SecretAccessKey = dns01.SecretAccessKey.MustGetPlain()
	config.HostedZoneID = dns01.HostedZoneID
	config.Region = dns01.Region

	provider, err := route53.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}

func dns01CreateProviderTencentCloud(dns01 *entity.AcmeDnsProviderTencentCloud) (challenge.Provider, error) {
	config := tencentcloud.NewDefaultConfig()
	config.SecretID = dns01.SecretID
	config.SecretKey = dns01.SecretKey.MustGetPlain()
	config.Region = dns01.Region

	provider, err := tencentcloud.NewDNSProviderConfig(config)
	if err != nil {
		return nil, apperrors.Wrap(err)
	}
	return provider, nil
}
