import https from 'https'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const file = fs.createWriteStream(path.join(__dirname, '../src/assets/artalk-go.example.yml'))

https.get(
  'https://raw.githubusercontent.com/ArtalkJS/ArtalkGo/master/artalk-go.example.yml',
  (resp) => {
    resp.pipe(file)

    file.on('finish', () => {
      file.close()
      console.log("\nArtalkGo 'artalk-go.example.yml' file download completed.\n")
    })
  }
).on('error', (e) => {
  console.error("Failed to download 'artalk-go.example.yml' file:\n\n", e, "\n");
})
