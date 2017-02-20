#!/bin/bash
openssl genrsa -out authserver.key 2048
openssl req -key authserver.key -new -out authserver.req 
openssl x509 -req -in authserver.req -CA ../ClusterCA/ca.pem  -CAkey ../ClusterCA/ca-key.pem -CAserial ../ClusterCA/ca.srl -out authserver.pem
