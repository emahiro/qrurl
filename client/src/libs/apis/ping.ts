import { transport } from "./transport";
import { createPromiseClient } from "@bufbuild/connect";
import { PingService } from "../../../gen/proto/ping/v1/ping_connectweb";

const client = createPromiseClient(PingService, transport);

export const Ping = async () => {
  console.log(import.meta.env.VITE_TRANSPORT_URL);

  const resp = await client.ping({ message: "Hello" }, {});
  return resp;
};
