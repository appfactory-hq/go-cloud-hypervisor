FROM golang:latest

RUN apt-get update && apt-get install -y \
    build-essential \
    libssl-dev \
    make \
    git \
    curl \
    wget \
    cpu-checker \
    qemu-utils \
    libcap2-bin \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN wget https://github.com/cloud-hypervisor/cloud-hypervisor/releases/download/v32.0/cloud-hypervisor-static
RUN chmod +x cloud-hypervisor-static && mv cloud-hypervisor-static /usr/bin/cloud-hypervisor

RUN mkdir -p /data

RUN wget https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img && \
    qemu-img convert -p -f qcow2 -O raw jammy-server-cloudimg-amd64.img jammy-server-cloudimg-amd64.raw && \
    ls -la && \
    rm jammy-server-cloudimg-amd64.img && \
    mv jammy-server-cloudimg-amd64.raw /data/jammy-server-cloudimg-amd64.raw

#RUN wget https://cloud-images.ubuntu.com/focal/current/focal-server-cloudimg-amd64.img && \
#    qemu-img convert -p -f qcow2 -O raw focal-server-cloudimg-amd64.img focal-server-cloudimg-amd64.raw && \
#    ls -la && \
#    rm focal-server-cloudimg-amd64.img && \
#    mv focal-server-cloudimg-amd64.raw /data/focal-server-cloudimg-amd64.raw

RUN wget https://github.com/cloud-hypervisor/rust-hypervisor-firmware/releases/download/0.4.2/hypervisor-fw && \
    mv hypervisor-fw /data/hypervisor-fw

RUN setcap cap_net_admin+ep /usr/bin/cloud-hypervisor

CMD [ "/src/.build/run.sh" ]

ADD . /src
