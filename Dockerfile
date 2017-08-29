FROM scratch

LABEL maintainer="estafette.io" \
      description="The estafette-extension-github-status component is an Estafette extension to update build status in Github for builds handled by Estafette CI"

COPY ca-certificates.crt /etc/ssl/certs/
COPY estafette-extension-github-status /

ENTRYPOINT ["/estafette-extension-github-status"]