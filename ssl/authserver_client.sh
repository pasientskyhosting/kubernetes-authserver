#!/bin/bash
openssl genrsa -out authclient.key 2048
openssl req -key authclient.key -new -out authclient.req
openssl x509 -req -in authclient.req -CA ../ClusterCA/ca.pem  -CAkey ../ClusterCA/ca-key.pem -CAserial ../ClusterCA/ca.srl -out authclient.pem 
