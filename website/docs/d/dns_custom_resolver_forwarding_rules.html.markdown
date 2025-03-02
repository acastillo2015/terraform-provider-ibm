---
subcategory: "DNS Services"
layout: "ibm"
page_title: "IBM : Forwarding Rules"
description: |-
  Manages IBM Cloud Infrastructure private domain name service forwarding rules.
---

# ibm_dns_custom_resolver_forwarding_rules

Provides a read-only data source for forwarding rules. You can then reference the fields of the data source in other resources within the same configuration using interpolation syntax. For more information about forwarding rules, refer to [list-forwarding-rules](https://cloud.ibm.com/apidocs/dns-svcs#list-forwarding-rules)

## Example Usage

```terraform
data "ibm_dns_custom_resolver_forwarding_rules" "dns_custom_resolver_forwarding_rules" {
  instance_id = ibm_dns_custom_resolver_forwarding_rule.dns_custom_resolver_forwarding_rule.instance_id
  resolver_id = ibm_dns_custom_resolver_forwarding_rule.dns_custom_resolver_forwarding_rule.resolver_id
}
```

## Argument Reference

Review the argument reference that you can specify for your data source.

- `instance_id` - (Required, String) The GUID of the private DNS service instance.
- `resolver_id` - (Required, String) The unique identifier of a custom resolver.

## Attribute Reference

In addition to the argument references list, you can access the following attribute references after your data source are created.

- `forwarding_rules` (List) List of forwarding rules.

	Nested scheme for `forwarding_rules`:
	- `description` - (String) Descriptive text of the forwarding rule.
	- `forward_to` - (String) The upstream DNS servers will be forwarded to.
	- `match` - (String) The matching zone or hostname.
	- `rule_id` - (String) Identifier of the forwarding rule.
	- `type` - (String) Type of the forwarding rule.

	
	
	

