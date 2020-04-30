const webpack = require("webpack");
const path = require("path");
const CopyPlugin = require("copy-webpack-plugin");

const API_URL = {
  production: JSON.stringify(""),
  development: JSON.stringify(""),
};

const environment =
  process.env.NODE_ENV === "production" ? "production" : "development";

module.exports = {
  entry: path.resolve(__dirname, "src", "index.jsx"),
  output: {
    path: path.resolve(__dirname, "output"),
    filename: "varq.js",
  },
  resolve: {
    extensions: [".js", ".jsx"],
  },
  plugins: [
    new CopyPlugin([
      { from: "src/assets/varq.svg", to: "./assets/varq.svg" },
      { from: "src/assets/favicons/", to: "./" },
    ]),
    new webpack.DefinePlugin({
      API_URL: API_URL[environment],
    }),
  ],
  module: {
    rules: [
      {
        test: /\.jsx/,
        use: {
          loader: "babel-loader",
          options: { presets: ["@babel/preset-react", "@babel/preset-env"] },
        },
      },
      {
        test: /\.scss/,
        use: ["style-loader", "css-loader", "postcss-loader", "sass-loader"],
      },
      {
        test: /\.css/,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.(woff|woff2|eot|ttf|svg)(\?.*$|$)/,
        use: [
          {
            loader: "file-loader",
            options: {
              name: "[name].[ext]",
              outputPath: "fonts/",
            },
          },
        ],
      },
    ],
  },
  devServer: {
    contentBase: "./src",
    publicPath: "/",
  },
};
