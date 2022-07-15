package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
)

func DataSourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description: "Get exist RedFox Namespace, Like https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",

		ReadContext: dataSourceNamespaceRead,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metadata": api_object.ApiObjectMeta(),
			"spec":     api_object.NamespaceDataSourceSpec(),
		},
	}
}

func dataSourceNamespaceRead(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	d.SetId(d.Get("metadata.0.name").(string))

	return resourceNamespaceRead(ctx, d, rawClient)
}
