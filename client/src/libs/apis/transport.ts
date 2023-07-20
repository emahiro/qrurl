import { createConnectTransport } from "@bufbuild/connect-web";

// init gRPC Client
export const transport = createConnectTransport({
  baseUrl:
    process.env.NODE_ENV === "production"
      ? process.env.TRANSPORT_URL || ""
      : "http://localhost:3000" || "",
});
