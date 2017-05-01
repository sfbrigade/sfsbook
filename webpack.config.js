var webpack = require('webpack');
var OptimizeCssAssetsPlugin = require('optimize-css-assets-plugin');
var path = require('path');

// This dir represents the directory path of the bundle file output.
var BUILD_DIR = path.resolve(__dirname, 'client/dist/assets');
// This dir holds the directory path of the application's front end.
var APP_DIR = path.resolve(__dirname, 'client/src');

var config = {
    entry: APP_DIR + '/index.js',
    output: {
        path: BUILD_DIR,
        filename: 'bundle.js'
    },
    devServer: {
        inline: true,
        contentBase: 'client/dist',
        port: 3000
    },
    module: {
        loaders: [
            {
                test: /.js$/,
                include: APP_DIR,
                exclude: /(node_modules)/,
                loader: 'babel',
                query: {
                    presets: [ 'latest', 'stage-0', 'react']
                }
            },
            {
                test: /\.json$/,
                exclude: /(node_modules)/,
                loader: 'json-loader'
            },
            {
                test: /.css$/,
                loader: 'style-loader'
            },
            {
                test: /\.scss/,
                loader: 'style-loader!css-loader!postcss-loader!sass-loader'
            }
        ]
    },
    plugins: [
        new OptimizeCssAssetsPlugin({
            assetNameRegExp: /\.optimize\.css$/g,
            cssProcessor: require('cssnano'),
            cssProcessorOptions: {discardComments: {removeAll: true}},
            canPrint: true
        })
    ]
};

module.exports = config;