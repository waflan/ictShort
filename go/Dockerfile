FROM golang:alpine
WORKDIR /app
EXPOSE 80
# EXPOSE 443

ADD ./wait /wait 

RUN chmod +x /wait&&\
    apk --update add tzdata&&\
    cp  /usr/share/zoneinfo/Asia/Tokyo /etc/localtime&&\
    apk del tzdata&&\
    rm -rf /var/cache/apk/*

ADD ./app /app

# RUN apk add -U git
# RUN go get github.com/go-sql-driver/mysql
# RUN go get github.com/dgrijalva/jwt-go

CMD ["go","run","main.go"]