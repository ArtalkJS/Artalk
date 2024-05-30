// TODO
export function inject({ config }: { config: Record<string, any> }) {
  const h2 = document.createElement('h2')
  h2.textContent = 'Hello from plugin-kit-runtime'
  document.body.appendChild(h2)

  const p = document.createElement('p')
  p.textContent = `config: ${JSON.stringify(config)}`
  document.body.appendChild(p)

  import.meta.hot?.send('artalk-plugin-kit:remote-add')

  import.meta.hot?.on('artalk-plugin-kit:update', (newConfig: RuntimeBootConfig) => {
    p.textContent = `config: ${JSON.stringify(newConfig)}`
  })
}
