package api

type DataApi struct {
	Code int
	Data any
}

func NewDataApi(code int, data any) *DataApi {
	return &DataApi{Code: code, Data: data}
}
