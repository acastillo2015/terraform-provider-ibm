package ibm

import (
	"context"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description:   "Redirect URIs that can be used as callbacks of App ID authentication flow",
		CreateContext: resourceIBMAppIDRedirectURLsCreate,
		ReadContext:   resourceIBMAppIDRedirectURLsRead,
		UpdateContext: resourceIBMAppIDRedirectURLsUpdate,
		DeleteContext: resourceIBMAppIDRedirectURLsDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDRedirectURLsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	urls, _, err := appIDClient.GetRedirectUrisWithContext(ctx, &appid.GetRedirectUrisOptions{
		TenantID: &tenantID,
	})
	if err != nil {
		return diag.Errorf("Error loading AppID Cloud Directory redirect urls: %s", err)
	}

	if err := d.Set("urls", urls.RedirectUris); err != nil {
		return diag.Errorf("Error setting AppID Cloud Directory redirect urls: %s", err)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDRedirectURLsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	urls := d.Get("urls")

	redirectURLs := expandStringList(urls.([]interface{}))
	_, err = appIDClient.UpdateRedirectUrisWithContext(ctx, &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: redirectURLs,
		},
	})

	if err != nil {
		return diag.Errorf("Error updating AppID Cloud Directory redirect URLs: %s", err)
	}

	d.SetId(tenantID)
	return resourceIBMAppIDRedirectURLsRead(ctx, d, meta)
}

func resourceIBMAppIDRedirectURLsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	urls := d.Get("urls")

	redirectURLs := expandStringList(urls.([]interface{}))
	_, err = appIDClient.UpdateRedirectUrisWithContext(ctx, &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: redirectURLs,
		},
	})

	if err != nil {
		return diag.Errorf("Error updating AppID Cloud Directory redirect URLs: %s", err)
	}

	return resourceIBMAppIDRedirectURLsRead(ctx, d, meta)
}

func resourceIBMAppIDRedirectURLsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	_, err = appIDClient.UpdateRedirectUrisWithContext(ctx, &appid.UpdateRedirectUrisOptions{
		TenantID: &tenantID,
		RedirectUrisArray: &appid.RedirectURIConfig{
			RedirectUris: []string{},
		},
	})

	if err != nil {
		return diag.Errorf("Error resetting AppID Cloud Directory redirect URLs: %s", err)
	}

	d.SetId("")

	return nil
}
