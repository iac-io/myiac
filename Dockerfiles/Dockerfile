FROM golang:1.15-alpine as myiac-builder

# Placeholders to include extra folders in final container (i.e. charts)
ARG EXTRA_WORKDIR_ORIG=/workdir/Dockerfiles
ARG EXTRA_WORKDIR_DEST=/home/app/Dockerfiles

ADD . /workdir

RUN \
    cd /workdir && \
    go build -o /usr/bin/myiac cmd/myiac/myiac.go

FROM google/cloud-sdk:alpine

COPY --from=myiac-builder /usr/bin/myiac /usr/bin/myiac
#COPY --from=myiac-builder /workdir/*.json /home/app/
COPY --from=myiac-builder /workdir/charts-dns /home/app/charts

ARG CURRENT_HELM_VERSION=3.1.2
ARG CURRENT_TERRAFORM_VERSION=0.12.29
ARG CURRENT_KUBE_PROMPT_VERSION=1.0.11
ENV TERRAFORM_VERSION=$CURRENT_TERRAFORM_VERSION \
    HELM_VERSION=$CURRENT_HELM_VERSION \
    KUBE_PROMPT_VERSION=$CURRENT_KUBE_PROMPT_VERSION

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
    mkdir -p /home/app/go && \
    /usr/bin/helm plugin install https://github.com/databus23/helm-diff --version 3.1.3

ADD Dockerfiles/container-files /

USER root

RUN \
    chown -R app:app /home/app

USER app

CMD ['/usr/bin/myiac']