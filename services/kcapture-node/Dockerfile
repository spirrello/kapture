FROM alpine:latest
LABEL maintainer "Stefano Pirrello <spirrello@gmail.com>"
COPY Dockerfile /Dockerfile
COPY .bashrc /root/.bashrc

RUN apk --no-cache --update --verbose add grep bash tcpdump && \
    rm -rf /var/cache/apk/* /tmp/* /sbin/halt /sbin/poweroff /sbin/reboot

ENTRYPOINT ["sh"]