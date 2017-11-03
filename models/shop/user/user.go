package user

type serviceProvider struct {}

var (
	Service *serviceProvider
)

type User struct {
	Id       string
	Username string
	Password string
	Gender  bool
	Age     int
	Address string
	Avatar  string
	Email   string
}

