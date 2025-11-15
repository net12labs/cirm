package service

type ServiceUnit struct {
	Mode   RunMode
	OnExit func()
	// Other fields here
}

type Service struct {
	ServiceUnit
	// Other fields here
}

type ServiceSubUnit struct {
	ServiceUnit
	// Other fields here
}

func NewService() *Service {
	svc := &Service{}
	svc.ServiceUnit = ServiceUnit{
		Mode: RunMode{
			items: make(map[string]string),
		},
	}
	return svc
}
