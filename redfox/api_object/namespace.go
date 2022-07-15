package api_object

import (
	"fmt"

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

func NamespaceResourceSpec() *schema.Schema {
	return &schema.Schema{
		Description: "Namespace Spec Block",
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: NamespaceSpecFields(),
		},
	}
}

func NamespaceDataSourceSpec() *schema.Schema {
	return &schema.Schema{
		Description: "Namespace Spec Block",
		Type:        schema.TypeList,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: NamespaceSpecFields(),
		},
	}
}

func MarshalNamespaceSpec(raw []any) (*namespaces.NamespaceSpec, error) {
	if raw == nil {
		return nil, fmt.Errorf("namespace Block Should not be null")
	}

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
