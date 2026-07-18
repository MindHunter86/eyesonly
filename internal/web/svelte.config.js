// svelte.config.js
export default {
    // Svelte options
    extensions: ['.svelte'],
    compilerOptions: {},
    preprocess: [],
    onwarn: (warning, handler) => handler(warning),
    // plugin options
    vitePlugin: {
        exclude: [],
        experimental: {}
    }
};