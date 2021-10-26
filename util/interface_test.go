package util

import (
	"reflect"
	"testing"
)

func TestMapToStruct(t *testing.T) {
	type args struct {
		data   map[string]interface{}
		result interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test map to struct",
			want: nil,
			args: args{
				data: map[string]interface{}{"string": "is string"},
				result: "is string",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapToStruct(tt.args.data, tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToStruct() = %v, want %v", got, tt.want)
			}
		})
	}
}
