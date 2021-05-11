FROM google/cloud-sdk:alpine

ARG UID=100
ARG GID=100
ARG TERRAFORM_VERSION=0.15.3
ARG HELM_VERSION=3.1.2
ARG KUBE_PROMPT_VERSION=1.0.10

ENV GOPATH="/home/app/go"

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
        mtr && \
    gcloud components install kubectl && \
    cd /tmp && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/bin && \
    rm -f terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    wget https://get.helm.sh/helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    tar -zxvf helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    mv linux-amd64/helm /usr/bin && \
    rm -f helm-v${HELM_VERSION}-linux-amd64.tar.gz && \
    addgroup -g $GID -S app && adduser -u $UID -S app -G app -h /home/app && \
    chown -R app:app /home/app && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* && \
    wget https://github.com/c-bata/kube-prompt/releases/download/v${KUBE_PROMPT_VERSION}/kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    unzip kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    rm -f kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    chmod +x kube-prompt && \
    mv ./kube-prompt /usr/local/bin/kube-prompt

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
