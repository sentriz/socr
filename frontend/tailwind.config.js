module.exports = {
  purge: {
    enabled: !!process.env.PRODUCTION,
    content: ['index.html', './src/**/*.vue', './src/**/*.css'],
  },
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
    },
  },
  variants: {},
  plugins: [],
}
