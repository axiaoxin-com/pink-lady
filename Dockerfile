FROM alpine

WORKDIR /srv/pink-lady

ADD ./dist/pink-lady.tar.gz /srv/

EXPOSE 4869 4870
ENTRYPOINT ["./apiserver", "-p", ".", "-c", "config.default", "-t", "toml"]
