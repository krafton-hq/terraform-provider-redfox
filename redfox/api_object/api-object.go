package api_object

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func ApiObjectMeta() *schema.Schema {
	return &schema.Schema{
		Description: "Api-Object Metadata Block",
		Type:        schema.TypeSet,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		ConfigMode:  schema.SchemaConfigModeBlock,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Description: "Resource Name, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/names/",
					Type:        schema.TypeString,
					Required:    true,
					ForceNew:    true,
				},
				"namespace": {
					Description: "Resource Namespace use only Namespaced Resource, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",
					Type:        schema.TypeString,
					Optional:    true,
					ForceNew:    true,
				},
				"labels": {
					Description: "Resource Annotations, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/",
					Type:        schema.TypeMap,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"annotations": {
					Description: "Resource Annotations, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/",
					Type:        schema.TypeMap,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

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

func MarshalApiObjectMeta(set *schema.Set) (*idl_common.ObjectMeta, error) {
	if set == nil {
		return nil, fmt.Errorf("namespace Block Should not be null")
	}

	rawItem := set.List()[0]
	buf, err := json.Marshal(rawItem)
	if err != nil {
		return nil, fmt.Errorf("unmarshal Terraform Object to Json Failed: %v", err.Error())
	}

	meta := &idl_common.ObjectMeta{}
	err = json.Unmarshal(buf, meta)
	if err != nil {
		return nil, fmt.Errorf("marshal Json to Go Struct Failed: %v", err.Error())
	}

	return meta, nil
}
