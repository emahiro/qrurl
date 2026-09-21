import { createClient } from "@connectrpc/connect";
import { PingService } from "../../../gen/proto/ping/v1/ping_pb";
import { transport } from "./transport";

const client = createClient(PingService, transport);

export const Ping = async () => {
  const resp = await client.ping({});
  return resp;
};
