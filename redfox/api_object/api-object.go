package api_object

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func ApiObjectMetaFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
			Description: "Resource Labels, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/",
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
	}
}

func ApiObjectMeta() *schema.Schema {
	return &schema.Schema{
		Description: "Api-Object Metadata Block",
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: ApiObjectMetaFields(),
		},
	}
}

func MarshalApiObjectMeta(raw []any) (*idl_common.ObjectMeta, error) {
	if raw == nil {
		return nil, fmt.Errorf("metadata Block Should not be null")
	}

	buf, err := json.Marshal(raw[0])
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

func UnmarshalApiObjectMeta(metadata *idl_common.ObjectMeta) ([]any, error) {
	if metadata == nil {
		return nil, fmt.Errorf("metadata Block Should not be null")
	}

	buf, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("unmarshal Go Struct to Json Failed: %v", err.Error())
	}

	raw := map[string]any{}
	err = json.Unmarshal(buf, &raw)
	if err != nil {
		return nil, fmt.Errorf("marshal Json to Go Map Failed: %v", err.Error())
	}

	return []any{raw}, nil
}

func LabelSelector() *schema.Schema {
	return &schema.Schema{
		Description: "Resource Label Selectors, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/",
		Type:        schema.TypeMap,
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}
}

func MarshalLabelSelectors(raw map[string]any) map[string]string {
	if raw == nil {
		return nil
	}

	selector := map[string]string{}
	for key, value := range raw {
		selector[key] = value.(string)
	}

	return selector
}

func BuildNamespaceObjectId(namespace, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func ParseNamespacedObjectId(id string) (namespace string, name string, found bool) {
	return strings.Cut(id, "/")
}

func BuildClusterObjectId(name string) string {
	return name
}

func ParseClusterObjectId(id string) (name string) {
	return id
}
