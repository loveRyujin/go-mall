package app

import "encoding/json"

type Pagination struct {
	page     int
	pageSize int
	total    int
}

func (p *Pagination) Page() int {
	return p.page
}

func (p *Pagination) PageSize() int {
	return p.pageSize
}

func (p *Pagination) SetTotal(total int) {
	p.total = total
}

func (p *Pagination) Offset() int {
	return (p.page - 1) * p.pageSize
}

func (p *Pagination) MarshalJson() ([]byte, error) {
	return json.Marshal(
		struct {
			Page     int `json:"page"`
			PageSize int `json:"page_size"`
			Total    int `json:"total"`
		}{
			Page:     p.page,
			PageSize: p.pageSize,
			Total:    p.total,
		},
	)
}

func (p *Pagination) UnmarshalJson(data []byte) error {
	var temp struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
		Total    int `json:"total"`
	}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	p.page = temp.Page
	p.pageSize = temp.PageSize
	p.total = temp.Total
	return nil
}
