package template

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
)

type MarshalResourceSpec[Spec any] func([]any) (Spec, error)
type UnmarshalResourceSpec[Spec any] func(spec Spec) ([]any, error)

type AssembleResource[Spec any, Res any] func(gvk *idl_common.GroupVersionKindSpec, metadata *idl_common.ObjectMeta, spec Spec) Res
type DisassembleResource[Spec any, Res any] func(res Res) (*idl_common.GroupVersionKindSpec, *idl_common.ObjectMeta, Spec)

type GenericResource[Spec any, Res any] struct {
	ResourceName       string
	ResourceNamePlural string
	Description        string

	IsNamespaced bool
	SpecSchema   map[string]*schema.Schema
	Timeouts     *schema.ResourceTimeout

	SpecMarshaller  MarshalResourceSpec[Spec]
	SpecUnmarshaler UnmarshalResourceSpec[Spec]
	ResAssembler    AssembleResource[Spec, Res]
	ResDisassembler DisassembleResource[Spec, Res]
	GvkOption       GvkOption

	Getter  func(context.Context, redfox_helper.ClientHelper, *api_object.ResourceId) (Res, *idl_common.CommonRes, error)
	Lister  func(context.Context, redfox_helper.ClientHelper, *idl_common.ListObjectReq) ([]Res, *idl_common.CommonRes, error)
	Creator func(context.Context, redfox_helper.ClientHelper, Res) (*idl_common.CommonRes, error)
	Updator func(context.Context, redfox_helper.ClientHelper, Res) (*idl_common.CommonRes, error)
	Deleter func(context.Context, redfox_helper.ClientHelper, *api_object.ResourceId) (*idl_common.CommonRes, error)
}

type GvkOption struct {
	UsePredefined bool
	GvkResolver   func(redfox_helper.ClientHelper, *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error)
}

func (r *GenericResource[Spec, Res]) Resource() *schema.Resource {
	return &schema.Resource{
		Description: r.Description,

		ReadContext:   r.ReadContextResource,
		DeleteContext: r.DeleteContext,
		CreateContext: r.CreateUpdateContext,
		UpdateContext: r.CreateUpdateContext,
		Timeouts:      r.Timeouts,

		Schema: map[string]*schema.Schema{
			"api_version": api_object.ApiVersion(r.GvkOption.UsePredefined),
			"kind":        api_object.Kind(r.GvkOption.UsePredefined),
			"metadata":    api_object.ApiObjectMeta(r.IsNamespaced),
			"spec": {
				Description: fmt.Sprintf("%s Spec Block", r.ResourceName),
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem:        &schema.Resource{Schema: r.SpecSchema},
			},
		},
	}
}

func (r *GenericResource[Spec, Res]) DataSource() *schema.Resource {
	return &schema.Resource{
		Description: r.Description,

		ReadContext: r.ReadContextDataSource,
		Timeouts:    r.Timeouts,

		Schema: map[string]*schema.Schema{
			"api_version": api_object.ApiVersion(r.GvkOption.UsePredefined),
			"kind":        api_object.Kind(r.GvkOption.UsePredefined),
			"metadata":    api_object.ApiObjectMeta(r.IsNamespaced),
			"spec": {
				Description: fmt.Sprintf("%s Spec Block", r.ResourceName),
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Resource{Schema: r.SpecSchema},
			},
		},
	}
}

func (r *GenericResource[Spec, Res]) DataSources() *schema.Resource {
	schemeRes := &schema.Resource{
		Description: fmt.Sprintf("List %s", r.ResourceName),

		ReadContext: r.ReadContextDataSources,
		Timeouts:    r.Timeouts,

		Schema: map[string]*schema.Schema{
			"api_version":     api_object.ApiVersion(r.GvkOption.UsePredefined),
			"kind":            api_object.Kind(r.GvkOption.UsePredefined),
			"label_selectors": api_object.LabelSelector(),
			strings.ToLower(r.ResourceNamePlural): {
				Type:        schema.TypeList,
				Description: r.Description,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Description: "RedFox ApiVersion, Same as ...",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"kind": {
							Description: "RedFox Kind, Same as ...",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"metadata": api_object.ApiObjectMeta(r.IsNamespaced),
						"spec": {
							Description: fmt.Sprintf("%s Spec Block", r.ResourceName),
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Resource{Schema: r.SpecSchema},
						},
					},
				},
			},
		},
	}
	if r.IsNamespaced {
		schemeRes.Schema["namespace"] = &schema.Schema{
			Description: "Resource Namespace use only Namespaced Resource, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",
			Type:        schema.TypeString,
			Optional:    true,
		}
	}
	return schemeRes
}

