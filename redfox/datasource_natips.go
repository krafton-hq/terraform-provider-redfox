package redfox

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceRedfoxNatIps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedfoxNatIpsRead,

		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Description: fmt.Sprintf("Namespace defines the space within which name of the %s must be unique.", "natip"),
				Optional:    true,
				Default:     "",
			},
			"selector": {
				Type:        schema.TypeList,
				Description: "A list of selectors which will be used to find ClusterRoles and create the rules.",
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: labelSelectorFields(true),
				},
			},
			"items": {
				Type:        schema.TypeList,
				Description: "List of NatIps",
				Computed:    true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func dataSourceRedfoxNatIpsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}
	namespace := d.Get("namespace").(string)

	labelSelector := expandLabelSelector(d.Get("selector").([]any))
	kubeGenericSelector, err := metav1.LabelSelectorAsSelector(labelSelector)
	if err != nil {
		tflog.Info(ctx, fmt.Sprintf("Convert LabelSelector to KubeGenericSelector Failed: %#v", err))
		return diag.FromErr(err)
	}

	outs, err := conn.MetadataV1alpha1().NatIps(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: kubeGenericSelector.String(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var attrs []any
	for index, natIp := range outs.Items {
		att := map[string]any{}
		att["metadata"] = flattenMetadata(natIp.ObjectMeta, d, meta, fmt.Sprintf("items.%d.", index))

		spec, err := flattenNatIpSpec(natIp.Spec, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		att["spec"] = spec

		attrs = append(attrs, att)
	}
	err = d.Set("items", attrs)
	if err != nil {
		return diag.FromErr(err)
	}

	hashId, err := hashTerraformObjects(outs.Items)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashId)

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
