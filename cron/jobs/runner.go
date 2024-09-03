package jobs

type Runner interface {
	Run()
	GetPattern() string
}
