package util

import (
	"reflect"
	"testing"
)

func TestBytesToString(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "B2S 1",
			args: args{b: []byte{'a', '0', 65, 'p', 'D'}},
			want: "a0ApD",
		},
		{
			name: "B2S 2",
			args: args{b: []byte{65, 66, 'C', '#', '+'}},
			want: "ABC#+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToString(tt.args.b); got != tt.want {
				t.Errorf("BytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringToBytes(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		wantB []byte
	}{
		{
			name:  "S2B 1",
			args:  args{s: "a0ApD"},
			wantB: []byte{'a', '0', 65, 'p', 'D'},
		},
		{
			name:  "S2B 2",
			args:  args{s: "ABC#+"},
			wantB: []byte{65, 66, 'C', '#', '+'},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := StringToBytes(tt.args.s); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("StringToBytes() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
