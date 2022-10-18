# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.238.1/containers/go/.devcontainer/base.Dockerfile

ARG VARIANT="1.19-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# Copy in cross-platform scripts. These scripts may install different binaries
# depending on the architecture of the local users computer (like if they are on an M1 or a junk Macbook)
COPY ./container_init /root/init_scripts
RUN /root/init_scripts/install_aws.sh
