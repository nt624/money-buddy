package auth

import (
"context"
"log"

firebase "firebase.google.com/go/v4"
"firebase.google.com/go/v4/auth"
"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() error {
	ctx := context.Background()

	// サービスアカウントキーのパス
	opt := option.WithCredentialsFile("firebase-admin-key.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	FirebaseAuth, err = app.Auth(ctx)
	if err != nil {
		return err
	}

	log.Println("Firebase Admin initialized successfully")
	return nil
}
