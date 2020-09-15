FROM golang:1.13-alpine as myiac-builder
RUN apk add --no-cache git bash sudo curl
WORKDIR /go/src/github.com/dfernandezm/myiac
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

FROM google/cloud-sdk:alpine
RUN apk --update add openjdk7-jre
RUN gcloud components install kubectl
ENV TERRAFORM_VERSION=0.12.17
RUN apk update && \
    apk add curl jq python bash ca-certificates git openssl unzip wget util-linux vim && \
    cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin
RUN wget https://get.helm.sh/helm-v2.16.1-linux-amd64.tar.gz && \
    tar -zxvf helm-v2.16.1-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/bin

# https://stackoverflow.com/questions/49955097/how-do-i-add-a-user-when-im-using-alpine-as-a-base-image

RUN addgroup -S app && adduser -S app -G app -h /home/app
USER app
RUN mkdir -p /home/app
WORKDIR /home/app
COPY --from=myiac-builder /go/bin/myiac .
COPY --chown=app:app account.json .
COPY --chown=app:app internal/helperScripts helperScripts
COPY --chown=app:app charts charts
COPY --chown=app:app entrypoint.sh /home/app
RUN chmod +x /home/app/entrypoint.sh
CMD /home/app/entrypoint.sh

