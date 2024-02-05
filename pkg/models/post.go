package models

type HostMetric struct { //Формируем структуру ответа в формате json
	CPULoad      float64 `json:"cpuLoad"` // только camelCase
	ThreadsCount int     `json:"threadsCount"`
}

type Post struct {
	ID       int   `json:"id,-"`               // не прочитает через Marshall, даже если будет в теле
	Elements []int `json:"elements,omitempty"` // omitempty - если поле Elements равно nil или пустому слайсу,
}
