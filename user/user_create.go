package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
)

// GOOGLE_CLOUD_PROJECT is automatically set by the Cloud Functions runtime.
var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

// client is a Firestore client, reused between function invocations.
var client *firestore.Client

func init() {
	// Use the application default credentials.
	conf := &firebase.Config{ProjectID: projectID}

	// Use context.Background() because the app/client should persist across
	// invocations.
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalf("firebase.NewApp: %v", err)
	}

	client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}
}

// CreateUser is creation user info
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// 最終的に firestore の client を明示的に Close する
	defer client.Close()

	user := &User{}
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		// リクエストの json がでコードできなかった時
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Decord Error.\n")
		log.Printf("JSON Decord Error: %+v", err)
		return
	}
	if user.UUID == "" {
		// json 内に name プロパティがなかった時
		// UUID の生成
		u, err := uuid.NewRandom()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Generate UUID Error.\n")
			log.Printf("Generate UUID Error: %+v", err)
			return
		}
		user.UUID = u.String()
	}

	// firestore に既に登録されているか

	collection := "Users"
	c := r.Context()
	if _, err := client.Collection(collection).Doc(user.UUID).Get(c); err == nil {
		// エラーじゃないときは多分いる
		// TODO: code = NotFound desc で判定
		w.WriteHeader(http.StatusForbidden)
		io.WriteString(w, fmt.Sprintf("Exist User UUID: %s", user.UUID))
		log.Printf("Exist User UUID: %s", user.UUID)
		return
	}

	if _, err := client.Collection(collection).Doc(user.UUID).Set(c, user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Generate User Error.\n")
		log.Printf("Get User Document Error: %+v", err)
		return
	}

	log.Printf("Crented User UUID: %s", user.UUID)

	res, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "JSON Marshal Error.\n")
		log.Printf("JSON Marshal Error: %+v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

	return
}
