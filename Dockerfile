FROM scratch
LABEL maintainer="Stefano Pirrello <spirrello@opentext.com>"

COPY kcapture /
ENTRYPOINT ["kcapture"]
