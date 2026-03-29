package logger

import (
	"testing"

	"github.com/EliasBlind/EduFlow/internal/journal_service/config"
)

func TestMustLoad(t *testing.T) {
	tests := []struct {
		name      string
		env       config.Env
		wantPanic bool
	}{
		{
			name:      "Susses start for Local",
			env:       config.EnvLocal,
			wantPanic: false,
		},
		{
			name:      "Susses start for Dev",
			env:       config.EnvDev,
			wantPanic: false,
		},
		{
			name:      "Susses start for Prod",
			env:       config.EnvProd,
			wantPanic: false,
		},
		{
			name:      "Panic where unknown env",
			env:       "unknown_env",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MustLoad() panic = %v, wantPanic %v", r, tt.wantPanic)
				}
			}()

			log := MustLoad(tt.env)

			if !tt.wantPanic && log == nil {
				t.Error("MustLoad() returned nil even though it shouldn't be like this")
			}
		})
	}
}
