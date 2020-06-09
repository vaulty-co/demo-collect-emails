# demo-collect-emails

## Prerequisites

* Docker and Docker compose
* API Key from active Mailgun account

Create .env file with Mailgun API Key and Vaulty configuration:

```
MG_API_KEY=xxxxxxxxxxx
PROXY_PASS=12345
ENCRYPTION_KEY=776f726420746f206120736563726574
```

## Run without Vaulty

You can build and run the demo app with docker compose:

```docker-compose up```

## Run with Vaulty

To run demo app and Vaulty instance:

```
docker-compose --file docker-compose-vaulty.yml up
```

Get into backend container:

```
docker-compose --file docker-compose-vaulty.yml exec backend sh
```

and then inside backend container run following to add Vaulty's CA into system certificates:

```
cp ./vaulty/ca.cert /usr/local/share/ca-certificates/vaulty.cert && update-ca-certificates
```

then navigate your browser to http://127.0.0.1:3001.
