const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    'internal/templates/*.templ',
    'internal/components/*.templ',
  ],
  theme: {
    extend: {
      fontFamily: {
        'mono': ['JetBrains Mono', 'Courier New', 'monospace'],
      },
      colors: {
        primary: colors.blue,
        secondary: colors.yellow,
        neutral: colors.gray,
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
  ],
}