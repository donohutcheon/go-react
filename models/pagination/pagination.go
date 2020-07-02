package pagination

import (
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response/types"
)

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "asc"
	SortDirectionDesc SortDirection = "desc"
)

type Parameters struct {
	FetchFrom  string
	Page       int64
	FetchCount int64
	SortField  string
	SortDir    SortDirection
	isInfinite bool
}

type Sortable interface {
	GetSortFields() map[string]bool
	GetPagination() Parameters
	SetSortParameters(Parameters)
}

func ParsePagination(logger *log.Logger, queryParams url.Values, entity Sortable) error {
	page := int64(10)
	fetchCount := int64(10)
	fetchFrom := "0"
	sortField := "id"
	sortDir := SortDirectionAsc
	isInfinite := true
	var err error

	if _, ok := queryParams["from"]; ok {
		fetchFrom = queryParams.Get("from")
		if err != nil {
			return err
		}
	}

	if _, ok := queryParams["page"]; ok {
		page, err = strconv.ParseInt(queryParams.Get("page"), 10, 64)
		if err != nil {
			return err
		}
		if page < 0 {
			fields := []types.ErrorField{
				{
					Name:    "page",
					Message: "negative page value",
					Direct:  true,
				},
			}
			return errors.NewError("invalid pagination parameters", fields, http.StatusBadRequest )
		}

		isInfinite = false
	}

	if _, ok := queryParams["count"]; ok {
		fetchCount, err = strconv.ParseInt(queryParams.Get("count"), 10, 64)
		if err != nil {
			return err
		}
		if page < 0 {
			fields := []types.ErrorField{
				{
					Name:    "count",
					Message: "invalid page count value",
					Direct:  true,
				},
			}
			return errors.NewError("invalid pagination parameters", fields, http.StatusBadRequest )
		}
	}

	if _, ok := queryParams["sortField"]; ok {
		sortField = queryParams.Get("sortField")
		_, ok := entity.GetSortFields()[sortField]
		if !ok {
			fields := []types.ErrorField{
				{
					Name:    "sortField",
					Message: "invalid sort field",
					Direct:  true,
				},
			}
			return errors.NewError("invalid sort field", fields, http.StatusBadRequest )
		}

	}

	if _, ok := queryParams["sortDir"]; ok {
		sortDir = SortDirection(queryParams.Get("sortDir"))
		if sortDir != SortDirectionAsc && sortDir != SortDirectionDesc {
			fields := []types.ErrorField{
				{
					Name:    "sortDir",
					Message: "invalid sort direction",
					Direct:  true,
				},
			}
			return errors.NewError("invalid sort direction", fields, http.StatusBadRequest )
		}
	}

	entity.SetSortParameters(
		Parameters{
			FetchFrom:  fetchFrom,
			Page:       page,
			FetchCount: fetchCount,
			SortField:  sortField,
			SortDir:    sortDir,
			isInfinite: isInfinite,
		})
	logger.Printf("Pagination: { page: %d, count: %d}", page, fetchCount)
	return nil
}
