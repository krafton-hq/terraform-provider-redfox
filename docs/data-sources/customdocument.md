---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "redfox_customdocument Data Source - terraform-provider-redfox"
subcategory: ""
description: |-
  RedFox CustomDocument
---

# redfox_customdocument (Data Source)

RedFox CustomDocument
```terraform
data "redfox_customdocument" "customdocument" {
  metadata {
    name      = "seoul-cluster"
    namespace = "infra"
  }
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_version` (String) RedFox ApiVersion, Same as ...
- `kind` (String) RedFox Kind, Same as ...
- `metadata` (Block List, Min: 1, Max: 1) Api-Object Metadata Block (see [below for nested schema](#nestedblock--metadata))

### Optional

- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `id` (String) The ID of this resource.
- `spec` (List of Object) CustomDocument Spec Block (see [below for nested schema](#nestedatt--spec))

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

- `raw_json` (String)


