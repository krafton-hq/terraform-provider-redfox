package provider

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
)

func DataSourceNamespaces() *schema.Resource {
	return &schema.Resource{
		Description: "List RedFox Namespaces, Like https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",

		ReadContext: dataSourceNamespacesRead,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"label_selectors": api_object.LabelSelector(),
			"namespaces": {
				Type:        schema.TypeList,
				Description: "List of Namespaces",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metadata": api_object.ApiObjectMeta(),
						"spec":     api_object.NamespaceDataSourceSpec(),
					},
				},
			},
		},
	}
}

func dataSourceNamespacesRead(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	rawSelector := d.Get("label_selectors").(map[string]any)
	selector := api_object.MarshalLabelSelectors(rawSelector)

	client := rawClient.(redfox_helper.ClientHelper)
	res, err := client.Namespaces().ListNamespaces(ctx, &idl_common.ListObjectReq{LabelSelectors: selector})
	if err != nil {
		return diag.FromErr(err)
	}
	if res.CommonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("List Namespace Failed, status: %v, message: %s", res.CommonRes.Status, res.CommonRes.Message)
	}

	var rawNamespaces []any
	for _, namespace := range res.Namespaces {
		rawMetadata, err := api_object.UnmarshalApiObjectMeta(namespace.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}

		rawSpec, err := api_object.UnmarshalNamespaceSpec(namespace.Spec)
		if err != nil {
			return diag.FromErr(err)
		}

		rawNamespace := map[string]any{
			"metadata": rawMetadata,
			"spec":     rawSpec,
		}
		rawNamespaces = append(rawNamespaces, rawNamespace)
	}

	hashId, err := hashTerraformObjects(res.Namespaces)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(hashId)
	if err = d.Set("namespaces", rawNamespaces); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func hashTerraformObjects(a any) (string, error) {
	buf, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("unmarshal Terraform Object to Json Failed: %v", err.Error())
	}

	hash := sha256.New()
	hash.Write(buf)
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
