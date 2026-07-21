#!/usr/bin/env node

import https from 'node:https'
import fs from 'node:fs'
import process from 'node:process'
import path from 'node:path'
import { Buffer } from 'node:buffer'

const __dirname = path.dirname(new URL(import.meta.url).pathname)

const REGISTRY_URL = 'https://github.com/ArtalkJS/Community/releases/latest/download/registry.json'
const DOCS_DIST_PATH = path.join(__dirname, '../docs/docs/.vitepress/dist/plugins')
const MAX_REDIRECTS = 5
const MAX_RESPONSE_SIZE = 5 * 1024 * 1024
const REQUEST_TIMEOUT = 10_000
const ALLOWED_HOSTS = new Set([
  'github.com',
  'objects.githubusercontent.com',
  'release-assets.githubusercontent.com',
])

const download = (url, redirectCount = 0) => {
  return new Promise((resolve, reject) => {
    const parsedUrl = new URL(url)
    if (parsedUrl.protocol !== 'https:' || !ALLOWED_HOSTS.has(parsedUrl.hostname)) {
      reject({ error: true, message: `Refusing registry URL '${url}'` })
      return
    }

    const request = https.get(parsedUrl, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        if (redirectCount >= MAX_REDIRECTS) {
          response.resume()
          reject({ error: true, message: 'Too many registry redirects' })
          return
        }

        const redirectUrl = new URL(response.headers.location, parsedUrl).toString()
        response.resume()
        return resolve(download(redirectUrl, redirectCount + 1))
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

      const chunks = []
      let responseSize = 0
      response.on('data', (chunk) => {
        responseSize += chunk.length
        if (responseSize > MAX_RESPONSE_SIZE) {
          response.destroy()
          reject({ error: true, message: 'Registry response exceeds 5 MiB' })
          return
        }
        chunks.push(chunk)
      })

      response.on('end', () => {
        try {
          resolve(JSON.parse(Buffer.concat(chunks).toString('utf-8'))) // Try to parse the JSON
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

    request.setTimeout(REQUEST_TIMEOUT, () => {
      request.destroy(new Error(`Request timed out after ${REQUEST_TIMEOUT}ms`))
    })
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
