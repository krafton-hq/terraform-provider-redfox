package api_object

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/samber/lo"
)

func NatIpSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ip_type": {
			Description:  "IP Type, Can be either IPv4 or IPv6",
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice([]string{documents.IpType_Ipv4.String(), documents.IpType_Ipv6.String()}, false),
		},
		"cidrs": {
			Description: "Classless Inter-Domain Routing notated IP List, Must be end '/<bits>'",
			Type:        schema.TypeList,
			Required:    true,
			Elem: &schema.Schema{
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
		},
	}
}

func NatIpResourceSpec() *schema.Schema {
	return &schema.Schema{
		Description: "NatIp Spec Block",
		Type:        schema.TypeList,
		Required:    true,
		MinItems:    1,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: NatIpSpecFields(),
		},
	}
}

func MarshalNatIpSpec(raw []any) (*documents.NatIpSpec, error) {
	if raw == nil {
		return nil, fmt.Errorf("natip Block Should not be null")
	}

	spec := &documents.NatIpSpec{}
	rawItem := raw[0].(map[string]interface{})

	if rawIpType, found := rawItem["ip_type"]; found {
		if value, found := documents.IpType_value[rawIpType.(string)]; found {
			spec.Type = *documents.IpType(value).Enum()
		} else {
			return nil, fmt.Errorf("unexpected ip_type: '%v'", rawIpType)
		}
	}

	cidrs := lo.Map[any, string](rawItem["cidrs"].([]any), func(x any, _ int) string {
		return x.(string)
	})
	spec.Cidrs = cidrs

	return spec, nil
}

func UnmarshalNatIpSpec(spec *documents.NatIpSpec) ([]any, error) {
	if spec == nil {
		return nil, fmt.Errorf("natIp Block Should not be null")
	}

	rawCidrs := lo.Map[string, any](spec.Cidrs, func(x string, _ int) any {
		return x
	})

	rawItem := map[string]any{
		"ip_type": spec.Type.String(),
		"cidrs":   rawCidrs,
	}
	return []any{rawItem}, nil
}
