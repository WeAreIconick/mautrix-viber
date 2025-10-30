package matrix

type Bridge interface {
	Start() error
	Stop() error
}

