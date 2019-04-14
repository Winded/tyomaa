const path = require('path');

module.exports = {
    entry: './src/index.ts',
    module: {
        rules: [
            {
                test: /\.tsx?$/,
                use: 'ts-loader',
                exclude: /node_modules/
            },
            {
                test: /\.(css|png|jpg|gif)$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: '[path][name].[ext]'
                        },
                    },
                ],
            },
            {
                test: /index\.html$/,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: 'index.html'
                        },
                    },
                ],
            },
        ]
    },
    devServer: {
        contentBase: 'dist/www',
        historyApiFallback: true,
        port: 80,
        host: '0.0.0.0',
        disableHostCheck: true,
    },
    resolve: {
        extensions: ['.tsx', '.ts', '.js']
    },
    output: {
        filename: 'bundle.js',
        path: path.resolve(__dirname, 'dist/www')
    }
};