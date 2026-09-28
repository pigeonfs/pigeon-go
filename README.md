# Pigeon Go SDK

Official Go client for the Pigeon email API. Package layout follows [resend-go](https://github.com/resend/resend-go): `client.Emails.Send`, `client.Domains`, `client.Contacts`.

## Install

```bash
go get github.com/pigeonfs/pigeon-go
```

## Setup

```go
package main

import (
	"context"
	"fmt"
	"log"

	pigeon "github.com/pigeonfs/pigeon-go"
)

func main() {
	client := pigeon.NewClient("pg_xxxx")
	email, err := client.Emails.Send(context.Background(), &pigeon.SendEmailRequest{
		From:    "Ada <ada@yourdomain.com>",
		To:      []string{"person@example.com"},
		Subject: "hello world",
		Html:    "<strong>it works!</strong>",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Email %s has been sent\n", email.Id)
}
```

Set `PIGEON_BASE_URL` (default `http://localhost:4005/`) when the API is not on localhost.

Verify the sending domain in the Pigeon dashboard before sending.

## License

MIT
