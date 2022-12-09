package util

import (
	"testing"
)

func TestHash(t *testing.T) {
	type args struct {
		passwd []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "string: ABCDabcd",
			args:    args{passwd: []byte("ABCDabcd")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hash(tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !VerifyEncoded(tt.args.passwd, got) {
				t.Error("VerifyEncoded() false, want true")
			}
		})
	}
}

func TestHashString(t *testing.T) {
	type args struct {
		passwd []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "string: ABCDabcd",
			args:    args{passwd: []byte("ABCDabcd")},
			wantErr: false,
		},
		{
			name:    "string: abcdABCD",
			args:    args{passwd: []byte("abcdABCD")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HashString(tt.args.passwd)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !VerifyEncoded(tt.args.passwd, []byte(got)) {
				t.Error("VerifyEncoded() false, want true")
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	GenerateKey()
	GenerateKey()
	GenerateKey()
	GenerateKey()
}
