#!/usr/bin/env node

import https from 'node:https'
import fs from 'node:fs'
import process from 'node:process'
import path from 'node:path'

const __dirname = path.dirname(new URL(import.meta.url).pathname)

const REGISTRY_URL = 'https://github.com/ArtalkJS/Community/releases/latest/download/registry.json'
const DOCS_DIST_PATH = path.join(__dirname, '../docs/docs/.vitepress/dist/plugins')

const download = (url) => {
  return new Promise((resolve, reject) => {
    const request = https.get(url, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        const redirectUrl = response.headers.location
        response.resume()
        return resolve(download(redirectUrl))
      }

      if (response.statusCode !== 200) {
        response.resume()
        reject({
          error: true,
          message: `Failed to get '${url}'`,
          statusCode: response.statusCode,
        })
        return
      }

      let data = ''
      response.on('data', (chunk) => {
        data += chunk
      })

      response.on('end', () => {
        try {
          resolve(JSON.parse(data)) // Try to parse the JSON
        } catch (err) {
          reject({
            error: true,
            message: 'Failed to parse JSON response',
            detail: err.message,
          })
        }
      })
    })

    // Handle request errors
    request.on('error', (err) => {
      reject({
        error: true,
        message: 'Network error occurred',
        detail: err.message,
      })
    })

    // Ensure the request ends
    request.end()
  })
}

download(REGISTRY_URL)
  .then((data) => {
    const registryJSON = JSON.stringify(data, null, 2)

    if (process.argv.includes('--docs-build')) {
      fs.mkdir(DOCS_DIST_PATH, { recursive: true }, (err) => {
        if (err) {
          console.error('❌ Failed to create directory:', err)
          process.exit(1)
        }

        const outputFile = path.join(DOCS_DIST_PATH, 'registry.json')
        fs.writeFile(outputFile, registryJSON, (err) => {
          if (err) {
            console.error('❌ Failed to write registry.json:', err)
            process.exit(1)
          }
          console.log(`✅ The registry.json has been updated! Saved to: "${outputFile}"`)
        })
      })
    } else {
      console.log(registryJSON)
    }
  })
  .catch((err) => {
    console.error(JSON.stringify(err, null, 2))
    process.exit(1)
  })
