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

func NewNatIpOrigin() *template.GenericResource[*documents.NatIpSpec, *documents.NatIp] {
	return &template.GenericResource[*documents.NatIpSpec, *documents.NatIp]{
		ResourceName:       "NatIp",
		ResourceNamePlural: "NatIps",
		Description:        "RedFox NatIp",

		IsNamespaced: true,
		SpecSchema:   api_object.NatIpSpecFields(),
		Timeouts:     &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},

		SpecMarshaller:  api_object.MarshalNatIpSpec,
		SpecUnmarshaler: api_object.UnmarshalNatIpSpec,
		ResAssembler:    api_object.AssembleNatIp,
		ResDisassembler: api_object.DisassembleNatIp,
		GvkOption: template.GvkOption{
			UsePredefined: true,
			GvkResolver: func(client redfox_helper.ClientHelper, raw *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error) {
				return client.NatIpGvk(), nil
			},
		},

		Getter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*documents.NatIp, *idl_common.CommonRes, error) {
			res, err := client.NatIps().GetNatIp(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
			})
			if res == nil {
				return nil, nil, err
			}
			return res.NatIp, res.CommonRes, err
		},
		Lister: func(ctx context.Context, client redfox_helper.ClientHelper, request *idl_common.ListObjectReq) ([]*documents.NatIp, *idl_common.CommonRes, error) {
			res, err := client.NatIps().ListNatIps(ctx, request)
			if res == nil {
				return nil, nil, err
			}
			return res.NatIps, res.CommonRes, err
		},
		Creator: func(ctx context.Context, client redfox_helper.ClientHelper, natIp *documents.NatIp) (*idl_common.CommonRes, error) {
			return client.NatIps().CreateNatIp(ctx, &documents.DesiredNatIpReq{NatIp: natIp})
		},
		Updator: func(ctx context.Context, client redfox_helper.ClientHelper, natIp *documents.NatIp) (*idl_common.CommonRes, error) {
			return client.NatIps().UpdateNatIp(ctx, &documents.DesiredNatIpReq{NatIp: natIp})
		},
		Deleter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*idl_common.CommonRes, error) {
			return client.NatIps().DeleteNatIp(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
			})
		},
	}
}
