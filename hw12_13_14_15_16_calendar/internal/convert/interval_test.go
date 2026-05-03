package convert

import (
	"testing"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestIntervalFromProto(t *testing.T) {
	tests := []struct {
		name  string
		p     api.Interval
		want  usecases.Interval
		valid bool
	}{
		{
			name:  "day",
			p:     api.Interval_DAY,
			want:  usecases.Day,
			valid: true,
		},
		{
			name:  "week",
			p:     api.Interval_WEEK,
			want:  usecases.Week,
			valid: true,
		},
		{
			name:  "month",
			p:     api.Interval_MONTH,
			want:  usecases.Month,
			valid: true,
		},
		{
			name:  "year",
			p:     api.Interval_YEAR,
			want:  usecases.Year,
			valid: true,
		},
		{
			name:  "unknown",
			p:     -1,
			want:  usecases.InvalidInterval,
			valid: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IntervalFromProto(tt.p)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.valid, got.IsValid())
		})
	}
}
