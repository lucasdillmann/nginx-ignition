package letsencrypt

import (
	"context"

	"github.com/go-acme/lego/v4/challenge"

	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/alibaba"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/aws"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/azion"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/azure"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/baiducloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/bunny"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/civo"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/cloudflare"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/cloudns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/conoha"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/constellix"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/cpanel"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/desec"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/digitalocean"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/directadmin"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/dnsimple"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/dnsmadeeasy"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/dreamhost"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/duckdns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/dyn"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/dyndnsfree"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/easydns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/exoscale"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/freemyip"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/gandi"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/gcp"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/godaddy"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/hetzner"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/huaweicloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/hurricane"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/ibmcloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/iij"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/internetbs" //nolint:misspell
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/inwx"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/ionos"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/linode"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/loopia"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/namecheap"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/namedotcom"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/namesilo"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/netcup"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/netlify"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/ns1"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/oraclecloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/ovh"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/pdns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/plesk"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/porkbun"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/rackspace"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/sakuracloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/scaleway"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/tencentcloud"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/ultradns"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/vercel"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/vultr"
	"dillmann.com.br/nginx-ignition/certificate/letsencrypt/dns/yandex"
	"dillmann.com.br/nginx-ignition/core/common/core_error"
)

var (
	providers = []dns.Provider{
		&alibaba.Provider{},
		&azion.Provider{},
		&aws.Provider{},
		&azure.Provider{},
		&baiducloud.Provider{},
		&bunny.Provider{},
		&civo.Provider{},
		&cloudflare.Provider{},
		&cloudns.Provider{},
		&conoha.Provider{},
		&constellix.Provider{},
		&cpanel.Provider{},
		&desec.Provider{},
		&digitalocean.Provider{},
		&directadmin.Provider{},
		&dnsmadeeasy.Provider{},
		&dnsimple.Provider{},
		&dreamhost.Provider{},
		&duckdns.Provider{},
		&dyn.Provider{},
		&dyndnsfree.Provider{},
		&easydns.Provider{},
		&exoscale.Provider{},
		&freemyip.Provider{},
		&gandi.Provider{},
		&gcp.Provider{},
		&godaddy.Provider{},
		&hetzner.Provider{},
		&huaweicloud.Provider{},
		&hurricane.Provider{},
		&ibmcloud.Provider{},
		&iij.Provider{},
		&inwx.Provider{},
		&ionos.Provider{},
		&internetbs.Provider{}, //nolint:misspell
		&linode.Provider{},
		&loopia.Provider{},
		&namecheap.Provider{},
		&namedotcom.Provider{},
		&namesilo.Provider{},
		&netlify.Provider{},
		&netcup.Provider{},
		&ns1.Provider{},
		&oraclecloud.Provider{},
		&ovh.Provider{},
		&pdns.Provider{},
		&plesk.Provider{},
		&porkbun.Provider{},
		&rackspace.Provider{},
		&sakuracloud.Provider{},
		&scaleway.Provider{},
		&tencentcloud.Provider{},
		&ultradns.Provider{},
		&vercel.Provider{},
		&vultr.Provider{},
		&yandex.Provider{},
	}
)

func resolveProviderChallenge(ctx context.Context, domainNames []string, parameters map[string]any) (challenge.Provider, error) {
	providerId, _ := parameters[dnsProvider.ID].(string)

	for _, provider := range providers {
		if provider.ID() == providerId {
			return provider.ChallengeProvider(ctx, domainNames, parameters)
		}
	}

	return nil, core_error.New("Unknown DNS provider", true)
}
