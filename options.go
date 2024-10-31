package auth

type Option func(*options)

type options struct {
	admin bool
}

func OnlyAdmin() Option {
	return func(o *options) {
		o.admin = true
	}
}
