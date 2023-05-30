package action

type Request interface {
	Run() (string, error)
}
