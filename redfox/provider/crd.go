package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/crds"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/template"
)

func NewCrdOrigin() *template.GenericResource[*crds.CustomResourceDefinitionSpec, *crds.CustomResourceDefinition] {
	return &template.GenericResource[*crds.CustomResourceDefinitionSpec, *crds.CustomResourceDefinition]{
		ResourceName:       "Crd",
		ResourceNamePlural: "Crds",
		Description:        "RedFox Crd",

		IsNamespaced: false,
		SpecSchema:   api_object.CrdSpecFields(),
		Timeouts:     &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},

		SpecMarshaller:  api_object.MarshalCrdSpec,
		SpecUnmarshaler: api_object.UnmarshalCrdSpec,
		ResAssembler:    api_object.AssembleCrd,
		ResDisassembler: api_object.DisassembleCrd,
		GvkOption: template.GvkOption{
			UsePredefined: true,
			GvkResolver: func(client redfox_helper.ClientHelper, raw *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error) {
				return client.NamespaceGvk(), nil
			},
		},

		Getter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*crds.CustomResourceDefinition, *idl_common.CommonRes, error) {
			res, err := client.Crds().GetCustomResourceDefinition(ctx, &idl_common.SingleObjectReq{
				Name: id.Name,
			})
			if res == nil {
				return nil, nil, err
			}
			return res.Crd, res.CommonRes, err
		},
		Lister: func(ctx context.Context, client redfox_helper.ClientHelper, request *idl_common.ListObjectReq) ([]*crds.CustomResourceDefinition, *idl_common.CommonRes, error) {
			res, err := client.Crds().ListCustomResourceDefinitions(ctx, request)
			if res == nil {
				return nil, nil, err
			}
			return res.Crds, res.CommonRes, err
		},
		Creator: func(ctx context.Context, client redfox_helper.ClientHelper, crd *crds.CustomResourceDefinition) (*idl_common.CommonRes, error) {
			return client.Crds().CreateCustomResourceDefinition(ctx, &crds.CreateCustomResourceDefinitionReq{Crd: crd})
		},
		Updator: func(ctx context.Context, client redfox_helper.ClientHelper, crd *crds.CustomResourceDefinition) (*idl_common.CommonRes, error) {
			return client.Crds().UpdateCustomResourceDefinition(ctx, &crds.UpdateCustomResourceDefinitionReq{Crd: crd})
		},
		Deleter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*idl_common.CommonRes, error) {
			return client.Crds().DeleteCustomResourceDefinition(ctx, &idl_common.SingleObjectReq{
				Name: id.Name,
			})
		},
	}
}
