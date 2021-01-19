#!/bin/bash

########## START-OF common.sh
# source me: source <(curl -s "https://gist.githubusercontent.com/cenk1cenk2/0446f3be22a39c9f5fe5ee1cfb3cca63/raw/common.sh?$(date +%s)")

# Constants
RESET='\033[0m'
RED='\033[38;5;1m'
GREEN='\033[38;5;2m'
YELLOW='\033[38;5;3m'
MAGENTA='\033[38;5;5m'
CYAN='\033[38;5;6m'
SEPERATOR="\033[90m-------------------------${RESET}"

log() {
  stdout_print "${RESET}${*}"
}

stdout_print() {
  # 'is_boolean_yes' is defined in libvalidations.sh, but depends on this file so we cannot source it
  local bool="${CENK1CENK2_QUIET:-false}"
  # comparison is performed without regard to the case of alphabetic characters
  shopt -s nocasematch
  if ! [[ "$bool" = 1 || "$bool" =~ ^(yes|true)$ ]]; then
    echo -e "${1}"
  fi
}

log_debug() {
  # 'is_boolean_yes' is defined in libvalidations.sh, but depends on this file so we cannot source it
  local bool="${CENK1CENK2_DEBUG:-false}"
  # comparison is performed without regard to the case of alphabetic characters
  shopt -s nocasematch
  if [[ "$bool" = 1 || "$bool" =~ ^(yes|true)$ ]]; then
    log_this "${1:-}" "${MAGENTA}DEBUG${RESET}" "${2}"
  fi
}

log_this() {
  INFO="${1:-}"
  SCOPE="${2:-}"
  SEPERATOR_INSERT="${3:-}"

  DATA="${INFO}"

  if [ ! -z "${SCOPE}" ] && [ "${SCOPE}" != "false" ]; then
    DATA="[${SCOPE}] ${DATA}"
  fi

  if [ ! -z "${SEPERATOR_INSERT}" ]; then
    if [[ ${SEPERATOR_INSERT} == "top" ]] || [[ ${SEPERATOR_INSERT} == "both" ]]; then
      DATA="${SEPERATOR}\n${DATA}"
    fi

    if [[ ${SEPERATOR_INSERT} == "bottom" ]] || [[ ${SEPERATOR_INSERT} == "both" ]]; then
      DATA="${DATA}\n${SEPERATOR}"
    fi
  fi

  log "${DATA}"
}

log_start() {
  log_this "${1:-}" "${GREEN}START${RESET}" "${2:-}"
}

log_finish() {
  log_this "${1:-}" "${GREEN}FINISH${RESET}" "${2:-}"
}

log_error() {
  log_this "${1:-}" "${RED}ERROR${RESET}" "${2:-}"
}

log_warn() {
  log_this "${1:-}" "${YELLOW}WARN${RESET}" "${2:-}"
}

log_info() {
  log_this "${1:-}" "${CYAN}INFO${RESET}" "${2:-}"
}

log_interrupt() {
  log_this "${1:-}" "${RED}INTERRUPT${RESET}" "${2:-}"
}

log_wait() {
  log_this "${1:-}" "${YELLOW}WAIT${RESET}" "${2:-}"
}

log_divider() {
  stdout_print "${SEPERATOR}"
}

########## END-OF common.sh
