package mcd

var ShellFunction string = `# start: mono-cd
mcd() {
  TARGETPATH=$("$HOME/.local/share/omranjamal/mono-cd/mono-cd" $1)

  if [ ! -z "${TARGETPATH}" ] ; then
    cd "${TARGETPATH}"
  fi
}
# end: mono-cd`
