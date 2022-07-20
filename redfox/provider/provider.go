package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client_sdk "github.com/krafton-hq/red-fox/client-sdk"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"google.golang.org/grpc"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	natIpOrigin := NewNatIpOrigin()
	namespaceOrigin := NewNamespaceOrigin()
	crdOrigin := NewCrdOrigin()
	endpointOrigin := NewEndpointOrigin()
	customDocumentOrigin := NewCustomDocumentOrigin()

	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				"redfox_natip":           natIpOrigin.DataSource(),
				"redfox_natips":          natIpOrigin.DataSources(),
				"redfox_namespace":       namespaceOrigin.DataSource(),
				"redfox_namespaces":      namespaceOrigin.DataSources(),
				"redfox_crd":             crdOrigin.DataSource(),
				"redfox_crds":            crdOrigin.DataSources(),
				"redfox_endpoint":        endpointOrigin.DataSource(),
				"redfox_endpoints":       endpointOrigin.DataSources(),
				"redfox_customdocument":  customDocumentOrigin.DataSource(),
				"redfox_customdocuments": customDocumentOrigin.DataSources(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"redfox_natip":          natIpOrigin.Resource(),
				"redfox_namespace":      namespaceOrigin.Resource(),
				"redfox_crd":            crdOrigin.Resource(),
				"redfox_endpoint":       endpointOrigin.Resource(),
				"redfox_customdocument": customDocumentOrigin.Resource(),
			},
			Schema: map[string]*schema.Schema{
				"address": {
					Description: "Red-Fox Grpc Address, <domain>:<port>",
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("REDFOX_ADDR", ""),
				},
				"use_tls": {
					Description: "Flag of Use Tls Connection",
					Type:        schema.TypeBool,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("REDFOX_GRPC_WITH_TLS", true),
				},
			},
		}

		p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return providerConfigure(ctx, d, p.TerraformVersion, version)
		}
		return p
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, terraformVersion string, providerVersion string) (interface{}, diag.Diagnostics) {
	address := d.Get("address").(string)
	useTls := d.Get("use_tls").(bool)

	config := client_sdk.DefaultConfig()
	config.DialOptions = []grpc.DialOption{
		grpc.WithUserAgent(fmt.Sprintf("Terraform/%s RedFoxProvider/%s", terraformVersion, providerVersion)),
	}
	config.WithTls = useTls
	config.GrpcEndpoint = address
	client, err := client_sdk.NewClient(config)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	clientHelper, err := redfox_helper.NewClient(ctx, client)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return clientHelper, nil
}
