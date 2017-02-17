#Kubernetes authserver by Joakim Karlsson <jk@patientsky.com>
Simple lightweigh & database backed authserver written in GOLANG to be used with webhook token authentication mode in Kubernetes

More info about webhook token authentication [here](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)

##The following environment variables are used at startup
###__DB_HOST__
Mysql hostname
_Default: 127.0.0.1_

###__DB_PORT__
Mysql port
_Default: 3306_

###__DB_NAME__
Mysql DB name
_Default: auth_

###__DB_USER__
Mysql username
_Default: auth_

###__DB_PASS__
Mysql password
_Default: auth_



##Database preparation

```
CREATE DATABASE auth CHARACTER SET utf8 COLLATE utf8_general_ci;

USE auth;

CREATE TABLE `users` (
  `id` int(6) NOT NULL,
  `token` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `uid` int(6) NOT NULL,
  `groups` varchar(255) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `uid` (`uid`),
  UNIQUE KEY `token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

##JSON Requests & responses

###Unsuccessfull request response
```
{
  "apiVersion": "authentication.k8s.io/v1beta1",
  "kind": "TokenReview",
  "status": {
    "authenticated": false
  }
}
```

###Successfull response example
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

###Faulty request ( Json check failed )
```
{
  "status": "400",
  "details": "Invalid TokenReview ( Json decode failed )"
}
```