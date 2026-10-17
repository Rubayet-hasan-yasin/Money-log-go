package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnvDuration(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envVal       string
		fallbackKeys []string
		defaultVal   time.Duration
		expected     time.Duration
	}{
		{
			name:       "default value when unset",
			envKey:     "TEST_DUR_UNSET",
			envVal:     "",
			defaultVal: 10 * time.Minute,
			expected:   10 * time.Minute,
		},
		{
			name:       "standard duration minutes",
			envKey:     "TEST_DUR_MIN",
			envVal:     "10m",
			defaultVal: 5 * time.Minute,
			expected:   10 * time.Minute,
		},
		{
			name:       "duration in days suffix d",
			envKey:     "TEST_DUR_DAYS_D",
			envVal:     "7d",
			defaultVal: 24 * time.Hour,
			expected:   7 * 24 * time.Hour,
		},
		{
			name:       "duration in days suffix day",
			envKey:     "TEST_DUR_DAYS_DAY",
			envVal:     "7day",
			defaultVal: 24 * time.Hour,
			expected:   7 * 24 * time.Hour,
		},
		{
			name:       "duration in days suffix days",
			envKey:     "TEST_DUR_DAYS_DAYS",
			envVal:     "30days",
			defaultVal: 24 * time.Hour,
			expected:   30 * 24 * time.Hour,
		},
		{
			name:       "duration in hours",
			envKey:     "TEST_DUR_HOURS",
			envVal:     "168h",
			defaultVal: 1 * time.Hour,
			expected:   168 * time.Hour,
		},
		{
			name:         "fallback key matched",
			envKey:       "TEST_DUR_PRIMARY_NOT_SET",
			envVal:       "",
			fallbackKeys: []string{"TEST_DUR_FALLBACK"},
			defaultVal:   10 * time.Minute,
			expected:     7 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envVal != "" {
				_ = os.Setenv(tt.envKey, tt.envVal)
				defer os.Unsetenv(tt.envKey)
			}
			if len(tt.fallbackKeys) > 0 {
				_ = os.Setenv(tt.fallbackKeys[0], "7d")
				defer os.Unsetenv(tt.fallbackKeys[0])
			}

			result := getEnvDuration(tt.envKey, tt.defaultVal, tt.fallbackKeys...)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
