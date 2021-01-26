package newrelic

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"strings"
	"time"
)

func expandMutingRuleCreateInput(d *schema.ResourceData) (alerts.MutingRuleCreateInput, error) {
	createInput := alerts.MutingRuleCreateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok {
		createInput.Condition = expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))
	}

	if e, ok := d.GetOk("schedule"); ok {
		schedule, err := expandMutingRuleCreateSchedule(e.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return alerts.MutingRuleCreateInput{}, err
		}
		createInput.Schedule = &schedule
	}

	return createInput, nil
}

func expandMutingRuleCreateSchedule(cfg map[string]interface{}) (alerts.MutingRuleScheduleCreateInput, error) {
	schedule := alerts.MutingRuleScheduleCreateInput{}

	if startTime, ok := cfg["start_time"]; ok {
		rawStartTime := startTime.(string)
		if rawStartTime != "" {
			formattedStartTime, err := time.Parse("2006-01-02T15:04:05", rawStartTime)
			if err != nil {
				return alerts.MutingRuleScheduleCreateInput{}, err
			}
			schedule.StartTime = &alerts.NaiveDateTime{Time: formattedStartTime}
		}
	}

	if endTime, ok := cfg["end_time"]; ok {
		rawEndTime := endTime.(string)
		if rawEndTime != "" {
			formattedEndTime, err := time.Parse("2006-01-02T15:04:05", rawEndTime)
			if err != nil {
				return alerts.MutingRuleScheduleCreateInput{}, err
			}
			schedule.EndTime = &alerts.NaiveDateTime{Time: formattedEndTime}
		}
	}

	if timeZone, ok := cfg["time_zone"]; ok {
		schedule.TimeZone = timeZone.(string)
	}

	if repeat, ok := cfg["repeat"]; ok {
		r := repeat.(string)
		if r != "" {
			sr := alerts.MutingRuleScheduleRepeat(strings.ToUpper(r))
			schedule.Repeat = &sr

		}
	}

	if endRepeat, ok := cfg["end_repeat"]; ok {
		rawEndRepeat := endRepeat.(string)
		if rawEndRepeat != "" {
			formattedEndRepeat, err := time.Parse("2006-01-02T15:04:05", rawEndRepeat)
			if err != nil {
				return alerts.MutingRuleScheduleCreateInput{}, err
			}
			schedule.EndRepeat = &alerts.NaiveDateTime{Time: formattedEndRepeat}
		}

	}

	if repeatCount, ok := cfg["repeat_count"]; ok {
		r := repeatCount.(int)
		if r > 0 {
			schedule.RepeatCount = &r
		}
	}

	if weeklyRepeatDays, ok := cfg["weekly_repeat_days"]; ok {
		repeatDaysAsList := weeklyRepeatDays.(*schema.Set).List()
		if rdLen := len(repeatDaysAsList); rdLen > 0 {
			repeatDays := make([]alerts.DayOfWeek, rdLen)
			for i, day := range repeatDaysAsList {
				repeatDays[i] = alerts.DayOfWeek(strings.ToUpper(day.(string)))
			}
			schedule.WeeklyRepeatDays = &repeatDays
		}
	}
	return schedule, nil
}

func expandMutingRuleUpdateSchedule(cfg map[string]interface{}) (alerts.MutingRuleScheduleUpdateInput, error) {
	schedule := alerts.MutingRuleScheduleUpdateInput{}

	if startTime, ok := cfg["start_time"]; ok {
		rawStartTime := startTime.(string)
		if rawStartTime != "" {
			formattedStartTime, err := time.Parse("2006-01-02T15:04:05", rawStartTime)
			if err != nil {
				return alerts.MutingRuleScheduleUpdateInput{}, err
			}
			schedule.StartTime = &alerts.NaiveDateTime{Time: formattedStartTime}
		} else {
			schedule.StartTime = nil
		}
	}

	if endTime, ok := cfg["end_time"]; ok {
		rawEndTime := endTime.(string)
		if rawEndTime != "" {
			formattedEndTime, err := time.Parse("2006-01-02T15:04:05", rawEndTime)
			if err != nil {
				return alerts.MutingRuleScheduleUpdateInput{}, err
			}
			schedule.EndTime = &alerts.NaiveDateTime{Time: formattedEndTime}
		} else {
			schedule.EndTime = nil
		}
	}

	if timeZone, ok := cfg["time_zone"]; ok {
		if rawTimeZone := timeZone.(string); rawTimeZone != "" {
			schedule.TimeZone = &rawTimeZone
		}
	}

	if repeat, ok := cfg["repeat"]; ok {
		r := repeat.(string)
		if r != "" {
			sr := alerts.MutingRuleScheduleRepeat(strings.ToUpper(r))
			schedule.Repeat = &sr
		} else {
			schedule.Repeat = nil
		}
	}

	if endRepeat, ok := cfg["end_repeat"]; ok {
		rawEndRepeat := endRepeat.(string)
		if rawEndRepeat != "" {
			formattedEndRepeat, err := time.Parse("2006-01-02T15:04:05", rawEndRepeat)
			if err != nil {
				return alerts.MutingRuleScheduleUpdateInput{}, err
			}
			schedule.EndRepeat = &alerts.NaiveDateTime{Time: formattedEndRepeat}
		} else {
			schedule.EndRepeat = nil
		}
	}

	if repeatCount, ok := cfg["repeat_count"]; ok {
		r := repeatCount.(int)
		if r > 0 {
			schedule.RepeatCount = &r
		} else {
			schedule.EndRepeat = nil
		}
	}

	if weeklyRepeatDays, ok := cfg["weekly_repeat_days"]; ok {
		repeatDaysAsList := weeklyRepeatDays.(*schema.Set).List()
		if rdLen := len(repeatDaysAsList); rdLen > 0 {
			repeatDays := make([]alerts.DayOfWeek, rdLen)
			for i, day := range repeatDaysAsList {
				repeatDays[i] = alerts.DayOfWeek(strings.ToUpper(day.(string)))
			}
			schedule.WeeklyRepeatDays = &repeatDays
		} else {
			schedule.WeeklyRepeatDays = nil
		}
	}
	return schedule, nil
}

