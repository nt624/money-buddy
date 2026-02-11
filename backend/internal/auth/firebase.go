package auth

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() error {
	ctx := context.Background()

	// 環境変数からサービスアカウントキーのパスを取得
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credentialsPath == "" {
		credentialsPath = "firebase-admin-key.json"
	}
	opt := option.WithCredentialsFile(credentialsPath)

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
