package index

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const pageKey = "p"
const ippKey = "ipp"
const defaultIpp = 32
const firstPage = 1

type pagination struct {
	CurrentPage             int
	TotalPages              int
	ItemsPerPage            int
	requestPath             string
	requestValuesWitoutPage url.Values
}

func newPagination(r *http.Request) pagination {
	reqQuery := r.URL.Query()
	reqQuery.Del(pageKey)

	ipp := getIpp(r)

	return pagination{
		CurrentPage:             getPage(r),
		TotalPages:              0, // set afterwards
		ItemsPerPage:            ipp,
		requestPath:             r.URL.Path,
		requestValuesWitoutPage: reqQuery,
	}
}

func getPage(r *http.Request) int {
	urlPage := intInURLOrDefault(r, pageKey, firstPage)
	if urlPage < firstPage {
		urlPage = firstPage
	}
	return urlPage
}

func getIpp(r *http.Request) int {
	return intInURLOrDefault(r, ippKey, defaultIpp)
}

func intInURLOrDefault(r *http.Request, key string, def int) int {
	if p, err := strconv.Atoi(r.URL.Query().Get(key)); err == nil {
		return p
	}
	return def
}

func (p *pagination) PageURL(page int) string {
	// copy values to not alter the original request
	values := make(url.Values)
	for k, v := range p.requestValuesWitoutPage {
		values[k] = v
	}

	values.Add(pageKey, fmt.Sprint(page))
	return p.requestPath + "?" + values.Encode()
}

func (p *pagination) PreviousPageURL() string {
	return p.PageURL(p.CurrentPage - 1)
}

func (p *pagination) FirstPageURL() string {
	return p.PageURL(firstPage)
}

func (p *pagination) NextPageURL() string {
	return p.PageURL(p.CurrentPage + 1)
}

func (p *pagination) LastPageURL() string {
	return p.PageURL(p.TotalPages)
}

func (p *pagination) IsFirstPage() bool {
	return p.CurrentPage <= firstPage
}

func (p *pagination) IsLastPage() bool {
	return p.CurrentPage >= p.TotalPages
}

// for easier looping in template
func (p *pagination) PageList() []int {
	l := make([]int, p.TotalPages)
	for i := range l {
		l[i] = i + 1
	}
	return l
}

func (p *pagination) IsNecessary() bool {
	return p.TotalPages > 1
}
