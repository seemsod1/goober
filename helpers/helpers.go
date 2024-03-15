package helpers

import "crypto/rand"

func CheckPassword(pass string, tocmp string) bool {
	return pass == tocmp
}

func GeneratePageNumbers(currentPage, totalPages int) []int {
	var pages []int

	maxPages := 5
	start := currentPage - maxPages/2
	end := currentPage + maxPages/2

	if start < 1 {
		start = 1
	}

	if end > totalPages {
		end = totalPages
	}

	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}

	return pages
}

var chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func GenerateRandomString(length int) string {
	ll := len(chars)
	b := make([]byte, length)
	rand.Read(b)
	for i := 0; i < length; i++ {
		b[i] = chars[int(b[i])%ll]
	}
	return string(b)
}
