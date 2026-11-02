# Test image for itamae integration runs.
#
# It mirrors a freshly installed workstation: a regular user with passwordless
# sudo and only the packages a desktop install ships with (sudo, curl, CA certs).
# Everything else must be installed by itamae itself.
ARG BASE_IMAGE=ubuntu:24.04
FROM ${BASE_IMAGE}

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update \
    && apt-get install -y --no-install-recommends sudo ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

RUN useradd --create-home --shell /bin/bash tester \
    && echo 'tester ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/tester \
    && chmod 0440 /etc/sudoers.d/tester

USER tester
WORKDIR /home/tester
