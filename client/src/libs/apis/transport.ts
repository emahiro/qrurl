import { createConnectTransport } from "@bufbuild/connect-web";

// init gRPC Client
export const transport = createConnectTransport({
  baseUrl: import.meta.env.PROD
    ? (import.meta.env.VITE_TRANSPORT_URL as string)
    : "http://localhost:8080",
});
