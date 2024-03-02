package models

type PaginationData struct {
	CurrentPage int   // Current page number
	TotalPages  int   // Total number of pages
	PrevPage    int   // Previous page number
	NextPage    int   // Next page number
	HasPrev     bool  // Whether there is a previous page
	HasNext     bool  // Whether there is a next page
	Pages       []int // List of page numbers for pagination links
}
