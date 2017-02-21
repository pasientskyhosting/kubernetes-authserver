#!/bin/bash
openssl genrsa -out authclient.key 2048
openssl req -key authclient.key -new -out authclient.req -subj "/CN=kube-authserver" -config openssl.cnf
openssl x509 -req -in authclient.req -CA ca.pem  -CAkey ca-key.pem -CAserial ca.srl -out authclient.pem -days 3650 -extensions v3_req -extfile openssl.cnf
