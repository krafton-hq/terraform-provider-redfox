package redfox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func resourceRedfoxClusterStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedfoxClusterStatusApply,
		ReadContext:   resourceRedfoxClusterStatusRead,
		UpdateContext: resourceRedfoxClusterStatusApply,
		DeleteContext: resourceRedfoxClusterStatusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("cluster", false),
			"status": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
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
	}
}

func resourceRedfoxClusterStatusApply(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	status, err := expandClusterStatus(d.Get("status").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := &redfoxV1alpha1.Cluster{
		TypeMeta:   clusterTypeMeta,
		ObjectMeta: metadata,
		Status:     *status,
	}

	tflog.Info(ctx, fmt.Sprintf("Apply %s: %#v", clusterKind.Kind, cluster))

	buf, err := json.Marshal(cluster)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("%s Marshal error: %#v", clusterKind.Kind, err))
		return diag.FromErr(err)
	}
	out, err := conn.MetadataV1alpha1().Clusters(cluster.Namespace).Patch(ctx, cluster.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: defaultFieldManagerName}, "status")
	if err != nil {
		return diag.Errorf(fmt.Sprintf("Failed to create %s: %#v", clusterKind.Kind, err))
	}

	d.SetId(buildId(out.ObjectMeta))

	tflog.Info(ctx, fmt.Sprintf("Submitted new %s: %#v", clusterKind.Kind, out))

	return resourceRedfoxClusterStatusRead(ctx, d, meta)
}

func resourceRedfoxClusterStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	exists, err := resourceRedfoxClusterExists(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		d.SetId("")
		return diag.Diagnostics{}
	}

	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Reading %s %s", clusterKind.Kind, name))
	cluster, err := conn.MetadataV1alpha1().Clusters(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Received error: %#v", err))
		return diag.FromErr(err)
	}
	tflog.Info(ctx, fmt.Sprintf("Received %s: %#v", clusterKind.Kind, cluster))

	err = d.Set("metadata", flattenMetadata(cluster.ObjectMeta, d, meta))
	if err != nil {
		return diag.FromErr(err)
	}

	status, err := flattenClusterStatus(cluster.Status, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("status", status)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceRedfoxClusterStatusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	exists, err := resourceRedfoxClusterExists(ctx, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		d.SetId("")
		return nil
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting %s %s", clusterKind.Kind, name))

	patchs := PatchOperations{&RemoveOperation{Path: "/status"}}
	buf, err := patchs.MarshalJSON()
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = conn.MetadataV1alpha1().Clusters(namespace).Patch(ctx, name, types.JSONPatchType, buf, metav1.PatchOptions{FieldManager: defaultFieldManagerName}, "status")
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("%s %s deleted", clusterKind.Kind, name))

	d.SetId("")
	return nil
}
