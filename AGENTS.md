# AGENTS.md

This file contains instructions for AI agents (like Jules) working on this repository.

## Documentation Maintenance

*   **Update Documentation:** Whenever you make changes to the code (adding features, changing behaviors, or updating installation steps), you **MUST** verify if the documentation needs to be updated.
*   **Docs Location:** The documentation is located in the `docs/` directory and is a VitePress-based static site.
*   **Source of Truth:** While `README.md` and other root markdown files exist, the `docs/` directory (with markdown files at the root level) should be considered the primary source for the documentation website. If you update root markdown files, ensure the corresponding files in `docs/` are also updated.

## Commit Standards

*   **Conventional Commits:** You **MUST** use [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages.
    *   Format: `<type>[optional scope]: <description>`
    *   Common types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `chore`.
    *   Example: `feat: add new neovim plugin`
    *   Example: `fix(installer): resolve path issue in ubuntu`
