package redfox

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceRedfoxNatIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedfoxNatIpRead,

		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("natip", false),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_type": {
							Description:  "IP Type, Can be either IPv4 or IPv6",
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{string(redfoxV1alpha1.Ipv4), string(redfoxV1alpha1.Ipv6)}, false),
						},
						"cidrs": {
							Description: "Classless Inter-Domain Routing notated IP List, Must be end '/<bits>'",
							Type:        schema.TypeList,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceRedfoxNatIpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	om := metav1.ObjectMeta{
		Namespace: d.Get("metadata.0.namespace").(string),
		Name:      d.Get("metadata.0.name").(string),
	}
	d.SetId(buildId(om))

	return resourceRedfoxNatIpRead(ctx, d, meta)
}
