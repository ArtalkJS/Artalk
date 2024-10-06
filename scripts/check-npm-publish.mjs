#!/usr/bin/env node

import { promises as fs } from 'fs'
import path from 'path'
import { execSync } from 'child_process'
import process from 'process'

const __dirname = path.dirname(new URL(import.meta.url).pathname)

// Helper to run shell commands
const runCommand = (command) => {
  try {
    return execSync(command, { encoding: 'utf-8' }).trim()
  } catch (error) {
    return null
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

// Compare local version with the latest npm version
const checkVersionMismatch = async (projectPath) => {
  const packageJsonPath = path.join(projectPath, 'package.json')

  try {
    const packageJsonContent = await fs.readFile(packageJsonPath, 'utf-8')
    const packageJson = JSON.parse(packageJsonContent)
    const localVersion = packageJson.version
    const packageName = packageJson.name

    // Get the latest version from npm using pnpm info
    const npmVersion = runCommand(`pnpm info ${packageName} version`)

    if (localVersion === npmVersion) {
      console.log(`✅ ${packageName} is up to date (${npmVersion})`)
    } else {
      console.log(`❌ ${packageName} is outdated (local: ${localVersion}, npm: ${npmVersion})`)
    }

    if (npmVersion && localVersion !== npmVersion) {
      return { packageName, localVersion, latestVersion: npmVersion }
    }
  } catch (error) {
    console.error(`Failed to read package.json in ${projectPath}:`, error)
  }

  return null
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

// Main function to find outdated packages
const findOutdatedProjects = async () => {
  const specifiedProject = getArgs()
  let projects = await findNodeProjects(path.join(__dirname, '../ui'))
  console.log(`Found ${projects.length} projects under 'ui' directory.\n`)
  console.log('Checking npm publishes...\n')

  // Filter projects by the specified one, if provided
  if (specifiedProject) {
    projects = projects.filter((projectPath) => {
      const packageJsonPath = path.join(projectPath, 'package.json')
      const packageJson = JSON.parse(fs.readFileSync(packageJsonPath, 'utf-8'))
      return packageJson.name === specifiedProject
    })

    if (projects.length === 0) {
      console.log(`Project '${specifiedProject}' not found.`)
      return
    }
  }

  const outdatedProjects = []

  for (const project of projects) {
    const result = await checkVersionMismatch(project)
    if (result) {
      outdatedProjects.push(result)
    }
  }

  console.log('\n==================================================\n')

  if (outdatedProjects.length === 0) {
    console.log('✅ All projects have the latest versions pushed to npm.')
  } else {
    console.log('Projects with outdated versions:\n')
    outdatedProjects.forEach(({ packageName, localVersion, latestVersion }) => {
      console.log(`❌ ${packageName}: Local version ${localVersion}, NPM version ${latestVersion}`)
    })
    process.exit(1)
  }
}

findOutdatedProjects()
