import colors from 'picocolors'

const LOG_TAG = 'ArtalkPluginKit'

function error(message: string) {
  console.error(colors.red(`[${LOG_TAG}] ${message}`))
}

function warn(message: string) {
  console.warn(colors.yellow(`[${LOG_TAG}] ${message}`))
}

function info(message: string) {
  console.info(colors.blue(`[${LOG_TAG}] ${message}`))
}

function debug(message: string) {
  console.debug(colors.gray(`[${LOG_TAG}] ${message}`))
}

const logger = {
  error,
  warn,
  info,
  debug,
}

export default logger
