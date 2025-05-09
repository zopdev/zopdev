package overprovision

import "testing"

func Test_getGCPCredentials(t *testing.T) {
	tests := []struct {
		name    string
		creds   any
		wantErr bool
	}{
		{
			name:    "valid creds",
			creds:   map[string]any{"type": "service_account", "project_id": "zopdev"},
			wantErr: false,
		},
		{
			name:    "invalid creds",
			creds:   `{"invalid_field":"invalid_value"}`,
			wantErr: true,
		},
		{
			name:    "nil creds",
			creds:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getGCPCredentials(tt.creds)

			if (err != nil) != tt.wantErr {
				t.Errorf("getGCPCredentials() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("getGCPCredentials() got = %v, want non-nil", got)
			}
		})
	}
}
