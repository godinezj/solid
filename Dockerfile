# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM gobuffalo/buffalo:v0.14.2 as builder

RUN mkdir -p $GOPATH/src/bitbucket.org/godinezj/solid
WORKDIR $GOPATH/src/bitbucket.org/godinezj/solid

ADD . .
RUN go get $(go list ./... | grep -v /vendor/)
RUN buffalo build --static -o /bin/app

FROM ubuntu:latest

# Install dependencies
RUN apt-get update && apt-get install -y curl
RUN curl -sO http://archive.ubuntu.com/ubuntu/pool/universe/e/easy-rsa/easy-rsa_3.0.4-2_all.deb && \
    dpkg -i easy-rsa_3.0.4-2_all.deb && rm easy-rsa_3.0.4-2_all.deb && \
    apt-get clean && \
    ln -s /usr/share/easy-rsa/easyrsa /usr/local/bin

WORKDIR /bin/

COPY --from=builder /bin/app .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0

EXPOSE 3000

# Uncomment to run the migrations before running the binary:
# CMD /bin/app migrate; /bin/app
CMD exec /bin/app
