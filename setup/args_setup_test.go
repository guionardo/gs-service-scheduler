package setup

import (
	"os"
	"reflect"
	"testing"
)

func Test_getArgs(t *testing.T) {
	tmp, err := os.MkdirTemp(os.TempDir(), "gs-service-scheduler")
	if err != nil {
		t.Errorf("getArgs() error = %v", err)
	}
	defer os.RemoveAll(tmp)

	tests := []struct {
		name     string
		runArgs  []string
		wantArgs *Args
		wantErr  bool
	}{{
		name:     "Valid",
		runArgs:  []string{"--setup-folder=" + tmp},
		wantArgs: &Args{SetupFolder: tmp},
		wantErr:  false,
	}, {
		name:     "No setup folder",
		runArgs:  []string{},
		wantArgs: &Args{},
		wantErr:  true,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, err := getArgs(tt.runArgs)
			if (err != nil) != tt.wantErr {
				t.Errorf("getArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("getArgs() = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}
