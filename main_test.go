//go:generate goversioninfo -icon=./resource/icon.ico ./resource/versioninfo.json

package main

import (
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Test_getBinaryDirectory(t *testing.T) {
	type args struct {
		args []string
	}
	type test struct {
		name    string
		args    args
		wantDir string
		wantErr bool
	}

	currentDir, _ := filepath.Abs(filepath.Dir(""))

	tests := []test{
		{
			name:    "empty input",
			args:    args{[]string{}},
			wantDir: "",
			wantErr: true,
		},
		{
			name:    "empty string",
			args:    args{[]string{""}},
			wantDir: currentDir,
			wantErr: false,
		},
	}

	if runtime.GOOS == "windows" {
		windows := test{
			name:    "test path windows",
			args:    args{[]string{"C:\\test\\main.exe"}},
			wantDir: "C:\\test",
			wantErr: false,
		}
		tests = append(tests, windows)
	}

	if runtime.GOOS == "linux" {
		linux := test{
			name:    "test path linux",
			args:    args{[]string{"/tmp/test/main"}},
			wantDir: "/tmp/test",
			wantErr: false,
		}
		tests = append(tests, linux)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir, err := getBinaryDirectory(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("getBinaryDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDir != tt.wantDir {
				t.Errorf("getBinaryDirectory() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func Test_sortConfigBrowserPriority(t *testing.T) {
	type args struct {
		input []domain
	}
	type test struct {
		name       string
		args       args
		wantOutput []domain
		wantErr    bool
	}

	tests := []test{
		{
			name:       "empty input",
			args:       args{[]domain{}},
			wantOutput: []domain{},
			wantErr:    false,
		},
		{
			name:       "sort",
			args:       args{[]domain{{"2", "", 20}, {"1", "", 10}}},
			wantOutput: []domain{{"1", "", 10}, {"2", "", 20}},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOutput, err := sortConfigBrowserPriority(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortConfigBrowserPriority() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("sortConfigBrowserPriority() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func Test_debug(t *testing.T) {
	type args struct {
		debug bool
		a     []interface{}
	}
	type test struct {
		name    string
		args    args
		wantN   int // bytes written to stdout
		wantErr bool
	}

	logInterface := make([]interface{}, 1)
	logInterface[0] = "log"

	tests := []test{
		{
			name:    "debug on",
			args:    args{debug: true, a: logInterface},
			wantN:   4,
			wantErr: false,
		},
		{
			name:    "debug off",
			args:    args{debug: false, a: logInterface},
			wantN:   0,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotN, err := debug(tt.args.debug, tt.args.a...)
			if (err != nil) != tt.wantErr {
				t.Errorf("debug() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("debug() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_getUrl(t *testing.T) {
	type args struct {
		args   []string
		config configuration
	}
	type test struct {
		name    string
		args    args
		wantUrl string
		wantErr bool
	}

	defaultConfig := configuration{Debug: false}

	tests := []test{
		{
			name:    "empty parameters",
			args:    args{args: []string{}, config: defaultConfig},
			wantUrl: "",
			wantErr: true,
		},
		{
			name:    "invalid parameters",
			args:    args{args: []string{"1", "2"}, config: defaultConfig},
			wantUrl: "",
			wantErr: true,
		},
		{
			name:    "url discovery",
			args:    args{args: []string{"1", "2", "http://a.b", "4"}, config: defaultConfig},
			wantUrl: "http://a.b",
			wantErr: false,
		},
		{
			name:    "url with path",
			args:    args{args: []string{"http://a.b/path"}, config: defaultConfig},
			wantUrl: "http://a.b/path",
			wantErr: false,
		},
		{
			name:    "url with complex path",
			args:    args{args: []string{"http://a.b/path?a=1234#1234"}, config: defaultConfig},
			wantUrl: "http://a.b/path?a=1234#1234",
			wantErr: false,
		},
		{
			name:    "url with spaces",
			args:    args{args: []string{"http://a.b/path with spaces/"}, config: defaultConfig},
			wantUrl: "http://a.b/path with spaces/",
			wantErr: false,
		},
		{
			name:    "file url",
			args:    args{args: []string{"file:///A:/b/file.pdf"}, config: defaultConfig},
			wantUrl: "file:///A:/b/file.pdf",
			wantErr: false,
		},
		{
			name:    "file url with spaces",
			args:    args{args: []string{"file:///A:/b/file with spaces.pdf"}, config: defaultConfig},
			wantUrl: "file:///A:/b/file with spaces.pdf",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrl, err := getUrl(tt.args.args, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUrl != tt.wantUrl {
				t.Errorf("getUrl() = %v, want %v", gotUrl, tt.wantUrl)
			}
		})
	}
}

func Test_getFqdnFromUrl(t *testing.T) {
	type args struct {
		url    string
		config configuration
	}
	type test struct {
		name         string
		args         args
		wantProtocol string
		wantFqdn     string
		wantErr      bool
	}

	defaultConfig := configuration{Debug: false}

	tests := []test{
		{
			name:         "invalid url",
			args:         args{url: "a.b", config: defaultConfig},
			wantProtocol: "",
			wantFqdn:     "",
			wantErr:      true,
		},
		{
			name:         "get http fqdn",
			args:         args{url: "http://a.b", config: defaultConfig},
			wantProtocol: "http",
			wantFqdn:     "a.b",
			wantErr:      false,
		},
		{
			name:         "get https fqdn",
			args:         args{url: "https://a.b/", config: defaultConfig},
			wantProtocol: "https",
			wantFqdn:     "a.b",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotProtocol, gotFqdn, err := getFqdnFromUrl(tt.args.url, tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("getFqdnFromUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotProtocol != tt.wantProtocol {
				t.Errorf("getFqdnFromUrl() gotProtocol = %v, want %v", gotProtocol, tt.wantProtocol)
			}
			if gotFqdn != tt.wantFqdn {
				t.Errorf("getFqdnFromUrl() gotFqdn = %v, want %v", gotFqdn, tt.wantFqdn)
			}
		})
	}
}
