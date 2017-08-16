FROM scratch

MAINTAINER estafette.io

COPY ca-certificates.crt /etc/ssl/certs/
COPY estafette-extension-github-status /

ENTRYPOINT ["/estafette-extension-github-status"]