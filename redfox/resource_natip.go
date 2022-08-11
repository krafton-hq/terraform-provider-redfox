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

func resourceRedfoxNatIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRedfoxNatIpApply,
		ReadContext:   resourceRedfoxNatIpRead,
		UpdateContext: resourceRedfoxNatIpApply,
		DeleteContext: resourceRedfoxNatIpDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(30 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"metadata": namespacedMetadataSchema("natip", true),
			"spec": {
				Type:        schema.TypeList,
				Description: "Spec defines the specification of the desired behavior of the deployment. More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.9/#deployment-v1-apps",
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
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

var natipKind = kubeSchema.GroupVersionKind{Group: "metadata.sbx-central.io", Version: "v1alpha1", Kind: "NatIp"}
var natipTypeMeta = metav1.TypeMeta{
	Kind:       natipKind.Kind,
	APIVersion: natipKind.GroupVersion().String(),
}

func resourceRedfoxNatIpApply(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	spec, err := expandNatIpSpec(d.Get("spec").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	natIp := &redfoxV1alpha1.NatIp{
		TypeMeta:   natipTypeMeta,
		ObjectMeta: metadata,
		Spec:       *spec,
	}

	tflog.Info(ctx, fmt.Sprintf("Creating new NatIp: %#v", natIp))

	buf, err := json.Marshal(natIp)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("NatIp Marshal error: %#v", err))
		return diag.FromErr(err)
	}
	out, err := conn.MetadataV1alpha1().NatIps(natIp.Namespace).Patch(ctx, natIp.Name, types.ApplyPatchType, buf, metav1.PatchOptions{FieldManager: defaultFieldManagerName})
	if err != nil {
		return diag.Errorf("Failed to create NatIp: %s", err)
	}

	d.SetId(buildId(out.ObjectMeta))

	tflog.Info(ctx, fmt.Sprintf("Submitted new NatIp: %#v", out))

	return resourceRedfoxNatIpRead(ctx, d, meta)
}

func resourceRedfoxNatIpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	exists, err := resourceRedfoxNatIpExists(ctx, d, meta)
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

	tflog.Info(ctx, fmt.Sprintf("Reading deployment %s", name))
	natIp, err := conn.MetadataV1alpha1().NatIps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Received error: %#v", err))
		return diag.FromErr(err)
	}
	tflog.Info(ctx, fmt.Sprintf("Received NatIp: %#v", natIp))

	err = d.Set("metadata", flattenMetadata(natIp.ObjectMeta, d, meta))
	if err != nil {
		return diag.FromErr(err)
	}

	spec, err := flattenNatIpSpec(natIp.Spec, d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("spec", spec)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceRedfoxNatIpExists(ctx context.Context, d *schema.ResourceData, meta interface{}) (bool, error) {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return false, err
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return false, err
	}

	tflog.Info(ctx, fmt.Sprintf("Checking deployment %s", name))
	_, err = conn.MetadataV1alpha1().NatIps(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			return false, nil
		}
		tflog.Debug(ctx, fmt.Sprintf("Received error: %#v", err))
	}
	return true, err
}

func resourceRedfoxNatIpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := meta.(KubeClientsets).RedfoxClient()
	if err != nil {
		return diag.FromErr(err)
	}

	namespace, name, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Deleting deployment: %#v", name))

	err = conn.MetadataV1alpha1().NatIps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return diag.FromErr(err)
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err := conn.MetadataV1alpha1().NatIps(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
				return nil
			}
			return resource.NonRetryableError(err)
		}

		e := fmt.Errorf("NatIp (%s) still exists", d.Id())
		return resource.RetryableError(e)
	})
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Deployment %s deleted", name))

	d.SetId("")
	return nil
}
