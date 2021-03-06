---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "redfox_endpoint Data Source - terraform-provider-redfox"
subcategory: ""
description: |-
  RedFox Endpoint
---

# redfox_endpoint (Data Source)

RedFox Endpoint
```terraform
data "redfox_endpoint" "endpoint" {
  metadata {
    name      = "portal"
    namespace = "infra"
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `metadata` (Block List, Min: 1, Max: 1) Api-Object Metadata Block (see [below for nested schema](#nestedblock--metadata))

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `api_version` (String) RedFox ApiVersion, Same as ...
- `id` (String) The ID of this resource.
- `kind` (String) RedFox Kind, Same as ...
- `spec` (List of Object) Endpoint Spec Block (see [below for nested schema](#nestedatt--spec))

<a id="nestedblock--metadata"></a>
### Nested Schema for `metadata`

Required:

- `name` (String) Resource Name, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/names/
- `namespace` (String) Resource Namespace use only Namespaced Resource, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/

Optional:

- `annotations` (Map of String) Resource Annotations, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
- `labels` (Map of String) Resource Labels, Same as https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `default` (String)


<a id="nestedatt--spec"></a>
### Nested Schema for `spec`

Read-Only:

- `addresses` (List of Object) (see [below for nested schema](#nestedobjatt--spec--addresses))
- `ports` (List of Object) (see [below for nested schema](#nestedobjatt--spec--ports))

<a id="nestedobjatt--spec--addresses"></a>
### Nested Schema for `spec.addresses`

Read-Only:

- `url` (String)


<a id="nestedobjatt--spec--ports"></a>
### Nested Schema for `spec.ports`

Read-Only:

- `name` (String)
- `port` (Number)
- `protocol` (String)


