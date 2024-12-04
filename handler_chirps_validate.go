package main

import (
	"strings"
)

func profaneFilter(msg string) string {
	splitted := strings.Split(msg, " ")
	for i, word := range splitted {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			splitted[i] = "****"
		}
	}
	return strings.Join(splitted, " ")
}
