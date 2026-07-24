#!/usr/bin/env node

import { promises as fs, readFileSync } from 'node:fs'
import path from 'node:path'
import process from 'node:process'
import { fileURLToPath } from 'node:url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const REGISTRY_URL = 'https://registry.npmjs.org/'
const REQUEST_TIMEOUT = 10_000

const getPublishedPackage = async (packageName) => {
  const packageUrl = new URL(encodeURIComponent(packageName), REGISTRY_URL)
  const response = await fetch(packageUrl, {
    headers: { Accept: 'application/vnd.npm.install-v1+json' },
    signal: AbortSignal.timeout(REQUEST_TIMEOUT),
  })

  if (response.status === 404) {
    return { latestVersion: null, versions: new Set() }
  }
  if (!response.ok) {
    throw new Error(`npm registry returned HTTP ${response.status} for ${packageName}`)
  }

  const packageMetadata = await response.json()
  return {
    latestVersion: packageMetadata['dist-tags']?.latest || null,
    versions: new Set(Object.keys(packageMetadata.versions || {})),
  }
}

// Recursively find all directories containing package.json
const findNodeProjects = async (dir) => {
  const subdirs = await fs.readdir(dir, { withFileTypes: true })
  const projects = []

  for (const subdir of subdirs) {
    const res = path.resolve(dir, subdir.name)
    if (subdir.isDirectory()) {
      const packageJsonPath = path.join(res, 'package.json')
      try {
        await fs.access(packageJsonPath)
        const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
        const packageJson = JSON.parse(packageJsonContent)
        if (packageJson.private) continue // ignore private packages
        projects.push(res)
      } catch (err) {
        const nestedProjects = await findNodeProjects(res)
        projects.push(...nestedProjects)
      }
    }
  }

  return projects
}

// Check whether the exact local version has already been published to npm.
const checkVersionPublished = async (projectPath) => {
  const packageJsonPath = path.join(projectPath, 'package.json')

  try {
    const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
    const packageJson = JSON.parse(packageJsonContent)
    const localVersion = packageJson.version
    const packageName = packageJson.name

    const { latestVersion, versions } = await getPublishedPackage(packageName)

    if (versions.has(localVersion)) {
      const latestSuffix =
        localVersion === latestVersion ? '' : `; npm latest is ${latestVersion || 'unset'}`
      console.log(`✅ ${packageName}@${localVersion} is published${latestSuffix}`)
      return null
    }

    console.log(
      `❌ ${packageName}@${localVersion} is not published (npm latest: ${latestVersion || 'none'})`,
    )
    return { packageName, localVersion, latestVersion }
  } catch (error) {
    throw new Error(`Failed to check npm version for ${projectPath}`, { cause: error })
  }
}

// Parse command-line arguments to get the project name if provided
const getArgs = () => {
  const args = process.argv.slice(2)
  let specifiedProject = null

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '-F' && i + 1 < args.length) {
      specifiedProject = args[i + 1]
    }
  }

  return specifiedProject
}

// Main function to find unpublished local package versions
const findUnpublishedProjects = async () => {
  const specifiedProject = getArgs()
  let projects = await findNodeProjects(path.join(__dirname, '../ui'))
  console.log(`Found ${projects.length} projects under 'ui' directory.\n`)
  console.log('Checking npm publishes...\n')

  // Filter projects by the specified one, if provided
  if (specifiedProject) {
    projects = projects.filter((projectPath) => {
      const packageJsonPath = path.join(projectPath, 'package.json')
      const packageJson = JSON.parse(readFileSync(packageJsonPath, 'utf-8'))
      return packageJson.name === specifiedProject
    })

    if (projects.length === 0) {
      console.log(`Project '${specifiedProject}' not found.`)
      return
    }
  }

  const unpublishedProjects = []

  for (const project of projects) {
    const result = await checkVersionPublished(project)
    if (result) {
      unpublishedProjects.push(result)
    }
  }

  console.log('\n==================================================\n')

  if (unpublishedProjects.length === 0) {
    console.log('✅ Every local package version is published to npm.')
  } else {
    console.log('Local package versions not published to npm:\n')
    unpublishedProjects.forEach(({ packageName, localVersion, latestVersion }) => {
      console.log(`❌ ${packageName}: local ${localVersion}, npm latest ${latestVersion || 'none'}`)
    })
    process.exit(1)
  }
}

findUnpublishedProjects().catch((error) => {
  console.error(error)
  process.exitCode = 1
})
