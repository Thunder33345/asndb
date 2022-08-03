package asndb

type Option func(r *ASList)

func WithAssumeValid() Option {
	return func(r *ASList) {
		r.assumeValid = true
	}
}
