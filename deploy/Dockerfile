FROM alpine:3.8

WORKDIR /home/
COPY cr-product.bin .

RUN chmod +x cr-product.bin

EXPOSE 8084

CMD [ "./cr-product.bin" ]