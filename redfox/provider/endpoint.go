package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/template"
)

func NewEndpointOrigin() *template.GenericResource[*documents.EndpointSpec, *documents.Endpoint] {
	return &template.GenericResource[*documents.EndpointSpec, *documents.Endpoint]{
		ResourceName:       "Endpoint",
		ResourceNamePlural: "Endpoints",
		Description:        "RedFox Endpoint",

		IsNamespaced: true,
		SpecSchema:   api_object.EndpointSpecFields(),
		Timeouts:     &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},

		SpecMarshaller:  api_object.MarshalEndpointSpec,
		SpecUnmarshaler: api_object.UnmarshalEndpointSpec,
		ResAssembler:    api_object.AssembleEndpoint,
		ResDisassembler: api_object.DisassembleEndpoint,
		GvkOption: template.GvkOption{
			UsePredefined: true,
			GvkResolver: func(client redfox_helper.ClientHelper, raw *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error) {
				return client.EndpointGvk(), nil
			},
		},

		Getter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*documents.Endpoint, *idl_common.CommonRes, error) {
			res, err := client.Endpoints().GetEndpoint(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
			})
			if res == nil {
				return nil, nil, err
			}
			return res.Endpoint, res.CommonRes, err
		},
		Lister: func(ctx context.Context, client redfox_helper.ClientHelper, request *idl_common.ListObjectReq) ([]*documents.Endpoint, *idl_common.CommonRes, error) {
			res, err := client.Endpoints().ListEndpoints(ctx, request)
			if res == nil {
				return nil, nil, err
			}
			return res.Endpoints, res.CommonRes, err
		},
		Creator: func(ctx context.Context, client redfox_helper.ClientHelper, endpoint *documents.Endpoint) (*idl_common.CommonRes, error) {
			return client.Endpoints().CreateEndpoint(ctx, &documents.DesiredEndpointReq{Endpoint: endpoint})
		},
		Updator: func(ctx context.Context, client redfox_helper.ClientHelper, endpoint *documents.Endpoint) (*idl_common.CommonRes, error) {
			return client.Endpoints().UpdateEndpoint(ctx, &documents.DesiredEndpointReq{Endpoint: endpoint})
		},
		Deleter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*idl_common.CommonRes, error) {
			return client.Endpoints().DeleteEndpoint(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
			})
		},
	}
}
