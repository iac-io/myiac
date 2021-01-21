FROM golang:1.15-alpine as myiac-builder

# Dummy
ARG EXTRA_WORKDIR_COPY=/workdir/Dockerfiles
ARG EXTRA_WORKDIR_DEST=/home/app/Dockerfiles
ADD . /workdir

RUN \
    cd /workdir && \
    go build -o /usr/bin/myiac cmd/myiac/myiac.go

FROM google/cloud-sdk:alpine

COPY --from=myiac-builder /usr/bin/myiac /usr/bin/myiac
COPY --from=myiac-builder /workdir/*.json /home/app/
COPY --from=myiac-builder $EXTRA_WORKDIR_COPY $EXTRA_WORKDIR_DEST

ARG CURRENT_HELM_VERSION=3.1.2
ENV TERRAFORM_VERSION=0.12.29 \
    HELM_VERSION=$CURRENT_HELM_VERSION

RUN \
    apk --update add \
        openjdk7-jre \
        curl \
        jq \
        bash \
        ca-certificates \
        git \
        openssl \
        unzip \
        wget \
        util-linux \
        vim \
        zsh \
        go \
        busybox-extras \
        mtr \
        sed && \
    gcloud components install kubectl && \
    cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    wget https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    tar -zxvf helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/bin && \
    rm -f helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    addgroup -S app && adduser -S app -G app -h /home/app && \
    chown app:app /usr/bin/myiac && \
    chown app:app -R /home/app && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

USER app
CMD ['/usr/bin/myiac']