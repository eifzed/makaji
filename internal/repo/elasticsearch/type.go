package elasticsearch

type Query struct {
	Bool Bool `json:"bool"`
}

type Search struct {
	Query Query `json:"query"`
	From  int   `json:"from"`
	Size  int   `json:"size"`
}

type Bool struct {
	Should             []Should `json:"should"`
	MinimumShouldMatch int      `json:"minimum_should_match,omitempty"`
}

type Should struct {
	Match    map[string]string   `json:"match,omitempty"`
	Wildcard map[string]Wildcard `json:"wildcard,omitempty"`
}

type Wildcard struct {
	Value           interface{} `json:"value"`
	Boost           float32     `json:"boost"`
	Rewrite         string      `json:"rewrite,omitempty"`
	CaseInsensitive bool        `json:"case_insensitive"`
}

func (q *Search) SetPagination(page, limit int) {
	q.From = (page - 1) * limit
	q.Size = limit
}
