package HttpService

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}
