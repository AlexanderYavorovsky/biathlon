package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		want       Config
		wantErr    bool
		wantErrStr string
	}{
		{
			name: "valid_config",
			input: `
{
	"laps": 2,
	"lapLen": 3500,
	"penaltyLen": 150,
	"firingLines": 2,
	"start": "10:00:00.000",
	"startDelta": "00:01:30"
}`,
			want: Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       Time{time.Date(0, 1, 1, 10, 0, 0, 0, time.UTC)},
				StartDelta:  Duration{time.Minute + 30*time.Second},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)
			tmpFile, _ := os.CreateTemp("", tt.name)
			defer os.Remove(tmpFile.Name())
			tmpFile.Write([]byte(tt.input))

			got, err := parseConfig(tmpFile.Name())
			if tt.wantErr {
				assert.EqualError(err, tt.wantErrStr)
			} else {
				assert.Equal(tt.want.Laps, got.Laps)
				assert.Equal(tt.want.LapLen, got.LapLen)
				assert.Equal(tt.want.PenaltyLen, got.PenaltyLen)
				assert.Equal(tt.want.FiringLines, got.FiringLines)
				assert.Equal(tt.want.Start, got.Start)
				assert.Equal(tt.want.StartDelta, got.StartDelta)
			}
		})
	}
}
