var merge = require('webpack-merge')
var prodEnv = require('./prod.env')


module.exports = merge(prodEnv, {
    NODE_ENV: '"development"',
    OPENAPI_ADMIN_ADDR: '"http://127.0.0.1:1322"'
})