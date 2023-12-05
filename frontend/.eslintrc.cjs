module.exports = {
  root: true,
  env: {
    node: true,
    es2021: true
  },
  extends: [
    'eslint:recommended',
    '@vue/typescript/recommended',
    'plugin:vue/vue3-essential',
    '@vue/prettier',
    '@vue/standard'
  ],
  parserOptions: {
    parser: '@typescript-eslint/parser',
    ecmaVersion: 2020
    // project: resolve(__dirname, './tsconfig.json'),
    // tsconfigRootDir: __dirname,
    // extraFileExtensions: ['.vue']
  },
  rules: {
    // 避免auto import報錯，暫時關閉
    'no-undef': 'off',
    'import/no-unresolved': 'off',
    'import/named': 'off',
    'no-unused-vars': 'off',
    '@typescript-eslint/no-unused-vars': 'off',
    'vue/multi-word-component-names': 'off',
    'space-before-function-paren': 'off',
    '@typescript-eslint/no-explicit-any': 'off'
    // 這邊可以自訂規則
  }
}
