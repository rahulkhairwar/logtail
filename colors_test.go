package logtail

import (
	"image/color"
	"reflect"
	"testing"
)

func TestColor_HexCode(t *testing.T) {
	colors := getColors()
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{want: Yellow.hexCode},
		{want: Green.hexCode},
		{want: Red.hexCode},
		{want: Brown.hexCode},
	}
	for i := 0; i < len(colors); i++ {
		tests[i].name = colors[i].name
		tests[i].fields = colors[i].fields
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Color{
				hexCode: tt.fields.hexCode,
				rgba:    tt.fields.rgba,
			}
			if got := c.HexCode(); got != tt.want {
				t.Errorf("HexCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColor_RGBA(t *testing.T) {
	colors := getColors()
	tests := []struct {
		name   string
		fields fields
		want   color.RGBA
	}{
		{want: Yellow.rgba},
		{want: Green.rgba},
		{want: Red.rgba},
		{want: Brown.rgba},
	}
	for i := 0; i < len(colors); i++ {
		tests[i].name = colors[i].name
		tests[i].fields = colors[i].fields
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Color{
				hexCode: tt.fields.hexCode,
				rgba:    tt.fields.rgba,
			}
			if got := c.RGBA(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RGBA() = %v, want %v", got, tt.want)
			}
		})
	}
}

type colorTest struct {
	name   string
	fields fields
}

type fields struct {
	hexCode string
	rgba    color.RGBA
}

func getColors() []colorTest {
	return []colorTest{
		{
			name: "color yellow",
			fields: fields{
				hexCode: "FFFF00",
				rgba: color.RGBA{
					R: 255,
					G: 255,
					B: 0,
				},
			},
		},
		{
			name: "color green",
			fields: fields{
				hexCode: "00FF00",
				rgba: color.RGBA{
					R: 0,
					G: 255,
					B: 0,
				},
			},
		},
		{
			name: "color red",
			fields: fields{
				hexCode: "FF0000",
				rgba: color.RGBA{
					R: 255,
					G: 0,
					B: 0,
				},
			},
		},
		{
			name: "color brown",
			fields: fields{
				hexCode: "663300",
				rgba: color.RGBA{
					R: 102,
					G: 51,
					B: 0,
				},
			},
		},
	}
}
