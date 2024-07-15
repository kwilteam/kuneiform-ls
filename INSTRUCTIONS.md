# Build Instructions

## Package

To generate a vsix extension image, run either of the following commands:

```bash
npm run package
```

or

```bash
cd server && ./build.sh
cd ../
vsce package
```

You can directly upload the vsix image into your vscode extensions and start using the Kuneiform Language Server Extension.

## Publish

To publish the extension to the marketplace, run the following command:

```bash
vsce publish
```

or

```bash
npm run publish
```

To publish the extension manually, you can directly upload the vsix image into the [VScode Marketplace Publisher Page]{<https://marketplace.visualstudio.com/manage>}.
For more detailed instructions refer to the [VSCode Extension Publishing Guide]{<https://code.visualstudio.com/api/working-with-extensions/publishing-extension>}.
