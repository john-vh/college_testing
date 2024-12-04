package services

type BackgroundService interface {
	Start()
	Stop()
}
