package main

import "testing"

func TestCase(t *testing.T) {
	if s := toSnakeCase("TestTrois"); s != "test_trois" {
		t.Errorf("got %s", s)
	}
}
