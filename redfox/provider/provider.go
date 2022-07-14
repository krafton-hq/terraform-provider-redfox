package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	client_sdk "github.com/krafton-hq/red-fox/client-sdk"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/resources"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			DataSourcesMap: map[string]*schema.Resource{
				//"scaffolding_data_source": dataSourceScaffolding(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"redfox_namespace": resources.ResourceNamespace(),
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
