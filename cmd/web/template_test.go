package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	// Given
	tm := time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC)
	expected := "17 Mar 2022 at 10:15"
	// When
	hd := humanDate(tm)
	// Then
	if hd != expected {
		t.Errorf("got %q; want %q", hd, expected)
	}
}
