package clients

type HttpClient interface {
	Ping() error
}