func (r *GenericResource[Spec, Res]) CreateUpdateContext(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	client := rawClient.(redfox_helper.ClientHelper)
	gvk, err := r.GvkOption.GvkResolver(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	rawMetadata := d.Get("metadata").([]any)
	metadata, err := api_object.MarshalApiObjectMeta(rawMetadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec := d.Get("spec").([]any)
	spec, err := r.SpecMarshaller(rawSpec)
	if err != nil {
		return diag.FromErr(err)
	}

	res := r.ResAssembler(gvk, metadata, spec)
	var commonRes *idl_common.CommonRes
	if d.IsNewResource() {
		commonRes, err = r.Creator(ctx, client, res)
	} else {
		commonRes, err = r.Updator(ctx, client, res)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if commonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Create or Update %s Failed, status: %v, message: %s", r.ResourceName, commonRes.Status, commonRes.Message)
	}

	var id *api_object.ResourceId
	if r.IsNamespaced {
		id = api_object.NewResourceIdFull(gvk, metadata.Namespace, metadata.Name)
	} else {
		id = api_object.NewResourceId(gvk, metadata.Name)
	}
	d.SetId(id.String())
	return r.ReadContextResource(ctx, d, rawClient)
}

func (r *GenericResource[Spec, Res]) ReadContextResource(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	id, err := api_object.ParseResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := rawClient.(redfox_helper.ClientHelper)
	resource, commonRes, err := r.Getter(ctx, client, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if commonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Get %s Failed, status: %v, message: %s", r.ResourceName, commonRes.Status, commonRes.Message)
	}

	_, metadata, spec := r.ResDisassembler(resource)

	rawMetadata, err := api_object.UnmarshalApiObjectMeta(metadata)
	if err != nil {
		return diag.FromErr(err)
	}

	rawSpec, err := r.SpecUnmarshaler(spec)
	if err != nil {
		return diag.FromErr(err)
	}

	if r.GvkOption.UsePredefined {
		if err = d.Set("api_version", id.ApiVersion()); err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("kind", id.Gvk.Kind); err != nil {
			return diag.FromErr(err)
		}
	}
	if err = d.Set("metadata", rawMetadata); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("spec", rawSpec); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func (r *GenericResource[Spec, Res]) ReadContextDataSource(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	client := rawClient.(redfox_helper.ClientHelper)
	gvk, err := r.GvkOption.GvkResolver(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	name := d.Get("metadata.0.name").(string)

	var id *api_object.ResourceId
	if r.IsNamespaced {
		namespace := d.Get("metadata.0.namespace").(string)
		id = api_object.NewResourceIdFull(gvk, namespace, name)
	} else {
		id = api_object.NewResourceId(gvk, name)
	}

	d.SetId(id.String())

	return r.ReadContextResource(ctx, d, rawClient)
}

func (r *GenericResource[Spec, Res]) ReadContextDataSources(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	client := rawClient.(redfox_helper.ClientHelper)
	request := &idl_common.ListObjectReq{}
	if !r.GvkOption.UsePredefined {
		gvk, err := r.GvkOption.GvkResolver(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
		request.Gvk = gvk
	}
	if r.IsNamespaced {
		rawNs, ok := d.GetOk("namespace")
		if ok {
			request.Namespace = rawNs.(string)
		}
	}

	rawSelector := d.Get("label_selectors").(map[string]any)
	selector := api_object.MarshalLabelSelectors(rawSelector)
	request.LabelSelectors = selector

	resourceList, commonRes, err := r.Lister(ctx, client, request)
	if err != nil {
		return diag.FromErr(err)
	}
	if commonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("List %s Failed, status: %v, message: %s", r.ResourceName, commonRes.Status, commonRes.Message)
	}

	var rawResources []any
	for _, res := range resourceList {
		gvk, metadata, spec := r.ResDisassembler(res)

		rawMetadata, err := api_object.UnmarshalApiObjectMeta(metadata)
		if err != nil {
			return diag.FromErr(err)
		}

		rawSpec, err := r.SpecUnmarshaler(spec)
		if err != nil {
			return diag.FromErr(err)
		}

		rawResource := map[string]any{
			"api_version": gvk.Group + "/" + gvk.Version,
			"kind":        gvk.Kind,
			"metadata":    rawMetadata,
			"spec":        rawSpec,
		}
		rawResources = append(rawResources, rawResource)
	}

	hashId, err := hashTerraformObjects(resourceList)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(hashId)

	if err = d.Set(strings.ToLower(r.ResourceNamePlural), rawResources); err != nil {
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

func (r *GenericResource[Spec, Res]) DeleteContext(ctx context.Context, d *schema.ResourceData, rawClient interface{}) diag.Diagnostics {
	id, err := api_object.ParseResourceId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	client := rawClient.(redfox_helper.ClientHelper)

	commonRes, err := r.Deleter(ctx, client, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if commonRes.Status != idl_common.ResultCode_SUCCESS {
		return diag.Errorf("Delete %s Failed, status: %v, message: %s", r.ResourceName, commonRes.Status, commonRes.Message)
	}
	return nil
}
