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

#### `noCycleDeps`

Circular dependencies should not be allowed in the `provide` method. The method must not inject a dependency that it also provides, including indirect circular references (e.g., `a` -> `b` -> `c` -> `a`).

The best way to deal with this situation is to do some kind of refactor to avoid the cyclic dependencies.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide('foo', (foo) => {}, ['foo'])
}
```

```ts
import type { ArtalkPlugin } from 'artalk'

// foo.ts
const FooPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide('foo', (bar) => {}, ['bar'])
}

// bar.ts
const BarPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide('bar', (foo) => {}, ['foo'])
}
```

**✅ Pass**:

You can introduce a mediator to resolve circular dependencies. The mediator will handle interactions between the dependencies, breaking the direct circular relationship while maintaining their communication through the mediator.

```ts
import type { ArtalkPlugin } from 'artalk'

// foo.ts
const FooPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide('foo', () => {})
}

// bar.ts
const BarPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide('bar', () => {})
}

// mediator.ts
const MediatorPlugin: ArtalkPlugin = (ctx) => {
  ctx.provide(
    'mediator',
    (foo, bar) => {
      // ...
      // interact with foo and bar
    },
    ['foo', 'bar'],
  )
}
```

#### `noLifeCycleEventInNestedBlocks`

Life-cycle event listeners such as `created`, `mounted`, `updated`, and `destroyed` should not be defined inside nested blocks. They must be placed in the top-level scope of the `ArtalkPlugin` arrow function to ensure clarity and maintainability.

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

Event listeners should not be defined inside the `watchConf` effect function. They must be placed outside to ensure proper separation of concerns and to avoid unintended side effects.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  ctx.watchConf(['el'], (conf) => {
    ctx.on('update', () => {})
  })
}
```

**✅ Pass**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  ctx.on('update', () => {})

  ctx.watchConf(['el'], (conf) => {})
}
```

#### `noInjectInNestedBlocks`

The `inject` method should not be called inside nested blocks. It must be used at the top-level scope of the `ArtalkPlugin` arrow function. For better readability and maintainability, it is recommended to place the `inject` call at the beginning of the function.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  const fn = () => {
    const foo = ctx.inject('foo')
  }
}
```

**✅ Pass**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  const foo = ctx.inject('foo')
}
```

#### `noInjectOutsidePlugin`

The `inject` method should not be called outside the `ArtalkPlugin` arrow function. It must be used in the top-level scope of the `ArtalkPlugin` function to ensure the dependency injection remains readable and maintainable.

**⚠️ Fail**:

```ts
function fn(ctx) {
  const foo = ctx.inject('foo')
}
```

**✅ Pass**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {
  const foo = ctx.inject('foo')
}
```

#### `onePluginPerFile`

Multiple plugins should not be defined in the same file. Each plugin must be defined in its own separate file to improve code organization and maintainability.

**⚠️ Fail**:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {}
export const AnotherPlugin: ArtalkPlugin = (ctx) => {}
```

**✅ Pass**:

TestPlugin.ts:

```ts
import type { ArtalkPlugin } from 'artalk'

export const TestPlugin: ArtalkPlugin = (ctx) => {}
```

AnotherPlugin.ts:

```ts
import type { ArtalkPlugin } from 'artalk'

export const AnotherPlugin: ArtalkPlugin = (ctx) => {}
```

## License

[MIT](https://github.com/ArtalkJS/Artalk/blob/master/LICENSE)
