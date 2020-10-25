module.exports = {
  purge: {
    enabled: !!process.env.PRODUCTION,
    content: [
      'index.html',
      './src/**/*.vue',
      './src/**/*.css'
    ],
  },
  theme: {
    extend: {
      width: {
        'fit': 'fit-content'
      }
    },
  },
  variants: {},
  plugins: [],
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
    defaultLineHeights: true,
    standardFontWeights: true,
  },
}
