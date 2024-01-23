package models

type HostMetric struct { //Формируем структуру ответа в формате json
	CPULoad      float64 `json:"cpu_load"`
	ThreadsCount int     `json:"threads_count"`
}

type Post struct {
	ID       int   `json:"id,omitempty"`
	Elements []int `json:"elements,omitempty"` // omitempty - если поле Elements равно nil или пустому слайсу,
}
