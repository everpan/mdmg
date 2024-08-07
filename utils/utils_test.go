package utils

import "testing"

func TestSplitModuleVersion(t *testing.T) {
	tests := []struct {
		name  string
		arg   string
		want  string
		want1 string
	}{
		{"none", "", "", ""},
		{"echo", "echo", "echo", ""},
		{"only module", "moduleName", "moduleName", ""},
		{"only module", "moduleName-", "moduleName", ""},
		{"only module", "moduleName-subject", "moduleName-subject", ""},
		{"mod ver 1", "model-1.3", "model", "1.3"},
		{"mod ver 2", "model-sub1-1.3.9", "model-sub1", "1.3.9"},
		{"mod ver 3", "model-sub-1-1.3.9", "model-sub-1", "1.3.9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitModuleVersion(tt.arg)
			if got != tt.want {
				t.Errorf("SplitModuleVersion() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitModuleVersion() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
