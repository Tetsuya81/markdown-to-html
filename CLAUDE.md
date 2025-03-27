# CLAUDE.md - Development Guidelines

## Project Overview
Simple GitHub-styled Markdown to HTML converter with browser-side rendering using marked.js.

## Commands
- **Run locally**: Open index.html in any browser
- **Deploy**: Push to GitHub Pages (automatic via this repository)
- **Install dependencies**: Not required (uses CDN)
- **Linting**: Use ESLint with `npx eslint index.html` (install with `npm i -g eslint`)
- **Testing**: No automated tests implemented yet

## Code Style Guidelines
- **HTML**: Use double quotes, 4-space indentation, semantic tags
- **CSS**: Follow BEM naming, 4-space indentation, group selectors logically
- **JavaScript**: 
  - ES6+ syntax, 4-space indentation
  - Event delegation for UI components
  - Function declarations for named functions 
  - Descriptive variable names, camelCase convention
- **Error Handling**: Use try/catch blocks around file operations
- **Responsive Design**: Mobile-first approach with media queries
- **Imports**: Prefer CDN links for dependencies (marked.js, GitHub markdown CSS)
- **Documentation**: Add comments for complex functions, document parameters

## Architecture
Single page application with client-side Markdown parsing. Uses vanilla JavaScript with marked.js for conversion and GitHub markdown CSS for styling.