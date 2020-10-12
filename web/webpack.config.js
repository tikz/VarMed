const webpack = require("webpack");
const path = require("path");
const package = require("./package.json");
const CopyPlugin = require("copy-webpack-plugin");

const API_URL = {
  prod: JSON.stringify(""),
  dev: JSON.stringify("http://localhost:8888"),
};

module.exports = {
  entry: path.resolve(__dirname, "src", "index.jsx"),
  output: {
    path: path.resolve(__dirname, "output"),
    filename: "varmed.js",
  },
  resolve: {
    extensions: [".js", ".jsx"],
  },
  plugins: [
    new CopyPlugin([
      { from: "src/assets/varmed.svg", to: "./assets/varmed.svg" },
      { from: "src/assets/favicons/", to: "./" },
    ]),
    new webpack.DefinePlugin({
      API_URL: API_URL[process.env.NODE_ENV === "dev" ? "dev" : "prod"],
      JOBS_KEY: JSON.stringify("jobs-" + package.version),
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
