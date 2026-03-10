#!/bin/bash
set -e

# Install Docker
yum update -y
yum install -y docker
systemctl enable docker
systemctl start docker

docker volume create ssh_keys

docker run -d \
  --name ssh-portfolio \
  --restart always \
  -p 22:2222 \
  -v ssh_keys:/app/.ssh \
  ${dockerhub_image}
