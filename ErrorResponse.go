package appstoreconnect

type ErrorResponse struct {
	Errors []struct {
		Id     string `json:"id"`
		Status int    `json:"status"`
		Code   string `json:"code"`
		Title  string `json:"title"`
		Detail string `json:"detail"`
	} `json:"errors"`
}
