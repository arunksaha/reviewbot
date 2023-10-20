FROM ubuntu:22.04
WORKDIR /app
COPY reviewbot .
RUN apt update
RUN apt install -y ca-certificates
CMD ["./reviewbot"]
