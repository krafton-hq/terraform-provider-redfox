package api_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func CustomDocumentFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"raw_json": {
			Description:      "Custom Document Spec Raw Json",
			Type:             schema.TypeString,
			Required:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsJSON),
		},
	}
}

func MarshalCustomDocumentSpec(raw []any) (string, error) {
	if raw == nil {
		return "", fmt.Errorf("crd Block Should not be null")
	}

	rawItem := raw[0].(map[string]any)
	return rawItem["raw_json"].(string), nil
}

func UnmarshalCustomDocumentSpec(rawSpec string) ([]any, error) {
	return []any{map[string]any{
		"raw_json": rawSpec,
	}}, nil
}

func AssembleCustomDocument(gvk *idl_common.GroupVersionKindSpec, metadata *idl_common.ObjectMeta, rawSpec string) *documents.CustomDocument {
	return &documents.CustomDocument{
		ApiVersion: gvk.Group + "/" + gvk.Version,
		Kind:       gvk.Kind,
		Metadata:   metadata,
		RawSpec:    rawSpec,
	}
}

func DisassembleCustomDocument(cr *documents.CustomDocument) (*idl_common.GroupVersionKindSpec, *idl_common.ObjectMeta, string) {
	gvk, err := ParseGvk(cr.ApiVersion, cr.Kind)
	// Unreachable
	if err != nil {
		tflog.Error(context.TODO(), err.Error())
	}

	return gvk, cr.Metadata, cr.RawSpec
}
