FROM alpine:3.5

LABEL maintainer "Joakim Karlsson <jk@patientsky.com>"

RUN apk add --no-cache su-exec
ADD kubernetes-authserver /kubernetes-authserver
ADD scripts/entrypoint.sh /entrypoint.sh

RUN chmod a+x /kubernetes-authserver && chmod a+x /entrypoint.sh
RUN mkdir -p /etc/ssl/

ENTRYPOINT ["/entrypoint.sh"]
CMD [""]
