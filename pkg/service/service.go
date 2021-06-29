package service

type Service interface {
	Stop()
}
type stopFunc func()

func (fn stopFunc) Stop() {
	fn()
}
