import https from 'https'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

fs.copyFile(path.join(__dirname, '../../../../artalk-go.example.yml'), path.join(__dirname, '../src/assets/artalk-go.example.yml'), (err) => {
  if (!err) console.log("\nArtalkGo config file 'artalk-go.example.yml' loaded.\n")
  else console.error("Failed to load config file 'artalk-go.example.yml':\n\n", err, "\n");
})

// const file = fs.createWriteStream()

// https.get(
//   'https://raw.githubusercontent.com/ArtalkJS/ArtalkGo/master/artalk-go.example.yml',
//   (resp) => {
//     resp.pipe(file)

//     file.on('finish', () => {
//       file.close()
//       console.log("\nArtalkGo 'artalk-go.example.yml' file download completed.\n")
//     })
//   }
// ).on('error', (e) => {
//   console.error("Failed to download 'artalk-go.example.yml' file:\n\n", e, "\n");
// })
