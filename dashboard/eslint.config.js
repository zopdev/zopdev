// eslint.config.js
import js from '@eslint/js';
import react from 'eslint-plugin-react';
import reactHooks from 'eslint-plugin-react-hooks';
import prettier from 'eslint-plugin-prettier';
import promise from 'eslint-plugin-promise';
import importPlugin from 'eslint-plugin-import';

export default [
  js.configs.recommended,
  {
    ignores: ['dist/**', 'node_modules/**', 'build/**'], // Global ignores that apply to all rules
  },
  {
    // Add specific configuration for tailwind.config.js and other config files
    files: ['tailwind.config.js', 'vite.config.js', '*.config.js'],
    languageOptions: {
      globals: {
        require: 'readonly',
        module: 'readonly',
        process: 'readonly',
        __dirname: 'readonly',
      },
    },
  },
  {
    files: ['**/*.js', '**/*.jsx'],
    ignores: ['tailwind.config.js', 'vite.config.js', '*.config.js'],
    plugins: {
      react,
      'react-hooks': reactHooks,
      prettier,
      promise,
      import: importPlugin,
    },
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      parserOptions: {
        ecmaFeatures: {
          jsx: true,
        },
      },
      globals: {
        React: 'readonly', // Add React to globals
        document: 'readonly',
        window: 'readonly',
        console: 'readonly',
        require: 'readonly',
        module: 'readonly',
        process: 'readonly',
        __dirname: 'readonly',
        navigator: 'readonly',
        localStorage: 'readonly',
        sessionStorage: 'readonly',
        fetch: 'readonly',
        alert: 'readonly',
        MutationObserver: 'readonly',
        reportError: 'readonly',
        setTimeout: 'readonly',
      },
    },
    settings: {
      react: {
        version: 'detect',
      },
      'import/core-modules': ['dayjs', 'dayjs/plugin/timezone', 'dayjs/plugin/utc'],
    },
    rules: {
      'prettier/prettier': 'error',
      'import/first': 'error',
      'import/newline-after-import': 'error',
      'import/no-extraneous-dependencies': ['error', { devDependencies: true }],
      'import/no-duplicates': 'error',
      'react/jsx-no-useless-fragment': 'off',
      'react/jsx-uses-react': 'error',
      'react/jsx-uses-vars': 'error',
      'react/react-in-jsx-scope': 'error',
      'react-hooks/rules-of-hooks': 'error',
      'react-hooks/exhaustive-deps': 'off',
      'no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          ignoreRestSiblings: true,
          caughtErrors: 'none', // This ignores unused variables in catch blocks
        },
      ],
      // Modified to allow console.log during development if needed
      'no-restricted-syntax': [
        'error',
        {
          selector:
            "CallExpression[callee.object.name='console'][callee.property.name!=/^(error)$/]",
          message: 'Unexpected property on console object was called',
        },
      ],
    },
  },
];
