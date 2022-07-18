package api_object

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func ApiObjectMetaFields(isNamespaced bool) map[string]*schema.Schema {
	schemes := map[string]*schema.Schema{
		"name": {
			Description: "Resource Name, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/names/",
			Type:        schema.TypeString,
			Required:    true,
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

	if isNamespaced {
		schemes["namespace"] = &schema.Schema{
			Description: "Resource Namespace use only Namespaced Resource, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/",
			Type:        schema.TypeString,
			Required:    isNamespaced,
			ForceNew:    true,
		}
	}
	return schemes
}

func ApiObjectMeta(isNamespaced bool) *schema.Schema {
	return &schema.Schema{
		Description: "Api-Object Metadata Block",
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: ApiObjectMetaFields(isNamespaced),
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

func ApiVersion(computed bool) *schema.Schema {
	var validFunc schema.SchemaValidateDiagFunc
	if !computed {
		validFunc = validation.ToDiagFunc(validation.StringIsNotEmpty)
	}

	return &schema.Schema{
		Description:      "RedFox ApiVersion, Same as ...",
		Type:             schema.TypeString,
		Required:         !computed,
		Computed:         computed,
		ValidateDiagFunc: validFunc,
	}
}

func Kind(computed bool) *schema.Schema {
	var validFunc schema.SchemaValidateDiagFunc
	if !computed {
		validFunc = validation.ToDiagFunc(validation.StringIsNotEmpty)
	}

	return &schema.Schema{
		Description:      "RedFox Kind, Same as ...",
		Type:             schema.TypeString,
		Required:         !computed,
		Computed:         computed,
		ValidateDiagFunc: validFunc,
	}
}

type ResourceId struct {
	Gvk       *idl_common.GroupVersionKindSpec
	Namespace string
	Name      string
}

func NewResourceIdFull(gvk *idl_common.GroupVersionKindSpec, namespace string, name string) *ResourceId {
	return &ResourceId{Gvk: gvk, Namespace: namespace, Name: name}
}

func NewResourceId(gvk *idl_common.GroupVersionKindSpec, name string) *ResourceId {
	return &ResourceId{Gvk: gvk, Namespace: "@fox-system", Name: name}
}

func (i *ResourceId) String() string {
	return fmt.Sprintf("%s:%s:%s/%s/%s", i.Gvk.Group, i.Gvk.Version, i.Gvk.Kind, i.Namespace, i.Name)
}

func (i *ResourceId) ApiVersion() string {
	return fmt.Sprintf("%s/%s", i.Gvk.Group, i.Gvk.Version)
}

func ParseResourceId(id string) (*ResourceId, error) {
	args := strings.Split(id, "/")
	if len(args) != 3 {
		return nil, fmt.Errorf("ParseResourceIdFailed: ResourceId should contains 3 slash '/' but it has %d slashs, id: '%s'", len(args), id)
	}

	rawGvks := strings.Split(args[0], ":")
	if len(rawGvks) != 3 {
		return nil, fmt.Errorf("ParseResourceIdFailed: Gvk should contains 3 colon ':' but it has %d colon, id: '%s'", len(rawGvks), args[0])
	}

	return &ResourceId{
		Gvk: &idl_common.GroupVersionKindSpec{
			Group:   rawGvks[0],
			Version: rawGvks[1],
			Kind:    rawGvks[2],
		},
		Namespace: args[1],
		Name:      args[2],
	}, nil
}
