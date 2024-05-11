package elasticsearch

type Query struct {
	Bool  *Bool            `json:"bool,omitempty"`
	Range map[string]Range `json:"range,omitempty"`
}

type Range struct {
	GreaterThanOrEqualTo float64 `json:"gte,omitempty"`
	LessThanOrEqualTo    float64 `json:"lte,omitempty"`
	LessThan             float64 `json:"lt,omitempty"`
	GreaterThan          float64 `json:"gt,omitempty"`
	Boost                float32 `json:"boost,omitempty"`
	Format               string  `json:"format,omitempty"`
}

type Search struct {
	Query Query  `json:"query"`
	From  uint32 `json:"from"`
	Size  uint32 `json:"size"`
}

func (q *Query) AddRange(field string, rangeQ Range) {
	if q.Range == nil {
		q.Range = map[string]Range{}
	}
	existingRange, ok := q.Range[field]

	if !ok {
		q.Range[field] = rangeQ
		return
	}

	if rangeQ.GreaterThan > 0 {
		existingRange.GreaterThan = rangeQ.GreaterThan
	}
	if rangeQ.LessThan > 0 {
		existingRange.LessThan = rangeQ.LessThan
	}
	if rangeQ.GreaterThanOrEqualTo > 0 {
		existingRange.GreaterThanOrEqualTo = rangeQ.GreaterThanOrEqualTo
	}
	if rangeQ.LessThanOrEqualTo > 0 {
		existingRange.LessThanOrEqualTo = rangeQ.LessThanOrEqualTo
	}
	if rangeQ.Format != "" {
		existingRange.Format = rangeQ.Format
	}
	if rangeQ.Boost > 0 {
		existingRange.Boost = rangeQ.Boost
	}
	q.Range[field] = existingRange
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

func (q *Search) SetPagination(page, limit uint32) {
	q.From = (page - 1) * limit
	q.Size = limit
}
