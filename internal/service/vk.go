package service

type VK struct {
	logger LogPusher
}

type LogPusher interface {
	Info(args ...interface{})
	Error(args ...interface{})
	Warn(args ...interface{})
}

func NewVK(logger LogPusher) *VK {
	return &VK{
		logger: logger,
	}
}
