import { transport } from "./transport";
import { createPromiseClient } from "@bufbuild/connect";
import { PingService } from "../../../gen/proto/ping/v1/ping_connectweb";

const client = createPromiseClient(PingService, transport);

export const Ping = async () => {
  const resp = await client.ping({ message: "Hello" }, {});
  return resp;
};
