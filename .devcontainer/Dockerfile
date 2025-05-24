FROM golang:1.20

# Install necessary system packages using apt-get
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        git \
        curl \
        unzip \
        python3 \
        python3-venv \
        python3-pip \
        chromium \
        chromium-driver \
        jq \
        tree \
        bash \
        passwd && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create Python virtual environment and install required packages
ENV VENV_PATH=/opt/venv
RUN python3 -m venv $VENV_PATH && \
    $VENV_PATH/bin/pip install --no-cache-dir \
        beautifulsoup4 \
        selenium

# Set environment variables for using the virtual environment and Chromium
ENV PATH="$VENV_PATH/bin:$PATH"
ENV CHROME_BIN="/usr/bin/chromium-browser"

# Create application user
RUN useradd -ms /bin/bash cherry-blossom-hunters-app

# Install Fly.io CLI
USER cherry-blossom-hunters-app
ENV FLYCTL_INSTALL=/home/cherry-blossom-hunters-app/.fly
ENV PATH="${FLYCTL_INSTALL}/bin:${PATH}"
RUN curl -L https://fly.io/install.sh | sh
