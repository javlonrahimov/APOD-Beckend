package main

type Handlers struct {
	Users UserHandler
	HealthCheck HealthchekHandler
	Apods ApodHandler
}

func NewHandler(app *application) *Handlers {
	return &Handlers{
		Users: NewUserApi(app),
		HealthCheck: NewHealthcheckApi(app),
		Apods: NewApodApi(app),
	} 
}