package mcd

var HelpText = `

# MONO-CD

> CLI utility that lets you interactively pick the
  directory of interest that you want to cd into.

Usage:
  mcd                  Start interactive mode directory picker
  mcd [SEARCH_TERM]    Search by SEARCH_TERM and automatically cd into it, if there is only one match

  [UP] / [DOWN] Arrow Keys (interactive mode only):
    Lets you choose a directory of interest from the list

  Start Typing (interactive mode only):
    This will allow you to filter the suggested directories of interest

Flags:
  -h / --help                  Show this help message

  --shell [ALIAS]              Show the shell function code with name set to ALIAS (optional)
                               Mostly useful for manual installation.

  --install FILE [ALIAS]       Add or update the shell function in a shell startup file
                               like ~/.bashrc; Setting an ALIAS will change the function
                               name from mcd to ALIAS
`
