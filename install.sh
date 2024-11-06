#!/usr/bin/env sh

V="${VERSION:-vvvv}"
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
else
  echo "Unknown Architecture"
  exit
fi

URL="https://github.com/omranjamal/mono-cd/releases/latest/download/mono-cd_${V}_${DOWNLOAD_ARCH}"
INSTALL_PATH="$HOME/.local/share/omranjamal/mono-cd"

create_directory() {
  mkdir -p "$INSTALL_PATH" && return 0
}

download() {
  echo "> download: 📥 downloading mono-cd_${V}_${DOWNLOAD_ARCH}" && \
      curl -s -L -o "$INSTALL_PATH/mono-cd" "$URL" && return 0
}

update_permissions() {
  echo "> permissions: 💪 setting execution permission" && chmod +x "$INSTALL_PATH/mono-cd" && return 0
}

add_to_shell() {
  if [ -f "$HOME/.bashrc" ]; then
    echo "> detected: ~/.bashrc"
    echo "> shell function: ⚡ Adding to ~/.bashrc"
    $INSTALL_PATH/mono-cd --install "$HOME/.bashrc"

    echo   "\e[0m"
    echo   "    🚀 Run this, or re-start your bash terminal:"
    echo   "       $ \e[1msource ~/.bashrc\e[0m"
    echo   "\e[2m"
  fi

  if [ -f "$HOME/.zshrc" ]; then
    echo "> detected: ~/.zshrc"
    echo "> shell function: ⚡ Adding to ~/.zshrc"
    $INSTALL_PATH/mono-cd --install "$HOME/.zshrc"

    echo   "\e[0m"
    echo   "    🚀 Run this, or re-start your zsh terminal:"
    echo   "       $ \e[1msource ~/.zshrc\e[0m"
    echo   "\e[2m"
  fi
}

completed() {
  return 0
}

printf "\e[2m"

create_directory && \
  download && \
  update_permissions && \
  add_to_shell && \
  completed

printf "\e[0m"
