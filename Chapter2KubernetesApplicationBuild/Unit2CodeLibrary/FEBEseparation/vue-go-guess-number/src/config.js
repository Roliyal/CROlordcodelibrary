// config.js
const isProduction = process.env.NODE_ENV === "production";

const config = {
    apiUrl: isProduction ? "http://app-go-backend-service.cicd.svc.cluster.local:8081/check-guess" : "http://127.0.0.1:8081/check-guess",
};

export default config;
