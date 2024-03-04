package helpers

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
