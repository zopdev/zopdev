// eslint.config.js
import js from '@eslint/js';
import react from 'eslint-plugin-react';
import prettier from 'eslint-plugin-prettier';
import promise from 'eslint-plugin-promise';
import importPlugin from 'eslint-plugin-import';

export default [
  js.configs.recommended,
  {
    files: ['**/*.js', '**/*.jsx'],
    plugins: {
      react,
      prettier,
      promise,
      // In the flat config format, the plugin name should match how you reference it in rules
      // So for "import/first", the plugin key should be "import"
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
      'react-hooks/exhaustive-deps': 'off',
      'no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
          ignoreRestSiblings: true,
        },
      ],
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
