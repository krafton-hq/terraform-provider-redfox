package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/api_object"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/redfox_helper"
	"github.com/krafton-hq/terraform-provider-redfox/redfox/template"
)

func NewCustomDocumentOrigin() *template.GenericResource[string, *documents.CustomDocument] {
	return &template.GenericResource[string, *documents.CustomDocument]{
		ResourceName:       "CustomDocument",
		ResourceNamePlural: "CustomDocuments",
		Description:        "RedFox CustomDocument",

		IsNamespaced: true,
		SpecSchema:   api_object.CustomDocumentFields(),
		Timeouts:     &schema.ResourceTimeout{Default: schema.DefaultTimeout(1 * time.Minute)},

		SpecMarshaller:  api_object.MarshalCustomDocumentSpec,
		SpecUnmarshaler: api_object.UnmarshalCustomDocumentSpec,
		ResAssembler:    api_object.AssembleCustomDocument,
		ResDisassembler: api_object.DisassembleCustomDocument,
		GvkOption: template.GvkOption{
			UsePredefined: false,
			GvkResolver: func(client redfox_helper.ClientHelper, d *schema.ResourceData) (*idl_common.GroupVersionKindSpec, error) {
				apiVersion := d.Get("api_version").(string)
				kind := d.Get("kind").(string)
				gvk, err := api_object.ParseGvk(apiVersion, kind)
				if err != nil {
					return nil, err
				}

				res, err := client.RawClient().ListApiResources(context.TODO(), &idl_common.CommonReq{})
				if err != nil {
					return nil, err
				}
				if res.CommonRes.Status != idl_common.ResultCode_SUCCESS {
					return nil, fmt.Errorf("ListApiResourcesFailed, status: %v, message: %s", res.CommonRes.Status, res.CommonRes.Message)
				}

				for _, apiResource := range res.ApiResources {
					if apiResource.Gvk.Group == gvk.Group && apiResource.Gvk.Kind == gvk.Kind && apiResource.Gvk.Version == gvk.Version {
						return gvk, nil
					}
				}

				return nil, fmt.Errorf("not Found Compatible Gvk %v", res.ApiResources)
			},
		},

		Getter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*documents.CustomDocument, *idl_common.CommonRes, error) {
			res, err := client.CustomDocuments().GetCustomDocument(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
				Gvk:       id.Gvk,
			})
			if res == nil {
				return nil, nil, err
			}
			return res.Document, res.CommonRes, err
		},
		Lister: func(ctx context.Context, client redfox_helper.ClientHelper, request *idl_common.ListObjectReq) ([]*documents.CustomDocument, *idl_common.CommonRes, error) {
			res, err := client.CustomDocuments().ListCustomDocuments(ctx, request)
			if res == nil {
				return nil, nil, err
			}
			return res.Documents, res.CommonRes, err
		},
		Creator: func(ctx context.Context, client redfox_helper.ClientHelper, cr *documents.CustomDocument) (*idl_common.CommonRes, error) {
			return client.CustomDocuments().CreateCustomDocument(ctx, &documents.DesiredCustomDocumentReq{Document: cr})
		},
		Updator: func(ctx context.Context, client redfox_helper.ClientHelper, cr *documents.CustomDocument) (*idl_common.CommonRes, error) {
			return client.CustomDocuments().UpdateCustomDocument(ctx, &documents.DesiredCustomDocumentReq{Document: cr})
		},
		Deleter: func(ctx context.Context, client redfox_helper.ClientHelper, id *api_object.ResourceId) (*idl_common.CommonRes, error) {
			return client.CustomDocuments().DeleteCustomDocument(ctx, &idl_common.SingleObjectReq{
				Name:      id.Name,
				Namespace: &id.Namespace,
				Gvk:       id.Gvk,
			})
		},
	}
}
