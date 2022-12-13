package server

import (
	"github.com/spf13/viper"
	"reflect"
	"testing"
)

func Test_parseViper(t *testing.T) {
	viperr := viper.New()

	type args struct {
		*Config
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *Config
	}{
		{
			name: "without keys",
			args: args{Config: &Config{
				v:          nil,
				Address:    "0.0.0.0:3500",
				DBHost:     "db",
				DBPort:     5432,
				DBName:     "test",
				DBUsername: "username",
				DBPassword: "sec23tandsacr3d",
			}},
			wantConfig: &Config{
				v:          viperr,
				Address:    "0.0.0.0:3500",
				DBHost:     "db",
				DBPort:     5432,
				DBName:     "test",
				DBUsername: "username",
				DBPassword: "sec23tandsacr3d",
			},
		},
		{
			name: "with key string",
			args: args{Config: &Config{
				v:   nil,
				Key: "abc",
			}},
			wantConfig: &Config{
				v:        viperr,
				Key:      "abc",
				KeyBytes: []byte("abc"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			config := parseViper(viperr, tt.args.Config)

			if !reflect.DeepEqual(*config, *tt.wantConfig) {
				t.Fatalf("parseConfig() got: %+v, want: %+v", *config, *tt.wantConfig)
			}
		})
	}
}

func Test_getOrDefault_int(t *testing.T) {
	type args[T interface{}] struct {
		val        T
		defaultVal T
	}
	type testCase[T interface{}] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[int]{
		{
			name: "set 123 default 100",
			args: args[int]{val: 123, defaultVal: 100},
			want: 123,
		},
		{
			name: "set 0 default 100",
			args: args[int]{val: 0, defaultVal: 100},
			want: 100,
		},
		{
			name: "set 0 default 0",
			args: args[int]{val: 0, defaultVal: 0},
			want: 00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOrDefault(tt.args.val, tt.args.defaultVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOrDefault_string(t *testing.T) {
	type args[T interface{}] struct {
		val        T
		defaultVal T
	}
	type testCase[T interface{}] struct {
		name string
		args args[T]
		want T
	}
	tests := []testCase[string]{
		{
			name: "set user default root",
			args: args[string]{val: "user", defaultVal: "root"},
			want: "user",
		},
		{
			name: "set '' default root",
			args: args[string]{val: "", defaultVal: "root"},
			want: "root",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOrDefault(tt.args.val, tt.args.defaultVal); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfig_Apply(t *testing.T) {
	type fields struct {
		*Config
	}
	type args struct {
		conf *Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "apply Address: 0.0.0.0",
			fields: fields{
				Config: &Config{
					Address: "0.0.0.0",
				},
			},
			args: args{
				conf: &Config{},
			},
			wantErr: false,
		},
		{
			name: "apply DBName: db0",
			fields: fields{
				Config: &Config{
					DBName: "db0",
				},
			},
			args: args{
				conf: &Config{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				v:          tt.fields.v,
				Address:    tt.fields.Address,
				DBHost:     tt.fields.DBHost,
				DBPort:     tt.fields.DBPort,
				DBName:     tt.fields.DBName,
				DBUsername: tt.fields.DBUsername,
				DBPassword: tt.fields.DBPassword,
			}
			if err := c.Apply(tt.args.conf); (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(*c, *tt.args.conf) {
				t.Fatalf("Apply want: %+v. got: %+v", *c, *tt.args.conf)
			}
		})
	}
}
