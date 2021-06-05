FROM alpine

WORKDIR /srv/pink-lady

ADD ./dist/pink-lady.tar.gz /srv/

EXPOSE 4869 4870
ENTRYPOINT ["./apiserver", "-c", "./config.default.toml"]
