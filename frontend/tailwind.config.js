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
    extend: {},
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
