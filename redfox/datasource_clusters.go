package redfox

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func dataSourceRedfoxClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRedfoxClustersRead,

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
						"metadata": namespacedMetadataSchema("cluster", false),
						"spec": {
							Type:        schema.TypeList,
							Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_name": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"cluster_region": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"cluster_group": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"service_phase": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"service_tag": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"cluster_engine": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"infra_vendor": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
									"infra_account_id": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
									},
								},
							},
						},
						"status": {
							Type:        schema.TypeList,
							Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"apiserver": {
										Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
										Type:        schema.TypeList,
										Required:    true,
										MinItems:    1,
										MaxItems:    1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"endpoint": {
													Description:      "",
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
												},
												"ca_cert": {
													Description:      "",
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
												},
											},
										},
									},
									"service_account_issuer": {
										Description:      "",
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
									},
									"aws_iam_idps": {
										Description: "",
										Type:        schema.TypeMap,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
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

func dataSourceRedfoxClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	outs, err := conn.MetadataV1alpha1().Clusters(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: kubeGenericSelector.String(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	var attrs []any
	for index, cluster := range outs.Items {
		att := map[string]any{}
		att["metadata"] = flattenMetadata(cluster.ObjectMeta, d, meta, fmt.Sprintf("items.%d.", index))

		spec, err := flattenClusterSpec(cluster.Spec, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		att["spec"] = spec

		status, err := flattenClusterStatus(cluster.Status, d, meta)
		if err != nil {
			return diag.FromErr(err)
		}
		att["status"] = status

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
