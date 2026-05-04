package timeedit

// ---------------- Imports ----------------

import (
	"math"
	"net/url"
	"strings"
)

// ---------------- Structures ----------------

type Pair struct {
	Key   string
	Value string
}

// ---------------- Pre-made arrays ----------------

var tabledata = []Pair{
	{"h=t&sid=", "6="},
	{"objects=", "1="},
	{"sid=", "2="},
	{"&ox=0&types=0&fe=0", "3=3"},
	{"&types=0&fe=0", "5=5"},
	{"&h=t&p=", "4="},
}

var tabledataspecial = []Pair{
	{"=", "ZZZX1"},
	{"&", "ZZZX2"},
	{",", "ZZZX3"},
	{".", "ZZZX4"},
	{" ", "ZZZX5"},
	{"-", "ZZZX6"},
	{"/", "ZZZX7"},
	{"%", "ZZZX8"},
}

var pairs = []Pair{
	{"=", "Q"},
	{"&", "Z"},
	{",", "X"},
	{".", "Y"},
	{" ", "V"},
	{"-", "W"},
}

var pattern = []int{4, 22, 5, 37, 26, 17, 33, 15, 39, 11, 45, 20, 2, 40, 19, 36, 28, 38, 30, 41, 44, 42, 7, 24, 14, 27, 35, 25, 12, 1, 43, 23, 6, 16, 3, 9, 47, 46, 48, 50, 21, 10, 49, 32, 18, 31, 29, 34, 13, 8}

// ---------------- Scramble Functions ----------------

func tablespecial(query string, reverse bool) string {
	for _, p := range tabledataspecial {
		if reverse {
			query = strings.ReplaceAll(query, p.Value, p.Key)
		}
		if !reverse {
			query = strings.ReplaceAll(query, p.Key, p.Value)
		}
	}
	return query
}

func tableshort(query string, reverse bool) string {
	for _, p := range tabledata {
		if reverse {
			query = strings.ReplaceAll(query, p.Value, p.Key)
		}
		if !reverse {
			query = strings.ReplaceAll(query, p.Key, p.Value)
		}
	}
	return query
}

func modKey(ch rune, reverse bool) rune {
	// a-z
	if ch >= 97 && ch <= 122 {
		if reverse {
			return rune(97 + (int(ch)-106+26)%26)
		}
		return rune(97 + (int(ch)-88)%26)
	}

	// 1-9
	if ch >= 49 && ch <= 57 {
		if reverse {
			return rune(49 + (int(ch)-53+9)%9)
		}
		return rune(49 + (int(ch)-45)%9)
	}
	return ch
}

func scrambleChar(ch rune, reverse bool) rune {
	char := string(ch)
	for _, pair := range pairs {
		if char == pair.Key {
			return []rune(pair.Value)[0]
		}
		if char == pair.Value {
			return []rune(pair.Key)[0]
		}
	}
	return modKey(ch, reverse)
}

func swap(source []rune, from, to int) {
	if from < 0 || to < 0 || from >= len(source) || to >= len(source) {
		return
	}
	source[from], source[to] = source[to], source[from]
}

func swapPattern(query []rune, reverse bool) {
	steps := int(math.Ceil(float64(len(query) / len(pattern))))
	for step := 0; step < steps; step++ {
		for i := 1; i < len(pattern); i += 2 {
			a := pattern[i-1] + step*len(pattern)
			b := pattern[i] + step*len(pattern)

			if reverse {
				swap(query, a, b)
			} else {
				swap(query, b, a)
			}
		}
	}
}

func swapChar(query string, reverse bool) string {
	chars := []rune(query)

	for i := range chars {
		chars[i] = scrambleChar(chars[i], reverse)
	}

	swapPattern(chars, reverse)

	return string(chars)
}

// ---------------- Scramble ----------------

func Scramble(query string, reverse bool) string {
	if len(query) < 2 {
		return query
	}
	if strings.HasPrefix(query, "i=") {
		return query
	}

	result, _ := url.QueryUnescape(query)

	steps := []func(string, bool) string{
		tableshort,
		func(q string, r bool) string { return swapChar(q, r) },
		tablespecial,
	}

	if reverse {
		for i := len(steps) - 1; i >= 0; i-- {
			result = steps[i](result, reverse)
		}
	} else {
		for _, step := range steps {
			result = step(result, reverse)
		}
	}

	return url.QueryEscape(result)
}
