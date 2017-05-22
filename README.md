# Kubernetes webhook token authenticationserver
Simple lightweigh & database backed authserver written in GOLANG to be used with webhook token authentication mode in Kubernetes

More info about webhook token authentication [here](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)

**This project is under heavy development and is not production ready**

## Preparations

### Generate SSL Certificates
Provided in the SSL folder is some scripts to help you generate the certificates needed.

Copy your ca.pem from your existing cluster into the SSL folder and run authserver.sh & authserver_client.sh to generate certificates valid for 10 years.

If you wish to have shorter validity. feel free to modify the scripts.

Also if you plan to use a different service ip for the authserver update *openssl.cnf*

### Create TLS secret

A TLS secret is needed with the server certificates as kubernetes-authserver reads the certs from a volumeMount.

```
kubectl create secret tls authserver --cert=authserver.pem --key=authserver.key --namespace kube-system
```

### Creating users and groups
At the moment this has to be done manually then the tokengen utility is ran against the database to generate a token.

The tokens are saved in scrypt encrypted format in the database, meaning if you loose your token it's impossible to recover it and a new one has to be generated.

#### Creating the user
Start with inserting a user with empty token into the users table.

#### Creating groups
Create groups in the groups table

#### Mapping users to groups
The groups_mapping table contains the user to groups mappings.
Insert rows mapping userid to groupid

#### Assigning a token to a user
See the section regarding tokengen further down.
A admin UI to easier administrate this will probably be introduced in the future.

## The following environment variables are used at startup of the docker container.
### DB_HOST <string>
Mysql hostname  
**Default: 127.0.0.1**

### DB_PORT <int>
Mysql port  
**Default: 3306**

### DB_NAME <string>
Mysql DB name  
**Default: auth**

### DB_USER <string>
Mysql username  
**Default: auth**

### DB_PASS <string>
Mysql password  
**Default: auth**

### DB_CHARSET <bool>
Charset to use for database  
**Default: utf8**

## Database preparation
The sql/db-layout.sql contains the structure needed for authserver. It will create a DB named auth uppon importing.

## JSON Requests & responses
### Unsuccessfull request response
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": false
  }
}
```

### Successfull response example
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

### Faulty request ( Json check failed )
```
{
  "status": "400",
  "details": "Invalid TokenReview ( Json decode failed )"
}
```

## Command line options for kubernetes-authserver
### --host <string>
DB hostname / ip  
**default: 127.0.0.1**

### --port <int>
DB host port  
**default: 3306**

### --db <srtring>
DB databasename  
**default auth**

### --user <string>
DB username  
**default: auth**

### --pass <string>
DB password  
**default: auth**

### --charset <string>
DB charset  
**default: utf8**

### --https <bool>
Enable HTTPS access  
**default: true**

### --http <bool>
Enable HTTP access  
**default: true**

### --cert <string>
Path to TLS cert  
**default: /etc/ssl/tls.crt**

### --key <string>
Path to TLS private key  
**default: /etc/ssl/tls.key**

# utilities/tokengen.go
This is a small utility to generate auth tokens for use with the system.
It's very basic at the moment a example to run it is:

go run tokengen.go --host=192.168.2.62 --db=auth --user=auth --pass=auth --username=jk

Output will look like

```
2017/03/16 12:30:17 Username: jk
2017/03/16 12:30:17 Token: v864d329d5c8b9aw$ff1b4e4107fnd728b8169c3d89kdoebbb81933f32b09f4216211934895acea77
```
the token is copied into your kubectl config usually located in ~/.kube/config

*Example:*
```
apiVersion: v1
clusters:
- cluster:
    certificate-authority: ca.pem
    server: https://10.10.10.10
  name: default-cluster
contexts:
- context:
    cluster: default-cluster
    user: default
  name: default-system
current-context: default-system
kind: Config
preferences: {}
users:
- name: default
  user:
    token: "v864d329d5c8b9aw$ff1b4e4107fnd728b8169c3d89kdoebbb81933f32b09f4216211934895acea77"
```

More work is needed on this, alternativly a admin UI will be introduced.