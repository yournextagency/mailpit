package config

import "testing"

func TestDefaultMaxAge(t *testing.T) {
	if MaxAge != "30d" {
		t.Fatalf("expected a default max age of 30d, got %q", MaxAge)
	}

	defer func() { MaxAgeInHours = 0 }()

	if err := parseMaxAge(); err != nil {
		t.Fatal(err)
	}

	if MaxAgeInHours != 30*24 {
		t.Fatalf("expected 720 hours, got %d", MaxAgeInHours)
	}
}

func TestMaxAgeDisabled(t *testing.T) {
	orig := MaxAge
	defer func() {
		MaxAge = orig
		MaxAgeInHours = 0
	}()

	MaxAge = "0d"
	if err := parseMaxAge(); err != nil {
		t.Fatal(err)
	}

	if MaxAgeInHours != 0 {
		t.Fatalf("expected max age to be disabled, got %d hours", MaxAgeInHours)
	}
}
