package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/apis/namespaces"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
)

func ResourceNamespace() *schema.Resource {
	return &schema.Resource{
		Description: "RedFox Namespace Resource, Like https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",

		ReadContext:   resourceNamespaceRead,
		DeleteContext: resourceNamespaceDelete,
		CreateContext: resourceNamespaceCreateUpdate,
		UpdateContext: resourceNamespaceCreateUpdate,
		Timeouts: &schema.ResourceTimeout{
			Default: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"metadata": api_object.ApiObjectMeta(),
			"spec":     api_object.NamespaceResourceSpec(),
		},
	}
}

func resourceNamespaceCreateUpdate(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	rawMetadata := d.Get("metadata").([]any)
	metadata, err := api_object.MarshalApiObjectMeta(rawMetadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec := d.Get("spec").([]any)
	spec, err := api_object.MarshalNamespaceSpec(rawSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	client := rawClient.(redfox_helper.ClientHelper)
	nsGvk := client.NamespaceGvk()

	var res *idl_common.CommonRes
	if d.IsNewResource() {
		res, err = client.Namespaces().CreateNamespace(ctx, &namespaces.CreateNamespaceReq{
			Namespace: &namespaces.Namespace{
				ApiVersion: nsGvk.GetGroup() + "/" + nsGvk.GetVersion(),
				Kind:       nsGvk.GetKind(),
				Metadata:   metadata,
				Spec:       spec,
			},
		})
	} else {
		res, err = client.Namespaces().UpdateNamespace(ctx, &namespaces.UpdateNamespaceReq{
			Namespace: &namespaces.Namespace{
				ApiVersion: nsGvk.GetGroup() + "/" + nsGvk.GetVersion(),
				Kind:       nsGvk.GetKind(),
				Metadata:   metadata,
				Spec:       spec,
			},
		})
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if res.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Create Namespace Failed, status: %v, message: %s", res.Status, res.Message)
	}

	d.SetId(api_object.BuildClusterObjectId(metadata.Name))
	return resourceNamespaceRead(ctx, d, rawClient)
}

func resourceNamespaceRead(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	name := api_object.ParseClusterObjectId(d.Id())

	client := rawClient.(redfox_helper.ClientHelper)
	res, err := client.Namespaces().GetNamespace(ctx, &idl_common.SingleObjectReq{
		Name: name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if res.CommonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Create Namespace Failed, status: %v, message: %s", res.CommonRes.Status, res.CommonRes.Message)
	}

	rawMetadata, err := api_object.UnmarshalApiObjectMeta(res.Namespace.Metadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec, err := api_object.UnmarshalNamespaceSpec(res.Namespace.Spec)
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

func resourceNamespaceDelete(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	name := api_object.ParseClusterObjectId(d.Id())

	client := rawClient.(redfox_helper.ClientHelper)
	res, err := client.Namespaces().DeleteNamespaces(ctx, &idl_common.SingleObjectReq{
		Name: name,
	})
	if err != nil {
		return diag.FromErr(err)
	}
	if res.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Delete Namespace Failed, status: %v, message: %s", res.Status, res.Message)
	}
	return nil
}
