package logtail

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestHandleLogLine(t *testing.T) {
	log.SetFlags(0)

	type args struct {
		conf *Config
		l    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty key-set",
			args: args{
				conf: &Config{
					ParseKeys: nil,
				},
			},
			want: "failed to convert string to JSON, err:  unexpected end of JSON input\n",
		},
		{
			name: "no matching key",
			args: args{
				conf: &Config{
					ParseKeys: []string{"message2"},
				},
				l: `{
					"fruit": "Apple",
					"color": "Red"
				}`,
			},
			want: separator + "\n",
		},
		{
			name: "matching key 'fruit'",
			args: args{
				conf: &Config{
					ParseKeys: []string{"fruit"},
				},
				l: `{
					"fruit": "Apple",
					"color": "Red"
				}`,
			},
			want: "fruit = Apple\n" + separator + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.Buffer{}

			log.SetOutput(&b)
			HandleLogLine(tt.args.conf, tt.args.l)
			assert.Equal(t, tt.want, b.String())
		})
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(defaultFlags)
}

func Test_contains(t *testing.T) {
	_key := "key_1"
	_arr := []string{"key_2", _key, "key_3"}

	type args struct {
		arr []string
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil array",
			args: args{
				arr: nil,
				key: _key,
			},
			want: false,
		},
		{
			name: "empty array",
			args: args{
				arr: nil,
				key: _key,
			},
			want: false,
		},
		{
			name: "empty key",
			args: args{
				arr: _arr,
				key: "",
			},
			want: false,
		},
		{
			name: "array does not contain key",
			args: args{
				arr: _arr,
				key: "hello",
			},
			want: false,
		},
		{
			name: "array contains key",
			args: args{
				arr: _arr,
				key: _key,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, contains(tt.args.arr, tt.args.key))
		})
	}
}

func Test_printSelectedKeys(t *testing.T) {
	_mp := make(map[string]interface{})

	_mp["time"] = time.Now()
	_mp["message"] = "This is a sample message."
	log.SetFlags(0)

	type args struct {
		conf *Config
		m    map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty key-set",
			args: args{
				conf: &Config{
					ParseKeys: nil,
				},
				m: _mp,
			},
			want: separator + "\n",
		},
		{
			name: "no matching key",
			args: args{
				conf: &Config{
					ParseKeys: []string{"message2"},
				},
				m: _mp,
			},
			want: separator + "\n",
		},
		{
			name: "matching key 'message'",
			args: args{
				conf: &Config{
					ParseKeys: []string{"message"},
				},
				m: _mp,
			},
			want: "message = This is a sample message.\n" + separator + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.Buffer{}

			log.SetOutput(&b)
			printSelectedKeys(tt.args.conf, tt.args.m)
			assert.Equal(t, tt.want, b.String())
		})
	}

	log.SetOutput(os.Stderr)
	log.SetFlags(defaultFlags)
}

func Test_stringToJSON(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "empty json",
			wantErr: true,
		},
		{
			name: "valid json object",
			args: args{
				s: `{
					"fruit": "Apple",
					"color": "Red"
				}`,
			},
			want:    map[string]interface{}{"fruit": "Apple", "color": "Red"},
			wantErr: false,
		},
		{
			name: "invalid json object",
			args: args{
				s: `{
					"fruit": "Apple",
					"color": "Red",
				}`,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToJSON(tt.args.s)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("stringToJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToJSON() got = %v, want %v", got, tt.want)
			}
		})
	}
}
