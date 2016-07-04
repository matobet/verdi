var path = require('path')
var webpack = require('webpack')
var CleanWebpackPlugin = require('clean-webpack-plugin')
var CopyWebpackPlugin = require('copy-webpack-plugin')

const isProd = process.env.NODE_ENV === 'production'

var config = module.exports = {
  entry: [
    path.join(__dirname, 'frontend/src/index.jsx')
  ],
  output: {
    path: path.join(__dirname, 'frontend/dist'),
    filename: 'bundle.js'
  },
  module: {
    loaders: [
      { test: /\.jsx?$/, exclude: /node_modules/, loader: 'react-hot!babel' }
    ]
  },
  resolve: {
    extensions: ['', '.js', '.jsx']
  },
  plugins: [
    new CleanWebpackPlugin(['frontend/dist']),
    new CopyWebpackPlugin([
      { from: 'frontend/static' }
    ])
  ]
}

if (isProd) {
  config.devtool = 'source-map'
  config.plugins.push(
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify('production')
      }
    }),
    new webpack.optimize.DedupePlugin(),
    new webpack.optimize.UglifyJsPlugin({
      minimize: true,
      compress: { warnings: false }
    })
  )
} else /* isDev */ {
  config.devServer = {
    host: '0.0.0.0',
    contentBase: './frontend/dist',
    hot: true,
    historyApiFallback: true,
    proxy: {
      '/api/*': 'http://localhost:4000',
      '/ws': 'http://localhost:4000'
    }
  }
  config.entry.unshift(
    'webpack-dev-server/client?http://localhost:3000',
    'webpack/hot/only-dev-server'
  )
  config.plugins.push(
    new webpack.HotModuleReplacementPlugin()
  )
}
