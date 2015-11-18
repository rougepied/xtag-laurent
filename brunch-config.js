exports.config = {
  modules: {
    wrapper: false
  },
  files: {
    javascripts: {
      joinTo: {
        "vendor.js": /^(bower_components|vendor)/,
        "app.js": "app/**/*.js"
      }
    }
  },
  plugins: {
    uglify: {
      mangle: false
    },
    jshint: {
      pattern: /^app[\\\/].*\.js$/,
      options: {
        browser: true,
        esnext: true,
        strict: false,
        globalstrict: true,
        curly: true,
        eqeqeq: true,
        forin: true,
        undef: true,
        unused: true,
        predef: ['moment', 'xtag', 'console']
      },
      warnOnly: true
    }
  }
}
