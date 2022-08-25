package redfox

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redfoxV1alpha1 "github.com/krafton-hq/redfox/pkg/apis/redfox/v1alpha1"
	"github.com/samber/lo"
)

func expandClusterSpec(cluster []any) (*redfoxV1alpha1.ClusterSpec, error) {
	obj := &redfoxV1alpha1.ClusterSpec{}

	if len(cluster) == 0 || cluster[0] == nil {
		return obj, nil
	}

	in := cluster[0].(map[string]any)

	obj.ClusterName = in["cluster_name"].(string)
	obj.ClusterGroup = in["cluster_group"].(string)
	obj.ClusterEngine = in["cluster_engine"].(string)
	obj.ClusterRegion = in["cluster_region"].(string)
	obj.InfraAccountId = in["infra_account_id"].(string)
	obj.InfraVendor = in["infra_vendor"].(string)
	obj.ServicePhase = in["service_phase"].(string)
	obj.ServiceTag = in["service_tag"].(string)

	rawRoles := in["roles"].([]any)
	for _, role := range rawRoles {
		obj.Roles = append(obj.Roles, redfoxV1alpha1.ClusterRole(role.(string)))
	}
	return obj, nil
}

func flattenClusterSpec(in redfoxV1alpha1.ClusterSpec, d *schema.ResourceData, meta interface{}) ([]any, error) {
	att := map[string]any{}
	att["cluster_name"] = in.ClusterName
	att["cluster_group"] = in.ClusterGroup
	att["cluster_engine"] = in.ClusterEngine
	att["cluster_region"] = in.ClusterRegion
	att["infra_account_id"] = in.InfraAccountId
	att["infra_vendor"] = in.InfraVendor
	att["service_phase"] = in.ServicePhase
	att["service_tag"] = in.ServiceTag
	att["roles"] = in.Roles
	return []any{att}, nil
}

func expandClusterStatus(cluster []any) (*redfoxV1alpha1.ClusterStatus, error) {
	obj := &redfoxV1alpha1.ClusterStatus{}

	if len(cluster) == 0 || cluster[0] == nil {
		return obj, nil
	}

	in := cluster[0].(map[string]any)

	// service_account_issuer is required field
	obj.ServiceAccountIssuer = in["service_account_issuer"].(string)

	// apiserver is required field
	rawApiserver := in["apiserver"].([]any)
	if len(rawApiserver) == 0 || rawApiserver[0] == nil {
		return nil, fmt.Errorf("`apiserver` is required block")
	}
	inApiserver := rawApiserver[0].(map[string]any)
	obj.Apiserver.Endpoint = inApiserver["endpoint"].(string)
	obj.Apiserver.CaCert = inApiserver["ca_cert"].(string)

	// aws_iam_external_idps is optional field
	if rawAwsIdps, found := in["aws_iam_idps"]; found {
		obj.AwsIamIdps = lo.MapValues[string, any, string](rawAwsIdps.(map[string]any), func(v any, _ string) string {
			return v.(string)
		})
	}

	return obj, nil
}

func flattenClusterStatus(in redfoxV1alpha1.ClusterStatus, d *schema.ResourceData, meta interface{}) ([]any, error) {
	att := map[string]any{}
	att["service_account_issuer"] = in.ServiceAccountIssuer
	apiserver := map[string]string{}
	apiserver["endpoint"] = in.Apiserver.Endpoint
	apiserver["ca_cert"] = in.Apiserver.CaCert
	att["apiserver"] = []any{apiserver}
	att["aws_iam_idps"] = in.AwsIamIdps
	return []any{att}, nil
}
