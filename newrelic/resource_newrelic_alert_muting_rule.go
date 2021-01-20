package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/newrelic/newrelic-client-go/pkg/errors"
	"log"
)

func scheduleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"end_repeat": {
				Type: schema.TypeString,
				Optional: true,
				Description: "The datetime stamp when the MutingRule schedule should stop repeating.",
				//TODO: add validation func
			},
			"end_time": {
				Type: schema.TypeString,
				Optional: true,
				Description: "The datetime stamp representing when the MutingRule should end.",
				//TODO: add validation func
			},
			"repeat": {
				//TODO: should this be an enum type? Should we mention enum values in desc?
				Type: schema.TypeString,
				Optional: true,
				Description: "The frequency the MutingRule schedule repeats.",
				//TODO: add validation func
			},
			"repeat_count": {
				Type: schema.TypeInt,
				Optional: true,
				Description: "The number of times the MutingRule schedule should repeat.",
				//TODO: add validation func?
			},
			"start_time": {
				Type: schema.TypeString,
				Optional: true,
				Description: "The datetime stamp representing when the MutingRule should start.",
				//TODO: add validation func
			},
			"time_zone": {
				Type: schema.TypeString,
				Required: true,
				Description: "The time zone that applies to the MutingRule schedule.",
				//TODO: add validation func
			},
			"weekly_repeat_days": {
				Type: schema.TypeString,
				Optional: true,
				Description: "The day(s) of the week that a MutingRule should repeat when the repeat field is set to WEEKLY.",
				//TODO: Change to Type: schema.TypeList,
				//MinItems: 1,
				//MaxItems: 7,
				//Elem:
			},
		},
	}
}

func resourceNewRelicAlertMutingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNewRelicAlertMutingRuleCreate,
		Read:   resourceNewRelicAlertMutingRuleRead,
		Update: resourceNewRelicAlertMutingRuleUpdate,
		Delete: resourceNewRelicAlertMutingRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The account id of the MutingRule..",
			},
			"condition": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The condition that defines which violations to target.",
				MaxItems:    1,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The individual MutingRuleConditions within the group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attribute": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice([]string{"accountId", "conditionId", "policyId", "policyName", "conditionName", "conditionType", "conditionRunbookUrl", "product", "targetId", "targetName", "nrqlEventType", "tag", "nrqlQuery"}, false),
										Description:  "The attribute on a violation.",
									},
									"operator": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The operator used to compare the attribute's value with the supplied value(s).",
									},
									"values": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The value(s) to compare against the attribute's value.",
										MinItems:    1,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"operator": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The operator used to combine all the MutingRuleConditions within the group.",
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Whether the MutingRule is enabled.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the MutingRule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the MutingRule.",
			},
			"schedule": {
				Type:          schema.TypeList,
				MinItems:      1,
				MaxItems:      1,
				Optional:      true,
				Elem:          scheduleSchema(),
				Description:   "The time window when the MutingRule should actively mute violations.",
			},
		},
	}
}

func resourceNewRelicAlertMutingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(*ProviderConfig)
	client := providerConfig.NewClient

	createInput := expandMutingRuleCreateInput(d)

	accountID := selectAccountID(providerConfig, d)

	log.Printf("[INFO] Creating New Relic MutingRule alerts")

	created, err := client.Alerts.CreateMutingRule(accountID, createInput)
	if err != nil {
		return err
	}

	d.SetId(serializeIDs([]int{accountID, created.ID}))

	return resourceNewRelicAlertMutingRuleRead(d, meta)
}

func resourceNewRelicAlertMutingRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Reading New Relic MutingRule alerts")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	mutingRule, err := client.Alerts.GetMutingRule(accountID, mutingRuleID)
	if err != nil {
		if _, ok := err.(*errors.NotFound); ok {
			d.SetId("")
			return nil
		}

		return err
	}

	return flattenMutingRule(mutingRule, d)
}

func resourceNewRelicAlertMutingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient
	updateInput := expandMutingRuleUpdateInput(d)

	log.Printf("[INFO] Updating New Relic One alert muting rule.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	_, err = client.Alerts.UpdateMutingRule(accountID, mutingRuleID, updateInput)
	if err != nil {
		d.SetId("")
		return nil
	}

	return resourceNewRelicAlertMutingRuleRead(d, meta)
}

func resourceNewRelicAlertMutingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).NewClient

	log.Printf("[INFO] Deleting New Relic One muting rule alert.")

	ids, err := parseHashedIDs(d.Id())
	if err != nil {
		return err
	}

	accountID := ids[0]
	mutingRuleID := ids[1]

	err = client.Alerts.DeleteMutingRule(accountID, mutingRuleID)
	if err != nil {
		return err
	}

	return nil
}
