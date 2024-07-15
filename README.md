# Kuneiform for Visual Studio Code

Kuneiform VS Code extension provides language support for the Kuneiform language.

## Feature highlights

- Syntax highlighting
- TODO comments
- Code completion: Enabling this extenshion should automatically recommends completions for Kuneiform keywords and variables, or you can manually trigger completions with `Ctrl+Space`,
- Goto Definition: Supports jump to the definition of actions and procedures by right clicking on the action or procedure name and choosing `Go to Definition` from the context menu or `F12`.
- Diagnostics: Syntax errors are highlighted in the editor, and you can see the error message by hovering over the error. You can also see the error message in the `Problems` panel.

This extension uses a [kuneiform language server](https://github.com/kwilteam/kuneiform-ls.git) to provide the above features.
