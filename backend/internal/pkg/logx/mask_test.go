package logx

import "testing"

func TestMaskMid(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"123", "****"},
		{"1234", "****"},
		{"16123456", "1612****"},
		{"1612345678", "1612****"},
	}
	for _, c := range cases {
		if got := MaskMid(c.in); got != c.want {
			t.Errorf("MaskMid(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskTx(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"abcd1234", "****"},
		{"abcde12345", "abcd****2345"},
		{"4200001234202501012345678901", "4200****8901"},
	}
	for _, c := range cases {
		if got := MaskTx(c.in); got != c.want {
			t.Errorf("MaskTx(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestMaskOpenID(t *testing.T) {
	if got := MaskOpenID("o4GgaunHGD4K7XHQh4u9Iy9FlfVc"); got != "o4Gg****lfVc" {
		t.Errorf("MaskOpenID mismatch, got %q", got)
	}
	if got := MaskOpenID(""); got != "" {
		t.Errorf("MaskOpenID empty should return empty, got %q", got)
	}
	if got := MaskOpenID("short"); got != "****" {
		t.Errorf("MaskOpenID short should return ****, got %q", got)
	}
}

func TestMaskPhone(t *testing.T) {
	if got := MaskPhone("13800138000"); got != "138****8000" {
		t.Errorf("MaskPhone mismatch, got %q", got)
	}
	if got := MaskPhone("123"); got != "***" {
		t.Errorf("MaskPhone short mismatch, got %q", got)
	}
}
