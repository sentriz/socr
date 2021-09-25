const scrollbar = require('tailwind-scrollbar')

module.exports = {
  mode: 'jit',
  purge: ['index.html', './components/**/*.vue'],
  theme: {
    screens: {
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
    },
    extend: {
      width: {
        fit: 'fit-content',
      },
      minHeight: {
        32: '8rem',
        36: '9rem',
        40: '10rem',
      },
    },
  },
  plugins: [scrollbar],
}
