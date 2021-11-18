FROM alpine:3.8

WORKDIR /home/
COPY cr-adapter-service.bin .

RUN chmod +x cr-adapter-service.bin

EXPOSE 8081

CMD [ "./cr-adapter-service.bin" ]