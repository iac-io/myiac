FROM google/cloud-sdk:alpine

ENV TERRAFORM_VERSION=0.12.29 \
    HELM_VERSION=3.1.2 \
    GOPATH="/home/app/go"

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
        go && \
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
    chown -R app:app /home/app && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

USER app

RUN \
    sh -c "$(wget -O- https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k && \
    mkdir -p /home/app/.cache && \
    mkdir -p /home/app/go

ADD container-files /

USER root

RUN \
    chown -R app:app /home/app

USER app
