interface RuntimeBootConfig {
  test: string
}

interface LintResult {
  ok: boolean
  level?: 'error' | 'warn' | 'info'
  message?: string
}
