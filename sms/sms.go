package sms

import "log"

// ISMS is a interface
type ISMS interface {
	Send(phone string, message string) error
}

// SMSSvc is a struct
type SMSSvc struct {
}

// Send is a func
func (smsSvc *SMSSvc) Send(phone string, message string) error {
	// Just a fake function
	log.Println("phone:", phone, "message:", message)
	return nil
}
