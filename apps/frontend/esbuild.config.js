// Custom esbuild configuration to handle CommonJS modules
module.exports = {
  plugins: [
    {
      name: 'commonjs-externals',
      setup(build) {
        // Mark google-protobuf and grpc-web as external to prevent bundling issues
        build.onResolve({ filter: /^(google-protobuf|grpc-web)$/ }, args => {
          return { path: args.path, external: false }
        })
      }
    }
  ],
  define: {
    'global': 'globalThis'
  }
}
