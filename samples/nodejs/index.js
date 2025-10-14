#!/usr/bin/env node
// Print current package version (from package.json) then exit
const fs = require('fs')
const path = require('path')

try {
	const pkgPath = path.resolve(__dirname, 'package.json')
	const raw = fs.readFileSync(pkgPath, 'utf8')
	const pkg = JSON.parse(raw)
	const version = pkg && pkg.version ? String(pkg.version) : ''
	if (!version) {
		console.error('version not found')
		process.exit(2)
	}
	console.log(version)
	process.exit(0)
} catch (err) {
	console.error('error reading package.json:', err.message)
	process.exit(1)
}
