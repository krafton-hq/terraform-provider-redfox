package redfox

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	redfoxV1alpha1 "github.com/krafton-hq/red-fox/pkg/apis/redfox/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeSchema "k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

func resourceRedfoxCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedfoxClusterApply,
		ReadContext:   resourceRedfoxClusterRead,
		UpdateContext: resourceRedfoxClusterApply,
		DeleteContext: resourceRedfoxClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("cluster", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
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
		},
	}
}

var clusterKind = kubeSchema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "Cluster"}
var clusterTypeMeta = metav1.TypeMeta{
	Kind:       clusterKind.Kind,
	APIVersion: clusterKind.GroupVersion().String(),
}

func resourceRedfoxClusterApply(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandClusterSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	cluster := &redfoxV1alpha1.Cluster{
		TypeMeta:   clusterTypeMeta,
		ObjectMeta: metadata,
		Spec:       *spec,
	}

	tflog.Info(ctx, fmt.Sprintf("Apply %s: %#v", clusterKind.Kind, cluster))

	buf, err := json.Marshal(cluster)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("%s Marshal error: %#v", clusterKind.Kind, err))
		return diag.FromErr(err)
	}
	out, err := conn.MetadataV1alpha1().Clusters(cluster.Namespace).Patch(ctx, cluster.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: defaultFieldManagerName})
	if err != nil {
		return diag.Errorf(fmt.Sprintf("Failed to create %s: %#v", clusterKind.Kind, err))
	}

	d.SetId(buildId(out.ObjectMeta))

	tflog.Info(ctx, fmt.Sprintf("Submitted new %s: %#v", clusterKind.Kind, out))

	return resourceRedfoxClusterRead(ctx, d, meta)
}

func resourceRedfoxClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	spec, err := flattenClusterSpec(cluster.Spec, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("spec", spec)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceRedfoxClusterExists(ctx context.Context, d *schema.ResourceData, meta interface{}) (bool, error) {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return false, err
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	tflog.Info(ctx, fmt.Sprintf("Reading %s %s", clusterKind.Kind, name))
	_, err = conn.MetadataV1alpha1().Clusters(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			return false, nil
		}
		tflog.Debug(ctx, fmt.Sprintf("Received error: %#v", err))
	}
	return true, err
}

func resourceRedfoxClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting %s %s", clusterKind.Kind, name))

	err = conn.MetadataV1alpha1().Clusters(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.MetadataV1alpha1().Clusters(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		e := fmt.Errorf("%s (%s) still exists", clusterKind.Kind, d.Id())
		return resource.RetryableError(e)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("%s %s deleted", clusterKind.Kind, name))

	d.SetId("")
	return nil
}
