package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
	"strings"
	"testing"
	"time"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "help",
			expected: "Welcome to the Pokedex!",
		},
	}

	for _, c := range cases {

		cmd := "hello"
		r, w, _ := os.Pipe() // make a pipe with read and write variables
		os.Stdout = w        // replace stdout with pipe "w"

		if c.input == "help" {
			commands["help"].callback(&pokeapi.Client{}, cmd)

			// Read the output
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)

			// Check if output contains expected string
			if !strings.Contains(buf.String(), c.expected) {
				t.Errorf("Expected output to contain %q", c.expected)
			}
		}

	}
}

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := pokecache.NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestGet(t *testing.T) {
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "www.arronisfat.com",
			val: []byte("testing"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}
	for _, c := range cases {
		interval := 5 * time.Second
		cache := pokecache.NewCache(interval)
		cache.Add(c.key, c.val)
		value, res := cache.Get(c.key)
		if !res {
			t.Errorf("response not good")
			return
		}
		if string(c.val) != string(value) {
			t.Errorf("values are not the same")
			return
		}
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := pokecache.NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestReapLoop1(t *testing.T) {
	interval := 5 * time.Millisecond
	cache := pokecache.NewCache(interval)
	cache.Add("https://example.com", []byte("testdata"))
	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
	}

	time.Sleep(interval)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
	}
}
