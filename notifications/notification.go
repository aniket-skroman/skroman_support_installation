package notifications

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

type Notification struct {
	MsgTitle          string
	MsgBody           string
	RegistrationToken string
}

func (n *Notification) SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {

	ctx := context.Background()

	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		panic("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic("Firebase load error")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)

	return app, ctx, client
}

func (n *Notification) SendToToken(app *firebase.App) {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	for i := 0; i < 5; i++ {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title:    n.MsgTitle,
				Body:     n.MsgBody,
				ImageURL: "http://15.207.19.172:9000/api/device-file/media/upload-3760796307.png",
			},
			Token: n.RegistrationToken,
		}

		response, err := client.Send(ctx, message)
		if err != nil {
			log.Fatalln("Error throwing..", err)
		}
		fmt.Println("Successfully sent message:	", response)
	}
	// fmt.Println("Successfully sent message:", response)
}
