package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
)

func ResourceNatIp() *schema.Resource {
	return &schema.Resource{
		Description: "RedFox NatIp Resource, Like https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",

		ReadContext:   resourceNatIpRead,
		DeleteContext: resourceNatIpDelete,
		CreateContext: resourceNatIpCreateUpdate,
		UpdateContext: resourceNatIpCreateUpdate,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metadata": api_object.ApiObjectMeta(),
			"spec":     api_object.NatIpResourceSpec(),
		},
	}
}

func resourceNatIpCreateUpdate(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	rawMetadata := d.Get("metadata").([]any)
	metadata, err := api_object.MarshalApiObjectMeta(rawMetadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec := d.Get("spec").([]any)
	spec, err := api_object.MarshalNatIpSpec(rawSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	client := rawClient.(redfox_helper.ClientHelper)
	gvk := client.NatIpGvk()

	var res *idl_common.CommonRes
	if d.IsNewResource() {
		res, err = client.NatIps().CreateNatIp(ctx, &documents.DesiredNatIpReq{
			NatIp: &documents.NatIp{
				ApiVersion: gvk.GetGroup() + "/" + gvk.GetVersion(),
				Kind:       gvk.GetKind(),
				Metadata:   metadata,
				Spec:       spec,
			},
		})
	} else {
		res, err = client.NatIps().UpdateNatIp(ctx, &documents.DesiredNatIpReq{
			NatIp: &documents.NatIp{
				ApiVersion: gvk.GetGroup() + "/" + gvk.GetVersion(),
				Kind:       gvk.GetKind(),
				Metadata:   metadata,
				Spec:       spec,
			},
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if res.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Create NatIp Failed, status: %v, message: %s", res.Status, res.Message)
	}

	d.SetId(api_object.BuildNamespaceObjectId(metadata.Namespace, metadata.Name))
	return resourceNatIpRead(ctx, d, rawClient)
}

func resourceNatIpRead(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	namespace, name, found := api_object.ParseNamespacedObjectId(d.Id())
	if !found {
		return diag.Errorf("Can't Parse Terraform Object Id to Namespaced Object, id: '%s'", d.Id())
	}

	client := rawClient.(redfox_helper.ClientHelper)
	res, err := client.NatIps().GetNatIp(ctx, &idl_common.SingleObjectReq{
		Name:      name,
		Namespace: &namespace,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if res.CommonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Get NatIp Failed, status: %v, message: %s", res.CommonRes.Status, res.CommonRes.Message)
	}

	rawMetadata, err := api_object.UnmarshalApiObjectMeta(res.NatIp.Metadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec, err := api_object.UnmarshalNatIpSpec(res.NatIp.Spec)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("metadata", rawMetadata); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("spec", rawSpec); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceNatIpDelete(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	namespace, name, found := api_object.ParseNamespacedObjectId(d.Id())
	if !found {
		return diag.Errorf("Can't Parse Terraform Object Id to Namespaced Object, id: '%s'", d.Id())
	}

	client := rawClient.(redfox_helper.ClientHelper)
	res, err := client.NatIps().DeleteNatIp(ctx, &idl_common.SingleObjectReq{
		Name:      name,
		Namespace: &namespace,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if res.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Delete NatIp Failed, status: %v, message: %s", res.Status, res.Message)
	}
	return nil
}
