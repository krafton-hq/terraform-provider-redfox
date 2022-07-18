package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/apis/namespaces"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/template"
)

func NewNamespaceOrigin() *template.GenericResource[*namespaces.NamespaceSpec, *namespaces.Namespace] {
	return &template.GenericResource[*namespaces.NamespaceSpec, *namespaces.Namespace]{
		ResourceName:       "Namespace",
		ResourceNamePlural: "Namespaces",
		Description:        "RedFox Namespace",

		IsNamespaced: false,
		SpecSchema:   api_object.NamespaceSpecFields(),
		Timeouts:     &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},

		SpecMarshaller:  api_object.MarshalNamespaceSpec,
		SpecUnmarshaler: api_object.UnmarshalNamespaceSpec,
		ResAssembler:    api_object.AssembleNamespace,
		ResDisassembler: api_object.DisassembleNamespace,
		GvkOption: template.GvkOption{
			UsePredefined: true,
			GvkResolver: func(client redfox_helper.ClientHelper, raw *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error) {
				return client.NamespaceGvk(), nil
			},
		},

		Getter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*namespaces.Namespace, *idl_common.CommonRes, error) {
			res, err := client.Namespaces().GetNamespace(ctx, &idl_common.SingleObjectReq{
				Name: id.Name,
			})
			if res == nil {
				return nil, nil, err
			}
			return res.Namespace, res.CommonRes, err
		},
		Lister: func(ctx context.Context, client redfox_helper.ClientHelper, request *idl_common.ListObjectReq) ([]*namespaces.Namespace, *idl_common.CommonRes, error) {
			res, err := client.Namespaces().ListNamespaces(ctx, request)
			if res == nil {
				return nil, nil, err
			}
			return res.Namespaces, res.CommonRes, err
		},
		Creator: func(ctx context.Context, client redfox_helper.ClientHelper, namespace *namespaces.Namespace) (*idl_common.CommonRes, error) {
			return client.Namespaces().CreateNamespace(ctx, &namespaces.CreateNamespaceReq{Namespace: namespace})
		},
		Updator: func(ctx context.Context, client redfox_helper.ClientHelper, namespace *namespaces.Namespace) (*idl_common.CommonRes, error) {
			return client.Namespaces().UpdateNamespace(ctx, &namespaces.UpdateNamespaceReq{Namespace: namespace})
		},
		Deleter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*idl_common.CommonRes, error) {
			return client.Namespaces().DeleteNamespaces(ctx, &idl_common.SingleObjectReq{
				Name: id.Name,
			})
		},
	}
}
