const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const { CleanWebpackPlugin } = require('clean-webpack-plugin');
const CopyWebpackPlugin = require('copy-webpack-plugin');

module.exports = (env, argv) => {
  const isProduction = argv.mode === 'production';
  
  return {
    // Entry point for the extension
    entry: './src/index.tsx',
    
    // Output configuration
    output: {
      filename: 'bundle.js',
      path: path.resolve(__dirname, 'dist')
    },
    
    // Development server
    devServer: {
      static: {
        directory: path.join(__dirname, 'dist'),
      },
      compress: true,
      port: 9000,
      hot: true,
      open: true
    },
    
    // Module resolution
    resolve: {
      extensions: ['.ts', '.tsx', '.js', '.jsx'],
      alias: {
        '@components': path.resolve(__dirname, 'src/components'),
        '@hooks': path.resolve(__dirname, 'src/hooks'),
        '@types': path.resolve(__dirname, 'src/types'),
        '@services': path.resolve(__dirname, 'src/services'),
        '@utils': path.resolve(__dirname, 'src/utils')
      }
    },
    
    // Loaders
    module: {
      rules: [
        // TypeScript and React
        {
          test: /\.tsx?$/,
          use: 'ts-loader',
          exclude: /node_modules/
        },
        // CSS
        {
          test: /\.css$/,
          use: ['style-loader', 'css-loader']
        },
        // Images and other assets
        {
          test: /\.(png|svg|jpg|jpeg|gif)$/i,
          type: 'asset/resource'
        },
        // Fonts
        {
          test: /\.(woff|woff2|eot|ttf|otf)$/i,
          type: 'asset/resource'
        }
      ]
    },
    
    // Plugins
    plugins: [
      // Clean dist folder before build
      new CleanWebpackPlugin(),
      
      // Generate HTML file
      new HtmlWebpackPlugin({
        template: './src/index.html',
        filename: 'snowflake-dropdown.html',
        inject: 'body'
      }),
      
      // Copy static files
      new CopyWebpackPlugin({
        patterns: [
          {
            from: 'images',
            to: 'images'
          }
        ]
      })
    ],
    
    // Externals - Don't bundle these, Azure DevOps provides them
    externals: [
      /^VSS\/.*/,
      /^TFS\/.*/,
      /^q$/,
      {
        'react': 'React',
        'react-dom': 'ReactDOM'
      }
    ],
    
    // Source maps for debugging
    devtool: isProduction ? false : 'source-map',
    
    // Optimization
    optimization: {
      minimize: isProduction,
      // Don't split chunks for extensions
      splitChunks: false
    }
  };
};