// Copyright (c) 2023 Andy Fusniak. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

// Package base58 provides a cryptographically secure random base58
// string generator.
//
// Base58 has an alphabet of 58 easily readable characters. It excludes
// letters that might look ambiguous when printed (0 – zero,
// I – capital i, O – capital o and l – lower-case L).
// Unlike Base64 it does not contain any URI reserved characters so is
// suitable for use in URL query parameters.
package base58

import (
	"crypto/rand"
)

var alphabet = [58]string{
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "A",
	"B", "C", "D", "E", "F", "G", "H", "J", "K", "L",
	"M", "N", "P", "Q", "R", "S", "T", "U", "V", "W",
	"X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g",
	"h", "i", "j", "k", "m", "n", "o", "p", "q", "r",
	"s", "t", "u", "v", "w", "x", "y", "z",
}

// RandString generates a cryptographically secure random base58 string
// of fixed length n.
func RandString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	var s string
	for _, j := range b {
		idx := int(j) % 58
		s = s + alphabet[idx]
	}
	return s, nil
}
