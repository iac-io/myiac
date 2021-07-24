FROM google/cloud-sdk:alpine

ARG TERRAFORM_VERSION=1.0.3
ARG HELM_VERSION=3.5.2
ARG KUBE_PROMPT_VERSION=1.0.11

ENV GOPATH="/root/go"

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
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* && \
    wget https://github.com/c-bata/kube-prompt/releases/download/v${KUBE_PROMPT_VERSION}/kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    unzip kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    rm -f kube-prompt_v${KUBE_PROMPT_VERSION}_linux_amd64.zip && \
    chmod +x kube-prompt && \
    mv ./kube-prompt /usr/local/bin/kube-prompt && \
    sh -c "$(wget -O- https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k && \
    /usr/bin/helm plugin install https://github.com/databus23/helm-diff --version 3.1.3

ADD container-files /
