FROM golang:alpine3.13

MAINTAINER dsMoon <deuksoo.mun@catenoid.net>

ENV APP_UID 2001
ENV APP_USER kollus
ENV APP_GROUP kollus

#redis 접속 관련 추가
#RUN apk add redis

# PID 2001로 kollus User 생성
RUN adduser $APP_USER -h /home/$APP_USER -u $APP_UID -D

# kollus에서 작업
WORKDIR /home/kollus

RUN mkdir /var/log/kollus && chmod -R 777 /var/log/kollus

COPY ./main /home/kollus/main

CMD ./main

EXPOSE 8888

