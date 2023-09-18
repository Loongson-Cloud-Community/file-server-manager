package main

import (
	"testing"
)

func TestParseUrl(t *testing.T) {
	type args struct {
		url string
	}

	tests := []struct {
		args    args
		name    string
		want1   string
		want2   string
		wantErr bool
	}{
		{
			args: args{
				url: "org/repo/ver/file",
			},
			name:    "1",
			want1:   "org/repo/ver",
			want2:   "file",
			wantErr: false,
		},
		{
			args: args{
				url: "org/repo/ver/file/",
			},
			name:    "1",
			want1:   "org/repo/ver",
			want2:   "file",
			wantErr: false,
		},
		{
			args: args{
				url: "file",
			},
			name:    "1",
			want1:   "",
			want2:   "",
			wantErr: true,
		},
		{
			args: args{
				url: "file/",
			},
			name:    "1",
			want1:   "",
			want2:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, file, err := parseUrl(tt.args.url)
			if err != nil {
				t.Logf("%v", err)
			}
			if (err != nil) != tt.wantErr || dir != tt.want1 || file != tt.want2 {
				t.Errorf("parseUrl(%s) got %s,%s; want %s,%s", tt.args.url, dir, file, tt.want1, tt.want2)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	in := []byte("Hello, World!")
	key := GenerateKey(in)
	t.Logf("%v", key)
}
