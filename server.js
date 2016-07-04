var webpack = require('webpack')
var WebpackDevServer = require('webpack-dev-server')
var proxy = require('http-proxy').createProxyServer()
var config = require('./webpack.config');

var server = new WebpackDevServer(webpack(config), {
  publicPath: config.output.publicPath,
  hot: true,
  historyApiFallback: true
})

server.listeningApp.on('upgrade', (req, socket) => {
    if (req.url.match('/ws')) {
      console.log('proxying ws', req.url);
      proxy.ws(req, socket, {'target': 'ws://localhost:4000/'});
    }
})

server.listen(3000, 'localhost', (err, result) => {
  if (err) {
    return console.error(err)
  }

  console.log('Listening at http://localhost:3000/')
})
