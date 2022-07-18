package redfox_helper

import (
	"context"
	"fmt"

	"github.com/krafton-hq/red-fox/apis/crds"
	"github.com/krafton-hq/red-fox/apis/documents"
	"github.com/krafton-hq/red-fox/apis/idl_common"
	"github.com/krafton-hq/red-fox/apis/namespaces"
	client_sdk "github.com/krafton-hq/red-fox/client-sdk"
)

type ClientHelper interface {
	Namespaces() namespaces.NamespaceServerClient
	NatIps() documents.NatIpServerClient
	Endpoints() documents.EndpointServerClient
	CustomDocuments() documents.CustomDocumentServerClient
	Crds() crds.CustomResourceDefinitionServerClient

	RawClient() *client_sdk.RedFoxClient
	NamespaceGvk() *idl_common.GroupVersionKindSpec
	NatIpGvk() *idl_common.GroupVersionKindSpec
	EndpointGvk() *idl_common.GroupVersionKindSpec
	CrdGvk() *idl_common.GroupVersionKindSpec
	ApiResources() []*idl_common.ApiResourceSpec
}

type clientStruct struct {
	redfoxClient *client_sdk.RedFoxClient
	gvks         map[string]*idl_common.GroupVersionKindSpec
	apiResources []*idl_common.ApiResourceSpec
}

func NewClient(ctx context.Context, redfoxClient *client_sdk.RedFoxClient) (*clientStruct, error) {
	res, err := redfoxClient.ApiResourcesServerClient.ListApiResources(ctx, &idl_common.CommonReq{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Supported ApiResources, error: %v", err.Error())
	}

	client := &clientStruct{
		redfoxClient: redfoxClient,
		apiResources: res.ApiResources,
		gvks:         map[string]*idl_common.GroupVersionKindSpec{},
	}

	for _, resource := range res.ApiResources {
		client.gvks[resource.Name] = resource.Gvk
	}

	return client, nil
}

func (c *clientStruct) Namespaces() namespaces.NamespaceServerClient {
	return c.redfoxClient.NamespaceServerClient
}

func (c *clientStruct) NatIps() documents.NatIpServerClient {
	return c.redfoxClient.NatIpServerClient
}

func (c *clientStruct) Endpoints() documents.EndpointServerClient {
	return c.redfoxClient.EndpointServerClient
}

func (c *clientStruct) Crds() crds.CustomResourceDefinitionServerClient {
	return c.redfoxClient.CustomResourceDefinitionServerClient
}

func (c *clientStruct) RawClient() *client_sdk.RedFoxClient {
	return c.redfoxClient
}

func (c *clientStruct) NamespaceGvk() *idl_common.GroupVersionKindSpec {
	return c.gvks["namespace.core"]
}

func (c *clientStruct) NatIpGvk() *idl_common.GroupVersionKindSpec {
	return c.gvks["natip.red-fox.sbx-central.io"]
}

func (c *clientStruct) EndpointGvk() *idl_common.GroupVersionKindSpec {
	return c.gvks["endpoint.red-fox.sbx-central.io"]
}

func (c *clientStruct) CrdGvk() *idl_common.GroupVersionKindSpec {
	return c.gvks["customresourcedefinition.red-fox.sbx-central.io"]
}

func (c *clientStruct) ApiResources() []*idl_common.ApiResourceSpec {
	return c.apiResources
}
