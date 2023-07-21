import { transport } from "./transport";
import { createPromiseClient } from "@bufbuild/connect";
import { QrUrlService } from "../../../gen/proto/qrurl/v1/qrurl_connectweb";

const client = createPromiseClient(QrUrlService, transport);

export const postQrCode = async (image: Uint8Array) => {
  const resp = await client.postQrCode({ image: image }, {});
  return resp;
};
