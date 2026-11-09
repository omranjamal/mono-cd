package mcd

var ShellFunction string = `# start: mono-cd
mcd() {
  TARGETPATH=$("$HOME/.local/share/omranjamal/mono-cd/mono-cd" $1)

  if [ ! -z "${TARGETPATH}" ] ; then
    cd "${TARGETPATH}"
  fi
}

if [ -n "${ZSH_VERSION:-}" ] && command -v compdef >/dev/null 2>&1; then
  _mcd() {
    _arguments '1:search term:'
  }
  compdef _mcd mcd
fi
# end: mono-cd`
