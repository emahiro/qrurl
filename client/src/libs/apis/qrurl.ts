import { createClient } from "@connectrpc/connect";
import { QrUrlService } from "../../../gen/proto/qrurl/v1/qrurl_pb";
import { transport } from "./transport";

const client = createClient(QrUrlService, transport);

export const postQrCode = async (image: string) => {
  const resp = await client.postQrCode({ image });
  return resp;
};
