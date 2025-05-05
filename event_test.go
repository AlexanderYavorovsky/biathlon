package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseEvent(t *testing.T) {
	testTime := time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		line       string
		want       Event
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "valid event",
			line: "[10:00:00.000] 11 2 Extra params",
			want: Event{
				ID:           11,
				CompetitorID: 2,
				Time:         testTime,
				ExtraParams:  []string{"Extra", "params"},
			},
		},
		{
			name:       "invalid event format",
			line:       "[10:00:00.000] 11",
			wantErr:    true,
			wantErrStr: "invalid line format: [10:00:00.000] 11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			got, err := parseEvent(tt.line)
			if tt.wantErr {
				assert.NotNil(err)
				assert.EqualError(err, tt.wantErrStr)
			} else {
				assert.Nil(err)
				assert.Equal(tt.want, got)
			}
		})
	}
}
