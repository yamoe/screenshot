
const path = require('path');
const webpack = require("webpack");
const HTMLWebpackPlugin = require('html-webpack-plugin');
const dist = path.resolve(__dirname, 'dist')

module.exports = {
  entry: {
    entry: './src/js/main.js'
  },
  output: {
    path: dist,
    filename: '[name].[hash].entry.js',
    chunkFilename: '[name].[hash].entry.js',
  },
  module: {
    rules: [{
      test: /\.js$/,
      include: [
        dist
      ],
      exclude: /node_modules/,
      use: {
        loader: 'babel-loader',
        options: {
          presets: ['env']
        }
      }
    },
    {
      test: /\.(png|jpe?g|gif|svg)(\?.*)?$/,
      use: [
        {
          loader: "file-loader",
          options: {
            context: "./img/",
            outputPath: "img/",
            name: "[path][name].[ext]"
          }
        }
      ]
    }
    ]
  },
  devServer: {
    hot: true,
    inline: true,
    contentBase: "dist",
    overlay: true,
    stats: {
      color: true
    }
  },
  plugins: [
    new webpack.HotModuleReplacementPlugin(),
    new HTMLWebpackPlugin({
      template: "./index.html"
    })
  ],
  devtool: 'source-map',
  mode: 'development'
};