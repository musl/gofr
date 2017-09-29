const autoprefixer = require('autoprefixer');
const webpack = require('webpack');

const CopyWebpackPlugin = require('copy-webpack-plugin');
const BabelMinifyPlugin = require('babel-minify-webpack-plugin');

module.exports = {
  entry: {
    main: __dirname + '/src/main.js',
  },
  output: {
    filename: 'bundle.js',
    path: __dirname + '/build',
  },
  module: {
    loaders: [
      {
        test: /\.js$/,
        loader: 'babel-loader',
        exclude: /(node_modules)/,
        options: {
          presets: [
            ['env', {'targets': {'browsers': ['last 2 versions']}}],
          ],
        },
      },
      { test: /\.css/, loader: 'style-loader!css-loader' },
      { test: /\.html/, loader: 'ractive-loader' },
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env': {
        NODE_ENV: JSON.stringify('production')
      }
    }),
    new webpack.optimize.OccurrenceOrderPlugin(),
    new BabelMinifyPlugin(),
    new CopyWebpackPlugin([
      { from: __dirname + '/static' },
    ]),
  ],
};
