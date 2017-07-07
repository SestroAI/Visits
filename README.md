# Sestro Visit Service

## Config

Please change the env variables in app.yaml based on the google/firebase project you want to deploy to

## Update vendor shared library

```
#Remove shared repo cache
$ sudo rm -r ~/.glide/cache/src/https-github.com-SestroAI-shared/
$ glide up
$ glide install -v
```

## Deploy

```
$ gcloud auth login
$ gcloud config set project <project-id:sestro-165123>
$ gcloud app deploy
```
