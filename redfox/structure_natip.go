package redfox

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	"github.com/samber/lo"
)

func expandNatIpSpec(natip []any) (*redfoxV1alpha1.NatIpSpec, error) {
	obj := &redfoxV1alpha1.NatIpSpec{}

	if len(natip) == 0 || natip[0] == nil {
		return obj, nil
	}

	in := natip[0].(map[string]any)

	// Optional field
	if rawIpType, found := in["ip_type"].(string); found {
		obj.IpType = redfoxV1alpha1.IpType(rawIpType)
	}

	// Required field
	cidrs := lo.Map[any, string](in["cidrs"].([]any), func(x any, _ int) string {
		return x.(string)
	})
	obj.Cidrs = cidrs
	return obj, nil
}

func flattenNatIpSpec(in redfoxV1alpha1.NatIpSpec, d *schema.ResourceData, meta interface{}) ([]any, error) {
	att := map[string]any{}
	att["ip_type"] = string(in.IpType)
	att["cidrs"] = in.Cidrs
	return []any{att}, nil
}