func expandMutingRuleUpdateInput(d *schema.ResourceData) (alerts.MutingRuleUpdateInput, error) {
	updateInput := alerts.MutingRuleUpdateInput{
		Enabled:     d.Get("enabled").(bool),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if e, ok := d.GetOk("condition"); ok {
		x := expandMutingRuleConditionGroup(e.([]interface{})[0].(map[string]interface{}))

		updateInput.Condition = &x
	}

	if e, ok := d.GetOk("schedule"); ok {
		schedule, err := expandMutingRuleUpdateSchedule(e.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return alerts.MutingRuleUpdateInput{}, err
		}
		updateInput.Schedule = &schedule
	}

	return updateInput, nil
}

func expandMutingRuleConditionGroup(cfg map[string]interface{}) alerts.MutingRuleConditionGroup {
	conditionGroup := alerts.MutingRuleConditionGroup{}
	var expandedConditions []alerts.MutingRuleCondition

	conditions := cfg["conditions"].([]interface{})

	for _, c := range conditions {
		var y = expandMutingRuleCondition(c)
		expandedConditions = append(expandedConditions, y)
	}

	conditionGroup.Conditions = expandedConditions

	if operator, ok := cfg["operator"]; ok {
		conditionGroup.Operator = operator.(string)
	}

	return conditionGroup
}

func expandMutingRuleCondition(cfg interface{}) alerts.MutingRuleCondition {
	conditionCfg := cfg.(map[string]interface{})
	condition := alerts.MutingRuleCondition{}

	if attribute, ok := conditionCfg["attribute"]; ok {
		condition.Attribute = attribute.(string)
	}

	if operator, ok := conditionCfg["operator"]; ok {
		condition.Operator = operator.(string)
	}
	if values, ok := conditionCfg["values"]; ok {
		condition.Values = expandMutingRuleValues(values.([]interface{}))
	}

	return condition
}

func expandMutingRuleValues(values []interface{}) []string {
	perms := make([]string, len(values))

	for i, values := range values {
		perms[i] = values.(string)
	}

	return perms
}

func flattenMutingRule(mutingRule *alerts.MutingRule, d *schema.ResourceData) error {

	x := d.Get("condition")
	configuredCondition := x.([]interface{})

	d.Set("enabled", mutingRule.Enabled)
	err := d.Set("condition", flattenMutingRuleConditionGroup(mutingRule.Condition, configuredCondition))
	if err != nil {
		return nil
	}

	d.Set("description", mutingRule.Description)
	d.Set("name", mutingRule.Name)
	d.Set("schedule", mutingRule.Schedule)
	//
	//s := schema.ResourceData{}
	//
	//if err := flattenSchedule(&s, mutingRule.Schedule); err != nil {
	//	return err
	//}
	//d.Set("schedule", s)
	return nil
}

//func flattenSchedule(d *schema.ResourceData, schedule *alerts.MutingRuleSchedule) error{
//	if schedule == nil {
//		return nil
//	}
//
//	if err := d.Set("start_time", schedule.StartTime.Format("2006-01-02T15:04:05")); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.start_time`: %v", err)
//	}
//
//	if err := d.Set("end_time", schedule.EndTime); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.end_time`: %v", err)
//	}
//
//	if err := d.Set("time_zone", schedule.TimeZone); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.time_zone`: %v", err)
//	}
//
//	if err := d.Set("repeat", schedule.Repeat); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.repeat`: %v", err)
//	}
//
//	if err := d.Set("end_repeat", schedule.EndRepeat); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.end_repeat`: %v", err)
//	}
//	if err := d.Set("repeat_count", schedule.RepeatCount); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.repeat_count`: %v", err)
//	}
//	if err := d.Set("weekly_repeat_days", schedule.WeeklyRepeatDays); err != nil {
//		return fmt.Errorf("[DEBUG] Error setting alert muting rule `schedule.weekly_repeat_days`: %v", err)
//	}
//
//	return nil
//}

func flattenMutingRuleConditionGroup(in alerts.MutingRuleConditionGroup, configuredCondition []interface{}) []map[string]interface{} {

	condition := []map[string]interface{}{
		{
			"operator": in.Operator,
		},
	}

	if len(in.Conditions) > 0 {

		condition[0]["conditions"] = handleImportFlattenCondition(in.Conditions)

	} else {

		condition[0]["conditions"] = flattenMutingRuleCondition(configuredCondition)
	}

	return condition
}
func handleImportFlattenCondition(conditions []alerts.MutingRuleCondition) []map[string]interface{} {
	var condition []map[string]interface{}

	for _, src := range conditions {
		dst := map[string]interface{}{
			"attribute": src.Attribute,
			"operator":  src.Operator,
			"values":    src.Values,
		}
		condition = append(condition, dst)
	}

	return condition
}
func flattenMutingRuleCondition(conditions []interface{}) []map[string]interface{} {
	var condition []map[string]interface{}

	for _, src := range conditions {
		x := src.(map[string]interface{})

		if x["values"] != nil && x["attributes"] != "" && x["operator"] != "" {
			dst := map[string]interface{}{
				"attribute": x["attribute"],
				"operator":  x["operator"],
				"values":    x["values"],
			}
			condition = append(condition, dst)
		}

	}

	return condition
}
