'use strict';

/**
 * Development-only webpack settings.
 */
const webpack = require('webpack');
const config = require('./webpack.config');

config.devServer = {
  hot: true,
  inline: true,
  contentBase: './build',
  historyApiFallback: true,
  proxy: { '/api': 'http://127.0.0.1:8081' },
};

config.devtool = 'cheap-module-eval-source-map';

config.plugins = [
  new webpack.optimize.OccurrenceOrderPlugin(),
  new webpack.HotModuleReplacementPlugin(),
  new webpack.NoEmitOnErrorsPlugin(),
  new webpack.NamedModulesPlugin(),
  new webpack.DefinePlugin({
    'process.env': {
      NODE_ENV: JSON.stringify('development')
    }
  }),
];

module.exports = config;
