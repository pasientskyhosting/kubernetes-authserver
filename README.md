# Kubernetes webhook token authenticationserver
Simple lightweigh & database backed authserver written in GOLANG to be used with webhook token authentication mode in Kubernetes

More info about webhook token authentication [here](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)

## Preparations

*Generate SSL Certificates

*Create TLS secret
```
kubectl create secret tls authserver --cert=authserver.pem --key=authserver.key --namespace kube-system
```


## The following environment variables are used at startup of the docker container.
#### _DB_HOST_
Mysql hostname  
_Default: 127.0.0.1_

#### _DB_PORT_
Mysql port  
_Default: 3306_

#### _DB_NAME_
Mysql DB name  
_Default: auth_

#### _DB_USER_
Mysql username  
_Default: auth_

#### _DB_PASS_
Mysql password  
_Default: auth_

## Database preparation
The sql/db-layout.sql contains the structure needed for authserver. It will create a DB named auth uppon importing.

## JSON Requests & responses
#### Unsuccessfull request response
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": false
  }
}
```

#### Successfull response example
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": true,
    "user": {
      "username": "janedoe@example.com",
      "uid": "42",
      "groups": [
        "developers",
        "qa"
      ],
    }
  }
}
```

#### Faulty request ( Json check failed )
```
{
  "status": "400",
  "details": "Invalid TokenReview ( Json decode failed )"
}
```

## Command line options for kubernetes-authserver
### --host <string>
DB hostname / ip, default 127.0.0.1

### --port <int>
DB host port, default 3306

### --db <srtring>
DB databasename, default 'auth'

### --user <string>
DB username, default 'auth'

### --pass <string>
DB password, default 'auth'

### --charset <string>
DB charset, defaut 'utf8'

### --https <bool>
Enable HTTPS access, default true

### --http <bool>
Enable HTTP access, default true

### --cert <string>
Path to TLS cert, default /etc/ssl/tls.crt

### --key <string>
Path to TLS private key, default /etc/ssl/tls.key


## utilities/tokengen.go
This is a small utility to generate auth tokens for use with the system.
It's very basic at the moment a example to run it is:

go run tokengen.go --host=192.168.2.62 --db=auth --user=auth --pass=auth --username=jk

More work is needed on this, alternativly a admin UI will be introduced.