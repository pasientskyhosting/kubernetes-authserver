#!/bin/bash
openssl genrsa -out authserver.key 2048
openssl req -key authserver.key -new -out authserver.req -subj "/CN=kube-authserver" -config openssl.cnf
openssl x509 -req -in authserver.req -CA ca.pem  -CAkey ca-key.pem -CAserial ca.srl -out authserver.pem -days 3650 -extensions v3_req -extfile openssl.cnf
