package api_object

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/apis/namespaces"
)

func NamespaceSpec() *schema.Schema {
	return &schema.Schema{
		Description: "Namespace Spec Block",
		Type:        schema.TypeSet,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		ConfigMode:  schema.SchemaConfigModeBlock,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"api_objects": {
					Description: "GroupVersionKind Spec Blocks",
					Type:        schema.TypeSet,
					Optional:    true,
					MinItems:    0,
					ConfigMode:  schema.SchemaConfigModeBlock,
					Elem:        GroupVersionKind(),
				},
			},
		},
	}
}

func MarshalNamespaceSpec(set *schema.Set) (*namespaces.NamespaceSpec, error) {
	if set == nil {
		return nil, fmt.Errorf("namespace Block Should not be null")
	}

	if set.Len() != 1 {
		return &namespaces.NamespaceSpec{}, nil
	}

	rawItem := set.List()[0].(map[string]interface{})
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

func MarshalGvks(set *schema.Set) ([]*idl_common.GroupVersionKindSpec, error) {
	if set == nil {
		return nil, fmt.Errorf("groupVersionKind Block Should not be null")
	}

	var gvks []*idl_common.GroupVersionKindSpec
	for _, rawGvk := range set.List() {
		buf, err := json.Marshal(rawGvk)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Terraform Object to Json Failed: %v", err.Error())
		}

		gvk := &idl_common.GroupVersionKindSpec{}
		err = json.Unmarshal(buf, gvk)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Struct Failed: %v", err.Error())
		}
		gvks = append(gvks, gvk)
	}

	return gvks, nil
}
