package paginate

var (
	DefaultMinPageSize uint64 = 10
	DefaultMaxPageSize uint64 = 100
)

type FilterParameter string

// Paginated struct represents a paginated response that will be passed to the repository layer
// to retrieve data based on pagination settings, filters, sorting, etc.
type Paginated struct {
	Page        uint64                     `json:"page"`
	PerPage     uint64                     `json:"per_page"`
	Total       uint64                     `json:"total"`
	Filters     map[FilterParameter]Filter `json:"filters"`
	SortColumn  string                     `json:"sort_column"`
	Decscending bool                       `json:"descending"`
}

/*
// Example of how the Paginated struct can be used in a database implementation:
// type PaginationSupportDB struct {
//     GetPaginated(ctx context.Context, p Paginated) (res []Result, total uint64, err error)
// }
// This interface defines how pagination and filtering data can be fetched from the database.
*/

type Filter struct {
	Operator FilterOperator
	Values   []interface{}
}

type FilterOperator int

const (
	FilterOperatorEqual FilterOperator = iota
	FilterOpratorNotEqual
	FilterOperatorGreater
	FilterOperatorGreaterEqual
	FilterOperatorLess
	FilterOperatorLessEqual
	FilterOperatorIN
	FilterOperatorNotIn
	FilterOperatorBetween
)

// the base structure for pagination requests from the client
type PaginateRequestBase struct {
	CurrentPage uint64                     `json:"current_page"`
	PageSize    uint64                     `json:"page_size"`
	Filters     map[FilterParameter]Filter `json:"filters"`
	SortColumn  string                     `json:"sort_column"`
	Decscending bool                       `json:"descending"`
}

// the base structure for paginated responses sent to the client
type PaginatedResponseBase struct {
	CurrentPage  uint64 `json:"current_page"`
	PageSize     uint64 `json:"page_size"`
	TotalNumbers uint64 `json:"total_numbers"`
	TotalPage    uint64 `json:"total_page"`
}

// BasicValidations just ensures that the pagination request is well-formed
// more complex validations should be done in the service layer based on the application's requirements
func (r *PaginateRequestBase) BasicValidations() error {
	if r.CurrentPage < 1 {
		r.CurrentPage = 1
	}
	if r.PageSize < DefaultMinPageSize {
		r.PageSize = DefaultMinPageSize
	}

	if r.PageSize > DefaultMaxPageSize {
		r.PageSize = DefaultMaxPageSize
	}

	// TODO: Add filters validation to ensure that the filters are well-formed

	return nil
}
