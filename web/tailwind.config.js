const scrollbar = require('tailwind-scrollbar')
const colors = require('tailwindcss/colors')

module.exports = {
  content: ['index.html', './components/**/*.vue'],
  theme: {
    fontFamily: {
      sans: ['Inconsolata', 'sans-serif'],
      serif: ['serif'],
      mono: ['monospace'],
    },
    screens: {
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
    },
    extend: {
      minHeight: {
        32: '8rem',
        36: '9rem',
        40: '10rem',
      },
      colors: {
        green: colors.emerald,
        yellow: colors.amber,
        purple: colors.violet,
        gray: colors.neutral,
      },
    },
  },
  plugins: [scrollbar],
}
