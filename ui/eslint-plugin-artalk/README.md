# eslint-plugin-artalk [![npm](https://img.shields.io/npm/v/eslint-plugin-artalk)](https://www.npmjs.com/package/eslint-plugin-artalk)

The ESLint plugin enforcing Artalk's development conventions.

It is a part of the [Plugin Development Kit](https://artalk.js.org/develop/plugin.html) for Artalk.

## Installation

```bash
pnpm add -D eslint-plugin-artalk
```

Since Artalk development is based on TypeScript and the plugin relies on it, you need to install `typescript` and `@typescript-eslint/parser`. For more details, refer to [TypeScript ESLint](https://typescript-eslint.io/getting-started/).

### Flat Configuration

Modify the `eslint.config.mjs` file in your project:

<!-- prettier-ignore -->
```js
import eslintJs from '@eslint/js'
import eslintTs from 'typescript-eslint'
import pluginArtalk from 'eslint-plugin-artalk'

export default eslintTs.config(
  eslintJs.configs.recommended,
  ...eslintTs.configs.recommended,
  {
    files: ['**/*.{ts,mts,cts,tsx,js,mjs,cjs}'],
    languageOptions: {
      parser: eslintTs.parser,
      parserOptions: {
        project: './tsconfig.json',
        tsconfigRootDir: __dirname,
        sourceType: 'module',
      },
    },
    plugins: {
      artalk: pluginArtalk,
    },
    rules: {
      ...pluginArtalk.configs.recommended.rules,
    },
  }
)
```

<!-- prettier-ignore-end -->

### Custom Configuration

You can customize the rules by modifying the `rules` field in the configuration:

```js
{
  plugins: {
    artalk: pluginArtalk,
  },
  rules: {
    'artalk/artalk-plugin': 'error',
  },
}
```

## Valid and Invalid Examples

### Rule `artalk-plugin`

The ESLint rule `artalk/artalk-plugin` enforces the conventions for Artalk plugins.

The ESLint rule is only enabled when a TypeScript file imports the `ArtalkPlugin` type from the `artalk` package and defines an arrow function variable with the type `ArtalkPlugin`, such as `const TestPlugin: ArtalkPlugin = (ctx) => {}`. The variable type must be `ArtalkPlugin`.

#### `noLifeCycleEventInNestedBlocks`

Should not allow life-cycle event listeners to be defined inside nested blocks.

The life-cycle event listeners are `created`, `mounted`, `updated`, and `destroyed` must be defined in the top-level scope of the ArtalkPlugin arrow function.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  const foo = () => {
    const bar = () => {
      ctx.on('updated', () => {})
    }
  }
}
```

**✅ Pass**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  ctx.on('updated', () => {})
}
```

#### `noEventInWatchConf`

Should not allow event listeners to be defined inside watchConf effect function.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['el'], (conf) => {
    ctx.on('update', () => {})
  })
}
```

## License

[MIT](https://github.com/ArtalkJS/Artalk/blob/master/LICENSE)
