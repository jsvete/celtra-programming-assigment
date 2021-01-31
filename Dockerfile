FROM golang:1.15.7-alpine3.13

RUN mkdir /code
WORKDIR /code

VOLUME [ "/code" ]

CMD ["/bin/bash"]