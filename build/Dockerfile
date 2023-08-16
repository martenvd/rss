FROM alpine:latest

ADD rss-linux /rss-linux

RUN apk add --no-cache tzdata

EXPOSE 8082

CMD [ "/rss-linux" ]
