package presenters

import (
	"brick/internal/adapters/dto/response"
	"math"
)

func calculatePagination(totalCount int, currentPage, pageSize, currentCount int) *response.PaginationResponse {
	if currentPage == 0 {
		currentPage = 1
	}

	nextPage := currentPage + 1

	switch {
	case pageSize > 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	firstPage := 1
	lastPage := totalPages

	if currentPage < firstPage {
		currentPage = firstPage
	}

	if currentPage > lastPage {
		currentPage = lastPage
	}

	if nextPage < firstPage {
		nextPage = firstPage
	}

	if nextPage > lastPage {
		nextPage = lastPage
	}

	return &response.PaginationResponse{
		CurrentPage:  currentPage,
		PageSize:     pageSize,
		TotalCount:   totalCount,
		TotalPages:   totalPages,
		FirstPage:    firstPage,
		NextPage:     nextPage,
		LastPage:     lastPage,
		CurrentCount: currentCount,
	}
}
