package api_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/crds"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func CrdSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"gvk": {
			Description: "GroupVersionKind Spec Blocks",
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			MaxItems:    1,
			ConfigMode:  schema.SchemaConfigModeBlock,
			Elem:        GroupVersionKind(),
		},
	}
}

func MarshalCrdSpec(raw []any) (*crds.CustomResourceDefinitionSpec, error) {
	if raw == nil {
		return nil, fmt.Errorf("crd Block Should not be null")
	}

	rawItem := raw[0].(map[string]interface{})
	var gvks []*idl_common.GroupVersionKindSpec
	if rawGvk, found := rawItem["gvk"]; found {
		var err error
		gvksSet := rawGvk.(*schema.Set)
		gvks, err = MarshalGvks(gvksSet)
		if err != nil {
			return nil, err
		}
	}
	spec := &crds.CustomResourceDefinitionSpec{
		Gvk: gvks[0],
	}
	return spec, nil
}

func UnmarshalCrdSpec(spec *crds.CustomResourceDefinitionSpec) ([]any, error) {
	if spec == nil {
		return nil, fmt.Errorf("crd Block Should not be null")
	}

	rawGvks, err := UnmarshalGvks([]*idl_common.GroupVersionKindSpec{spec.GetGvk()})
	if err != nil {
		return nil, err
	}

	rawItem := map[string]any{
		"gvk": rawGvks,
	}

	return []any{rawItem}, nil
}

func AssembleCrd(gvk *idl_common.GroupVersionKindSpec, metadata *idl_common.ObjectMeta, spec *crds.CustomResourceDefinitionSpec) *crds.CustomResourceDefinition {
	return &crds.CustomResourceDefinition{
		ApiVersion: gvk.Group + "/" + gvk.Version,
		Kind:       gvk.Kind,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func DisassembleCrd(crd *crds.CustomResourceDefinition) (*idl_common.GroupVersionKindSpec, *idl_common.ObjectMeta, *crds.CustomResourceDefinitionSpec) {
	gvk, err := ParseGvk(crd.ApiVersion, crd.Kind)
	// Unreachable
	if err != nil {
		tflog.Error(context.TODO(), err.Error())
	}

	return gvk, crd.Metadata, crd.Spec
}
