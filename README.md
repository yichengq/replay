# Reason Playground

The root directory holds the source of reason playground frontend.

The /sandbox directory holds the source of sandbox to compile and run reason programs.

# Local Development

## Building

```
# build the sandbox image
$ docker build -t sandbox sandbox/
```

## Running

```
# run reason playground frontend at the root directory
$ goapp serve -host 0.0.0.0 -port 8081 $PWD
# visit page at http://localhost:8081
```

```
# run the sandbox
$ docker run -d -p 8080:8080 sandbox
# get docker host ip, try boot2docker fallback on localhost.
$ DOCKER_HOST_IP=$(boot2docker ip || echo localhost)
# run some reason code
$ cat /path/to/code.re | curl --data @- $DOCKER_HOST_IP:8080/compile\?type=to_run
```

# Deployment

(available to admin only) To deploy the front-end, use `play/deploy.sh`.

```
$ gcloud --project replay-154206 app deploy app.yaml
$ gcloud --project replay-154206 app deploy sandbox/app.yaml
```

Use the Cloud Console's to check status:
	https://console.cloud.google.com/appengine/versions?project=replay-154206
