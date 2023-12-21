import https from 'https'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

fs.copyFile(path.join(__dirname, '../../../../artalk.example.yml'), path.join(__dirname, '../src/assets/artalk.example.yml'), (err) => {
  if (!err) console.log("\nArtalk config file 'artalk.example.yml' loaded.\n")
  else console.error("Failed to load config file 'artalk.example.yml':\n\n", err, "\n");
})

// const file = fs.createWriteStream()

// https.get(
//   'https://raw.githubusercontent.com/ArtalkJS/Artalk/master/artalk.example.yml',
//   (resp) => {
//     resp.pipe(file)

//     file.on('finish', () => {
//       file.close()
//       console.log("\nArtalk 'artalk.example.yml' file download completed.\n")
//     })
//   }
// ).on('error', (e) => {
//   console.error("Failed to download 'artalk.example.yml' file:\n\n", e, "\n");
// })
