const esbuild = require('esbuild');

esbuild.build({
  entryPoints: ['./client/client.js'], // Adjust the path to your main entry point
  bundle: true,
  platform: 'node',
  target: 'node14', // Adjust based on your target Node version
  outfile: 'dist/extension.js', // Output bundled file
  external: ['vscode'], // Exclude VS Code API from the bundle
  minify: true,
}).catch(() => process.exit(1));
