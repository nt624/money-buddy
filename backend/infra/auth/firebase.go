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

	// 環境変数から認証情報を取得
	credentialsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	credentialsJSON := os.Getenv("FIREBASE_CREDENTIALS_JSON")

	var opt option.ClientOption

	// JSON文字列が設定されている場合は優先して使用
	if credentialsJSON != "" {
		opt = option.WithCredentialsJSON([]byte(credentialsJSON))
		log.Println("Using Firebase credentials from FIREBASE_CREDENTIALS_JSON")
	} else if credentialsPath != "" {
		opt = option.WithCredentialsFile(credentialsPath)
		log.Printf("Using Firebase credentials from file: %s", credentialsPath)
	} else {
		// デフォルトのファイルパス
		opt = option.WithCredentialsFile("firebase-admin-key.json")
		log.Println("Using default Firebase credentials file: firebase-admin-key.json")
	}

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
