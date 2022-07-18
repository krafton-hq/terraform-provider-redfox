package api_object

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func GroupVersionKind() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group": {
				Description: "Object's Group Name, Like Kubernetes ApiGroup",
				Type:        schema.TypeString,
				Required:    true,
			},
			"version": {
				Description: "Object's Version, Like Kubernetes Version",
				Type:        schema.TypeString,
				Required:    true,
			},
			"kind": {
				Description: "Object's Kind Name, Like Kubernetes Kind",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
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

func UnmarshalGvks(gvks []*idl_common.GroupVersionKindSpec) (*schema.Set, error) {
	if gvks == nil {
		return nil, fmt.Errorf("groupVersionKind Block Should not be null")
	}

	var raws []any
	for _, gvk := range gvks {
		buf, err := json.Marshal(gvk)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Go Struct to Json Failed: %v", err.Error())
		}

		raw := map[string]any{}
		err = json.Unmarshal(buf, &raw)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Map Failed: %v", err.Error())
		}
		raws = append(raws, raw)
	}

	return schema.NewSet(schema.HashResource(GroupVersionKind()), raws), nil
}

func ParseGvk(apiVersion string, kind string) (*idl_common.GroupVersionKindSpec, error) {
	group, version, found := strings.Cut(apiVersion, "/")
	if !found {
		return nil, fmt.Errorf("can't Parse ApiVersion to ApiGroup and Version, api_version: '%s'", apiVersion)
	}

	return &idl_common.GroupVersionKindSpec{
		Group:   group,
		Version: version,
		Kind:    kind,
	}, nil
}
