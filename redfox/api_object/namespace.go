package api_object

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/apis/namespaces"
)

func NamespaceSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_objects": {
			Description: "GroupVersionKind Spec Blocks",
			Type:        schema.TypeSet,
			Optional:    true,
			MinItems:    0,
			ConfigMode:  schema.SchemaConfigModeBlock,
			Elem:        GroupVersionKind(),
		},
	}
}

func MarshalNamespaceSpec(raw []any) (*namespaces.NamespaceSpec, error) {
	if raw == nil {
		return nil, fmt.Errorf("namespace Block Should not be null")
	}

	// Namespace Spec is Optional
	if len(raw) != 1 {
		return &namespaces.NamespaceSpec{}, nil
	}

	rawItem := raw[0].(map[string]interface{})
	var gvks []*idl_common.GroupVersionKindSpec
	if rawApiObjects, found := rawItem["api_objects"]; found {
		var err error
		gvksSet := rawApiObjects.(*schema.Set)
		gvks, err = MarshalGvks(gvksSet)
		if err != nil {
			return nil, err
		}
	}
	spec := &namespaces.NamespaceSpec{
		ApiObjects: gvks,
	}
	return spec, nil
}

func UnmarshalNamespaceSpec(spec *namespaces.NamespaceSpec) ([]any, error) {
	if spec == nil {
		return nil, fmt.Errorf("namespace Block Should not be null")
	}

	rawGvks, err := UnmarshalGvks(spec.GetApiObjects())
	if err != nil {
		return nil, err
	}

	rawItem := map[string]any{
		"api_objects": rawGvks,
	}

	return []any{rawItem}, nil
}

func AssembleNamespace(gvk *idl_common.GroupVersionKindSpec, metadata *idl_common.ObjectMeta, spec *namespaces.NamespaceSpec) *namespaces.Namespace {
	return &namespaces.Namespace{
		ApiVersion: gvk.Group + "/" + gvk.Version,
		Kind:       gvk.Kind,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func DisassembleNamespace(namespace *namespaces.Namespace) (*idl_common.GroupVersionKindSpec, *idl_common.ObjectMeta, *namespaces.NamespaceSpec) {
	gvk, err := ParseGvk(namespace.ApiVersion, namespace.Kind)
	// Unreachable
	if err != nil {
		tflog.Error(context.TODO(), err.Error())
	}

	return gvk, namespace.Metadata, namespace.Spec
}
