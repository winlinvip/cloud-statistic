FROM ossrs/srs:ubuntu20 as build

ADD . /g
WORKDIR /g
RUN go build -mod vendor .
RUN mkdir -p  /usr/local/cloud-statistic
RUN cp cloud-statistic /usr/local/cloud-statistic/app
RUN cp -r static /usr/local/cloud-statistic/static
RUN cp index.tmpl.html /usr/local/cloud-statistic/index.tmpl.html

FROM ubuntu:focal as dist

COPY --from=build /usr/local/cloud-statistic /usr/local/cloud-statistic

WORKDIR /usr/local/cloud-statistic
CMD ["./app"]
