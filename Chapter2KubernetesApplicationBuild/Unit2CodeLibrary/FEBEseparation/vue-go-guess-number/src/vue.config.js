const { createProxyMiddleware } = require("http-proxy-middleware");
module.exports = {
  productionSourceMap: false,
};
module.exports = {
  devServer: {
    proxy: {
      "/api": {
        target: "http://app-go-backend-service.cicd.svc.cluster.local:8081",
        changeOrigin: true,
        pathRewrite: {
          "^/api": "",
        },
      },
    },
  },
};

