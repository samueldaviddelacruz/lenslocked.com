package email

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/mailgun/mailgun-go.v3"
)

const (
	welcomeSubject = "Welcome to Lenslocked.net"
)
const welcomeText = `
Hi there!
Welcome to Lenslocked.net! we really hope you enjoy using
our application!

Best,
Samy

`

const welcomeHTML = `
Hi there!<br/>
Welcome to Lenslocked.net! we really hope you enjoy using
our application!
<br/>
Best,<br/>
Samy

`

type ClientConfig func(*Client)

func WithMailgun(domain, apiKey, publicKey string) ClientConfig {
	return func(c *Client) {
		mg := mailgun.NewMailgun(domain, apiKey)
		c.mg = mg
	}
}

func WithSender(name, email string) ClientConfig {
	return func(c *Client) {
		c.from = buildEmail(name, email)
	}
}

func NewClient(opts ...ClientConfig) *Client {
	client := Client{

		from: "support@lenslocked.net",
	}
	for _, opt := range opts {
		opt(&client)
	}
	return &client
}

type Client struct {
	from string
	mg   mailgun.Mailgun
}

func (c *Client) Welcome(toName, toEmail string) error {
	message := c.mg.NewMessage(c.from, welcomeSubject, welcomeText, buildEmail(toName, toEmail))
	message.SetHtml(welcomeHTML)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := c.mg.Send(ctx, message)

	return err
}

func buildEmail(name, email string) string {
	if name == "" {
		return email
	}
	return fmt.Sprintf("%s <%s>", name, email)
}
