module.exports = {
  devServer: {
      host: '127.0.0.1',
      port: 9992,
      proxy: {
          '/echo': {
              target: `ws://172.18.10.46:8929`,
          }
      }
  },
  productionSourceMap: false,
}