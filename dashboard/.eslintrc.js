module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true,
    es6: true,
    jest: true,
  },
  extends: ['standard', 'next/core-web-vitals', 'prettier', 'plugin:prettier/recommended'],
  parser: '@babel/eslint-parser',
  parserOptions: {
    ecmaFeatures: {
      jsx: true,
    },
    sourceType: 'module',
  },
  settings: {
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
        selector: "CallExpression[callee.object.name='console'][callee.property.name!=/^(error)$/]",
        message: 'Unexpected property on console object was called',
      },
    ],
  },
  plugins: ['prettier', 'react', 'promise'],
  overrides: [
    {
      files: ['**/__test__/**/*.js'],
      rules: {
        'import/no-extraneous-dependencies': ['error', { devDependencies: true }],
      },
    },
  ],
};
