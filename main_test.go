package main

import (
	"reflect"
	"testing"
)

func Test_parseOptions(t *testing.T) {
	type args struct {
		options []string
	}
	tests := []struct {
		name    string
		args    args
		want    *Options
		wantErr bool
	}{
		// {
		// 	name: "empty",
		// 	args: args{
		// 		options: []string{},
		// 	},
		// 	want: &Options{},
		// },
		{
			name: "suppress",
			args: args{
				options: []string{"suppress=Name"},
			},
			want: &Options{
				Suppress: []struct {
					Model string
					Field string
				}{
					{
						Model: "",
						Field: "Name",
					},
				},
				EntityFileDetect: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptions(tt.args.options)
			got.sets = nil
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
