package ibm

import (
	"context"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIBMAppIDIDPSAMLMetadata() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieve SAML metadata",
		ReadContext: dataSourceIBMAppIDIDPSAMLMetadataRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"metadata": {
				Type:        schema.TypeString,
				Description: "SAML Metadata",
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPSAMLMetadataRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appidClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	metadata, _, err := appidClient.GetSAMLMetadataWithContext(ctx, &appid.GetSAMLMetadataOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID SAML metadata: %s", err)
	}

	if err := d.Set("metadata", metadata); err != nil {
		return diag.Errorf("Error setting AppID SAML metadata: %s", err)
	}

	d.SetId(tenantID)

	return nil
}
