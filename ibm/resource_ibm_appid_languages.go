package ibm

import (
	"context"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func resourceIBMAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		Description:   "User localization configuration",
		CreateContext: resourceIBMAppIDLanguagesCreate,
		ReadContext:   resourceIBMAppIDLanguagesRead,
		DeleteContext: resourceIBMAppIDLanguagesDelete,
		UpdateContext: resourceIBMAppIDLanguagesCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"languages": {
				Description: "The list of languages that can be used to customize email templates for Cloud Directory",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDLanguagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	langs, resp, err := appIDClient.GetLocalizationWithContext(ctx, &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing language configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID languages: %s", err)
	}

	d.Set("languages", langs.Languages)
	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDLanguagesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	languages := expandStringList(d.Get("languages").([]interface{}))

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: languages,
	}

	_, err = appIDClient.UpdateLocalizationWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID languages: %s", err)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDLanguagesRead(ctx, d, meta)
}

func resourceIBMAppIDLanguagesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.UpdateLocalizationOptions{
		TenantID:  &tenantID,
		Languages: []string{"en"}, // AppID default
	}

	_, err = appIDClient.UpdateLocalizationWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID languages: %s", err)
	}

	d.SetId("")

	return nil
}
