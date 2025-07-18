package command

type SendEmailCommand struct {
	FromEmail string
	ToEmails  []string
	Subject   string
	HtmlBody  string
}
