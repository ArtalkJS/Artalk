import ts from 'typescript'

export function runExportLint(codePath: string, expectedName: string): LintResult {
  const sourceFile = ts.createSourceFile(
    codePath,
    ts.sys.readFile(codePath) || '',
    ts.ScriptTarget.Latest,
    true,
  )

  const exportNames = new Set<string>()
  let containsDefaultExport = false

  const delintNode = (node: ts.Node) => {
    if (ts.isExportAssignment(node)) {
      containsDefaultExport = true
      if (ts.isIdentifier(node.expression)) {
        exportNames.add(node.expression.text)
      }
    } else if (ts.isVariableStatement(node)) {
      if (
        node.modifiers &&
        node.modifiers.some((mod) => mod.kind === ts.SyntaxKind.ExportKeyword)
      ) {
        for (const decl of node.declarationList.declarations) {
          if (ts.isIdentifier(decl.name)) {
            exportNames.add(decl.name.text)
          }
        }
      }
    } else if (ts.isExportDeclaration(node)) {
      if (node.exportClause && ts.isNamedExports(node.exportClause)) {
        for (const elem of node.exportClause.elements) {
          exportNames.add(elem.name.text)
        }
      } else {
        containsDefaultExport = true
      }
    }

    ts.forEachChild(node, delintNode)
  }

  delintNode(sourceFile)

  if (!exportNames.size) {
    return {
      ok: false,
      level: 'error',
      message: 'No plugin export found',
    }
  }

  if (containsDefaultExport) {
    return {
      ok: false,
      level: 'error',
      message: 'Default export is not allowed',
    }
  }

  if (!Array.from(exportNames).includes(expectedName)) {
    return {
      ok: false,
      level: 'error',
      message: `Exported names [${Array.from(exportNames)
        .map((a) => `"${a}"`)
        .join(', ')}] no one is "${expectedName}"`,
    }
  }

  return {
    ok: true,
  }
}
