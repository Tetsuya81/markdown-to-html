# Markdown to HTML Converter

A simple GitHub-styled Markdown to HTML converter with browser-side rendering.

https://tetsuya81.github.io/markdown-to-html/index.html


## Features

- Converts Markdown text to HTML in real-time
- GitHub-styled rendering
- Client-side processing (no server required)
- Mobile-responsive design

## Usage

1. Open `index.html` in any browser
2. Paste or type Markdown in the editor
3. See the HTML output instantly

## Dependencies

- [marked.js](https://marked.js.org/) - Markdown parser and compiler
- GitHub markdown CSS (via CDN)

## Development

- **Run locally (browser only)**: Open index.html in any browser
- **Run with Go server**: Run `go run server.go` and visit http://localhost:8080
  - Server features:
    - Automatically opens browser when started
    - Hot reload - automatically refreshes browser when files change
    - Monitors HTML, CSS, JS, and MD files for changes
- **Deploy**: Push to GitHub Pages (automatic via this repository)
- **Linting**: Use ESLint with `npx eslint index.html`

## License

See the [LICENSE](LICENSE) file for details.
