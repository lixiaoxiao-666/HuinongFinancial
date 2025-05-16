const { defineConfig } = require('@vue/cli-service')

module.exports = {
  transpileDependencies: true,
  publicPath: './',
  css: {
    loaderOptions: {
      postcss: {
        postcssOptions: {
          plugins: [
            require("postcss-pxtorem")({
              rootValue: 16, // 根元素字体大小
              propList: ["*"] // 需要转换的属性，这里设置为转换所有属性
            })
          ]
        }
      }
    }
  }
};