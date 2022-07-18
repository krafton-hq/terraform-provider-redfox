package api_object

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
)

func EndpointSpecFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"addresses": {
			Description: "An URL",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    0,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"url": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
				},
			},
		},
		"ports": {
			Description: "An Endpoint Port",
			Type:        schema.TypeList,
			Required:    true,
			MinItems:    0,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
					"port": {
						Type:             schema.TypeInt,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.IsPortNumber),
					},
					"protocol": {
						Type:             schema.TypeString,
						Required:         true,
						ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
					},
				},
			},
		},
	}
}

func MarshalEndpointSpec(raw []any) (*documents.EndpointSpec, error) {
	if raw == nil {
		return nil, fmt.Errorf("crd Block Should not be null")
	}

	spec := &documents.EndpointSpec{}
	rawItem := raw[0].(map[string]interface{})

	for _, rawAddress := range rawItem["addresses"].([]any) {
		buf, err := json.Marshal(rawAddress)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Terraform Object to Json Failed: %v", err.Error())
		}

		address := &documents.EndpointAddress{}
		err = json.Unmarshal(buf, address)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Struct Failed: %v", err.Error())
		}
		spec.Addresses = append(spec.Addresses, address)
	}

	for _, rawPorts := range rawItem["ports"].([]any) {
		buf, err := json.Marshal(rawPorts)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Terraform Object to Json Failed: %v", err.Error())
		}

		port := &documents.EndpointPort{}
		err = json.Unmarshal(buf, port)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Struct Failed: %v", err.Error())
		}
		spec.Ports = append(spec.Ports, port)
	}

	return spec, nil
}

func UnmarshalEndpointSpec(spec *documents.EndpointSpec) ([]any, error) {
	if spec == nil {
		return nil, fmt.Errorf("endpoint spec Block Should not be null")
	}

	var rawAddresses []any
	for _, address := range spec.Addresses {
		buf, err := json.Marshal(address)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Go Struct to Json Failed: %v", err.Error())
		}

		raw := map[string]any{}
		err = json.Unmarshal(buf, &raw)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Map Failed: %v", err.Error())
		}
		rawAddresses = append(rawAddresses, raw)
	}

	var rawPorts []any
	for _, port := range spec.Ports {
		buf, err := json.Marshal(port)
		if err != nil {
			return nil, fmt.Errorf("unmarshal Go Struct to Json Failed: %v", err.Error())
		}

		raw := map[string]any{}
		err = json.Unmarshal(buf, &raw)
		if err != nil {
			return nil, fmt.Errorf("marshal Json to Go Map Failed: %v", err.Error())
		}
		rawPorts = append(rawPorts, raw)
	}

	rawItem := map[string]any{
		"addresses": rawAddresses,
		"ports":     rawPorts,
	}

	return []any{rawItem}, nil
}

func AssembleEndpoint(gvk *idl_common.GroupVersionKindSpec, metadata *idl_common.ObjectMeta, spec *documents.EndpointSpec) *documents.Endpoint {
	return &documents.Endpoint{
		ApiVersion: gvk.Group + "/" + gvk.Version,
		Kind:       gvk.Kind,
		Metadata:   metadata,
		Spec:       spec,
	}
}

func DisassembleEndpoint(endpoint *documents.Endpoint) (*idl_common.GroupVersionKindSpec, *idl_common.ObjectMeta, *documents.EndpointSpec) {
	gvk, err := ParseGvk(endpoint.ApiVersion, endpoint.Kind)
	// Unreachable
	if err != nil {
		tflog.Error(context.TODO(), err.Error())
	}

	return gvk, endpoint.Metadata, endpoint.Spec
}
