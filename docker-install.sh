#!/usr/bin/env sh

PLATFORM="$(uname | awk '{print tolower($0)}')"
REPORTED_ARCH="$(uname -m)"

if [ "$REPORTED_ARCH" = "x86_64" ]; then
  DOWNLOAD_ARCH="amd64"
elif [ "$REPORTED_ARCH" = "i686" ]; then
  DOWNLOAD_ARCH="386"
elif [ "$REPORTED_ARCH" = "i386" ]; then
  DOWNLOAD_ARCH="386"
elif [ "$REPORTED_ARCH" = "arm" ]; then
  DOWNLOAD_ARCH="arm"
elif [ "$REPORTED_ARCH" = "armv7l" ]; then
  DOWNLOAD_ARCH="arm"
elif [ "$REPORTED_ARCH" = "aarch64" ]; then
  DOWNLOAD_ARCH="arm64"
elif [ "$REPORTED_ARCH" = "arm64" ]; then
  DOWNLOAD_ARCH="arm64"
else
  echo "Unknown Architecture"
  exit
fi

URL="https://github.com/omranjamal/mono-cd/releases/latest/download/mono-cd_${PLATFORM}_${DOWNLOAD_ARCH}"
INSTALL_PATH="$HOME/.local/share/omranjamal/mono-cd"

create_directory() {
  mkdir -p "$INSTALL_PATH" && return 0
}

download() {
  echo "> download: 📥 downloading mono-cd_${DOWNLOAD_ARCH}" && \
      wget -qO "$INSTALL_PATH/mono-cd" "$URL" && return 0
}

update_permissions() {
  echo "> permissions: 💪 setting execution permission" && chmod +x "$INSTALL_PATH/mono-cd" && return 0
}

add_to_shell() {
  touch "$HOME/.profile"
  echo "> shell function: ⚡ Adding to ~/.profile"
  $INSTALL_PATH/mono-cd --install "$HOME/.profile"
}

completed() {
  return 0
}

create_directory && \
  download && \
  update_permissions && \
  add_to_shell && \
  completed
