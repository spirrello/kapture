FROM scratch
LABEL maintainer="Stefano Pirrello <spirrello@liaison.com>"

COPY build/out/kcapture-amd64 /
ENTRYPOINT ["/kcapture-amd64-amd64"]

