package ibm

import (
	"context"
	"fmt"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
)

func resourceIBMAppIDActionURL() *schema.Resource {
	return &schema.Resource{
		Description:   "The custom url to redirect to when Cloud Directory action is executed.",
		CreateContext: resourceIBMAppIDActionURLCreate,
		ReadContext:   resourceIBMAppIDActionURLRead,
		DeleteContext: resourceIBMAppIDActionURLDelete,
		UpdateContext: resourceIBMAppIDActionURLUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"action": {
				Description:  "The type of the action: `on_user_verified` - the URL of your custom user verified page, `on_reset_password` - the URL of your custom reset password page",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on_user_verified", "on_reset_password"}, false),
				Required:     true,
				ForceNew:     true,
			},
			"url": {
				Description: "The action URL",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceIBMAppIDActionURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: AppID action URL ID should be a combination of tenantID/action", id)
	}

	tenantID := idParts[0]
	action := idParts[1]

	cfg, resp, err := appIDClient.GetCloudDirectoryActionURLWithContext(ctx, &appid.GetCloudDirectoryActionURLOptions{
		TenantID: &tenantID,
		Action:   &action,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Action URL from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID actionURL: %s", err)
	}

	if cfg.ActionURL != nil {
		d.Set("url", *cfg.ActionURL)
	}

	d.Set("tenant_id", tenantID)
	d.Set("action", action)

	return nil
}

func resourceIBMAppIDActionURLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)
	actionURL := d.Get("url").(string)

	input := &appid.SetCloudDirectoryActionOptions{
		TenantID:  &tenantID,
		Action:    &action,
		ActionURL: &actionURL,
	}

	_, _, err = appIDClient.SetCloudDirectoryActionWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting AppID Cloud Directory action URL: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, action))

	return resourceIBMAppIDActionURLRead(ctx, d, meta)
}

func resourceIBMAppIDActionURLDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)

	_, err = appIDClient.DeleteActionURLWithContext(ctx, &appid.DeleteActionURLOptions{
		TenantID: &tenantID,
		Action:   &action,
	})

	if err != nil {
		return diag.Errorf("Error deleting AppID Cloud Directory action URL: %s", err)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDActionURLUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIBMAppIDActionURLCreate(ctx, d, m)
}
