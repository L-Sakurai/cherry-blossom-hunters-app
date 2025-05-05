FROM golang:1.20

RUN apt-get update && \
    apt-get install -y git curl unzip && \
    rm -rf /var/lib/apt/lists/*

RUN useradd -m cherry-blossom-hunters-app

# Install Fly.io (do this before switching the user)
USER cherry-blossom-hunters-app
ENV FLYCTL_INSTALL=/home/cherry-blossom-hunters-app/.fly
ENV PATH="${FLYCTL_INSTALL}/bin:${PATH}"
RUN curl -L https://fly.io/install.sh | sh