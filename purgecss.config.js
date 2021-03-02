const purgecss = require('@fullhuman/postcss-purgecss')

module.exports = {
  content: ['./views/**/*.html'],
  css: ['./assets/master-tailwind.min.css'],
  output: './assets/tailwind.min.css',
  extractors: [
    {
      extractor: (content) => content.match(/[A-z0-9-:\/]+/g) || [],
      extensions: ["html"]
    }
  ],
  plugins: [
    purgecss({
      content: ['./views/**/*.html']
    })
  ]
}