// @ts-check
import path from 'node:path'
import url from 'node:url'
import eslintTs from 'typescript-eslint'
import pluginTS from '@typescript-eslint/eslint-plugin'
import pluginVue from 'eslint-plugin-vue'
import vueParser from 'vue-eslint-parser'
import pluginFunctional from 'eslint-plugin-functional/flat'
import globals from 'globals'
import eslintJs from '@eslint/js'
import eslintConfigPrettier from 'eslint-config-prettier'
import { FlatCompat } from '@eslint/eslintrc'
import { fixupPluginRules } from '@eslint/compat'

const __dirname = path.dirname(url.fileURLToPath(import.meta.url))

const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: eslintJs.configs.recommended,
})

function legacyPlugin(name, alias = name) {
  const plugin = compat.plugins(name)[0]?.plugins?.[alias]

  if (!plugin) {
    throw new Error(`Unable to resolve plugin ${name} and/or alias ${alias}`)
  }

  return fixupPluginRules(plugin)
}

export default eslintTs.config(
  eslintJs.configs.recommended,
  ...eslintTs.configs.recommended,

  // @ts-expect-error the type of `pluginVue` is not compatible with the latest `eslint` v9 package yet
  ...pluginVue.configs['flat/recommended'],

  ...compat.extends('plugin:import/typescript'),
  ...compat.extends('plugin:react-hooks/recommended'),

  // FIXME: TypeError SEE https://github.com/amilajack/eslint-plugin-compat/pull/609#issuecomment-2123734301
  // ...compat.extends('plugin:compat/recommended'),

  {
    ...pluginFunctional.configs.recommended,

    // FIXME: https://github.com/eslint-functional/eslint-plugin-functional/issues/766#issuecomment-1904715609
    rules: {
      ...pluginFunctional.configs.recommended.rules,
      'functional/immutable-data': 'off',
      'functional/no-return-void': 'off',
    },
  },

  eslintConfigPrettier,
  {
    // `ignores` key must been defined in a separate object without any other keys
    // see https://eslint.org/docs/latest/use/configure/ignore#ignoring-files
    // and https://github.com/eslint/eslint/discussions/18304
    ignores: [
      '**/node_modules',
      '**/dist',
      'public',
      'local',
      'test',
      '**/.vitepress/cache',
      '**/*.config.js',
      '**/*.config.ts',
      '**/*.d.ts',
    ],
  },
  {
    files: ['**/*.{ts,tsx,js,vue}'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: eslintTs.parser,
        project: ['./ui/*/tsconfig.json', './docs/*/tsconfig.json'],
        tsconfigRootDir: __dirname,
        globals: {
          ...globals.browser,
        },
        ecmaFeatures: {
          jsx: true,
        },
        sourceType: 'module',
        extraFileExtensions: ['.vue'],
      },
    },
    plugins: {
      '@typescript-eslint': pluginTS,

      // FIXME: eslint `import` plugin is not fully support eslint v9 yet
      // see https://github.com/import-js/eslint-plugin-import/issues/2948#issuecomment-2148832701
      import: legacyPlugin('eslint-plugin-import', 'import'),
    },
    rules: {
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
    },
    settings: {
      'import/resolver': {
        typescript: {
          project: ['ui/artalk/tsconfig.json'].map((p) => path.resolve(__dirname, p)),
        },
      },
      polyfills: ['AbortController'],
    },
  },
  {
    files: ['**/*.vue'],
    languageOptions: {
      globals: {
        // Auto Imports Support
        // SEE https://eslint.vuejs.org/user-guide/#auto-imports-support
        // SEE https://github.com/antfu/eslint-config/blob/e32301ac398896f20e1ec1f4f10a334687f8afc8/src/configs/vue.ts#L40-L55
        computed: 'readonly',
        defineEmits: 'readonly',
        defineExpose: 'readonly',
        defineProps: 'readonly',
        onMounted: 'readonly',
        onUnmounted: 'readonly',
        reactive: 'readonly',
        ref: 'readonly',
        shallowReactive: 'readonly',
        shallowRef: 'readonly',
        toRef: 'readonly',
        toRefs: 'readonly',
        watch: 'readonly',
        watchEffect: 'readonly',
        onUpdated: 'readonly',
        onBeforeMount: 'readonly',
        onBeforeUnmount: 'readonly',
        onBeforeUpdate: 'readonly',
        useRoute: 'readonly',
        useRouter: 'readonly',
        useStore: 'readonly',
        useI18n: 'readonly',
      },
    },
    rules: {},
  },
  {
    files: ['**/*.tsx'],
    rules: {},
  },
)
