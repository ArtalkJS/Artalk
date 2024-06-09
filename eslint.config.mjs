// @ts-check
import path from 'node:path'
import url from 'node:url'
import { fixupPluginRules } from '@eslint/compat'
// eslint-disable-next-line import-x/namespace
import { FlatCompat } from '@eslint/eslintrc'
import eslintJs from '@eslint/js'
import pluginTS from '@typescript-eslint/eslint-plugin'
import eslintConfigPrettier from 'eslint-config-prettier'
import pluginFunctional from 'eslint-plugin-functional/flat'
import pluginImportX from 'eslint-plugin-import-x'
import pluginReact from 'eslint-plugin-react'
import pluginReactHooks from 'eslint-plugin-react-hooks'
import pluginReactRefresh from 'eslint-plugin-react-refresh'
import pluginVue from 'eslint-plugin-vue'
import globals from 'globals'
import eslintTs from 'typescript-eslint'
import vueParser from 'vue-eslint-parser'

const __dirname = path.dirname(url.fileURLToPath(import.meta.url))
const tsProjects = ['./tsconfig.base.json', './ui/*/tsconfig.json', './docs/*/tsconfig.json']

const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: eslintJs.configs.recommended,
})

export default eslintTs.config(
  eslintJs.configs.recommended,
  ...eslintTs.configs.recommended,
  // @ts-expect-error the type of `pluginVue` is not compatible with the latest `eslint` v9 package yet
  ...pluginVue.configs['flat/recommended'],
  // FIXME: TypeError SEE https://github.com/amilajack/eslint-plugin-compat/pull/609#issuecomment-2123734301
  // ...compat.extends('plugin:compat/recommended'),
  // {
  //   ...pluginFunctional.configs.recommended,
  //   // FIXME: https://github.com/eslint-functional/eslint-plugin-functional/issues/766#issuecomment-1904715609
  //   rules: {
  //     ...pluginFunctional.configs.recommended.rules,
  //     'functional/immutable-data': 'off',
  //     'functional/no-return-void': 'off',
  //   },
  // },

  /* Global Ignores */
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

  /* TypeScript */
  {
    files: ['**/*.{ts,mts,cts,tsx,js,mjs,cjs,vue}'],
    languageOptions: {
      parser: vueParser,
      parserOptions: {
        parser: eslintTs.parser,
        project: tsProjects,
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
      'import-x': pluginImportX,
    },
    rules: {
      ...pluginImportX.configs.recommended.rules,
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/no-unused-vars': 'off',
      '@typescript-eslint/no-unsafe-declaration-merging': 'off',
      'vue/multi-word-component-names': 'off',
      'import-x/no-named-as-default-member': 'off',
      'import-x/no-named-as-default': 'off',
      'import-x/default': 'off', // fix https://github.com/import-js/eslint-plugin-import/issues/1800
      'import-x/order': 'warn',
      // 'import-x/no-default-export': 'warn'
    },
    settings: {
      'import-x/parsers': {
        '@typescript-eslint/parser': ['.ts', '.tsx'],
      },
      'import-x/resolver': {
        typescript: {
          project: tsProjects,
        },
        node: true,
      },
      polyfills: ['AbortController'],
    },
  },

  /* Vue */
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
        nextTick: 'readonly',
      },
    },
  },

  /* React */
  {
    files: ['docs/landing/**/*.tsx'],
    plugins: {
      react: pluginReact,
      // @ts-expect-error SEE https://github.com/facebook/react/issues/28313
      'react-hooks': fixupPluginRules(pluginReactHooks),
      'react-refresh': pluginReactRefresh,
    },
    rules: {
      ...pluginReact.configs.recommended.rules,
      ...pluginReact.configs['jsx-runtime'].rules,
      ...pluginReactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': 'warn',
      'react/prop-types': 'off',
    },
    languageOptions: {
      ...pluginReact.configs.recommended.languageOptions,
    },
    settings: {
      react: {
        version: '18.3',
      },
    },
  },

  eslintConfigPrettier, // disable conflicting rules with Prettier
)
