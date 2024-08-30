package iternal

type EmailHeader struct {
	Sender   string
	Receiver string
	Data     int64
}

type EmailSender interface {
	SendEmail(email string, header EmailHeader, body string) error
}
