# Step
1. Generate user Folder
2. Generate `go.mod`
    1. `$ cd {project folder}`
    2. `$ go mod init`
3. Generate `user/model.go` and `user/user_create.go`
4. Generate `go.sum`
    1. `$ cd {project folder}/user`
    2. `$ go mod download`
5. Generate `cloudbuild.yml`

## Activate GCP Services and APIs
- [Cloud Build](https://console.developers.google.com/apis/api/cloudbuild.googleapis.com/overview)
- [Cloud Functions](https://console.developers.google.com/apis/api/cloudfunctions.googleapis.com/overview)
- [Cloud Resource Manager](https://console.developers.google.com/apis/api/cloudresourcemanager.googleapis.com/overview)
- [Firestore](https://console.cloud.google.com/firestore)

## deploy
1. `$ cd {project folder}`
2. `$ gcloud builds submit . --config cloudbuild.yml`

# reference
- [base](https://github.com/r-konishi/cloud-functions-sample)
- https://godoc.org/cloud.google.com/go/firestore
- https://github.com/GoogleCloudPlatform/golang-samples/blob/master/functions/firebase/upper/upper.go
