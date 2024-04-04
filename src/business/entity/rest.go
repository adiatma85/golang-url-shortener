package entity

const (
	SystemUser         = 0
	SchedulerUser      = -1
	TooManyConnections = "Error 1040: Too many connections"
)

type HTTPResp struct {
	Message    HTTPMessage `json:"message"`
	Meta       Meta        `json:"metadata"`
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type HTTPMessage struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Meta struct {
	Path        string     `json:"path"`
	StatusCode  int        `json:"statusCode"`
	Status      string     `json:"status"`
	Message     string     `json:"message"`
	Timestamp   string     `json:"timestamp"`
	Error       *MetaError `json:"error,omitempty"`
	RequestID   string     `json:"requestId"`
	TimeElapsed string     `json:"timeElapsed"`
}

type MetaError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Pagination struct {
	CurrentPage     int64    `json:"currentPage"`
	CurrentElements int64    `json:"currentElements"`
	TotalPages      int64    `json:"totalPages"`
	TotalElements   int64    `json:"totalElements"`
	SortBy          []string `json:"sortBy"`
	CursorStart     *string  `json:"cursorStart,omitempty"`
	CursorEnd       *string  `json:"cursorEnd,omitempty"`
}

func (p *Pagination) ProcessPagination(limit int64) {
	if p.SortBy == nil {
		p.SortBy = []string{}
	}

	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}

	if limit < 1 {
		limit = 10
	}

	totalPage := p.TotalElements / limit
	if p.TotalElements%limit > 0 || p.TotalElements == 0 {
		totalPage++
	}

	p.TotalPages = 1
	if totalPage > 1 {
		p.TotalPages = totalPage
	}
}

type PaginationParam struct {
	GroupBy           []string `param:"-" db:"-"`
	SortBy            []string `param:"sort_by" db:"sort_by"`
	Limit             int64    `form:"limit" param:"limit" db:"limit"`
	Page              int64    `form:"page" param:"page" db:"page"`
	IncludePagination bool
}

type Authorize struct {
	Param      string
	IsParam    string
	ActionCode string
}
