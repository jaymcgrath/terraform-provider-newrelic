// +build integration

package newrelic

import (
	"testing"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/alerts"
	"github.com/stretchr/testify/require"
)

// This test currently fails due to the time zone being
// appended via formatting, but should be resolved with
// the custom unmarshal method being implemented.
func TestFlattenSchedule(t *testing.T) {
	t.Parallel()

	timestamp := time.Now()
	repeat := alerts.MutingRuleScheduleRepeat("WEEKLY")

	mockMutingRuleSchedule := alerts.MutingRuleSchedule{
		StartTime: &timestamp,
		EndTime:   &timestamp,
		TimeZone:  "America/Los_Angeles",
		Repeat:    &repeat,
		EndRepeat: &timestamp,
		WeeklyRepeatDays: &[]alerts.DayOfWeek{
			"MONDAY",
			"TUESDAY",
		},
	}

	mockScheduleConfig := map[string]interface{}{
		"start_time": "2021-01-21T15:30:00",
		"end_time":   "2021-01-21T15:30:00",
		"end_repeat": "2021-01-21T15:30:00",
		"time_zone":  "America/Los_Angeles",
		"repeat":     "WEEKLY",
		"weekly_repeat_days": []string{
			"MONDAY",
			"TUESDAY",
		},
	}

	result := flattenSchedule(&mockMutingRuleSchedule, mockScheduleConfig)

	require.Equal(t, []interface{}{mockScheduleConfig}, result)
}
