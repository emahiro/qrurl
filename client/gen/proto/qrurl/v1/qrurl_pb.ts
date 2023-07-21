// @generated by protoc-gen-es v1.3.0 with parameter "target=ts"
// @generated from file proto/qrurl/v1/qrurl.proto (package qrurl.v1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import type { BinaryReadOptions, FieldList, JsonReadOptions, JsonValue, PartialMessage, PlainMessage } from "@bufbuild/protobuf";
import { Message, proto3 } from "@bufbuild/protobuf";

/**
 * @generated from message qrurl.v1.PostQrCodeRequest
 */
export class PostQrCodeRequest extends Message<PostQrCodeRequest> {
  /**
   * @generated from field: bytes image = 1;
   */
  image = new Uint8Array(0);

  constructor(data?: PartialMessage<PostQrCodeRequest>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "qrurl.v1.PostQrCodeRequest";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "image", kind: "scalar", T: 12 /* ScalarType.BYTES */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PostQrCodeRequest {
    return new PostQrCodeRequest().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PostQrCodeRequest {
    return new PostQrCodeRequest().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PostQrCodeRequest {
    return new PostQrCodeRequest().fromJsonString(jsonString, options);
  }

  static equals(a: PostQrCodeRequest | PlainMessage<PostQrCodeRequest> | undefined, b: PostQrCodeRequest | PlainMessage<PostQrCodeRequest> | undefined): boolean {
    return proto3.util.equals(PostQrCodeRequest, a, b);
  }
}

/**
 * @generated from message qrurl.v1.PostQrCodeResponse
 */
export class PostQrCodeResponse extends Message<PostQrCodeResponse> {
  /**
   * @generated from field: string url = 1;
   */
  url = "";

  constructor(data?: PartialMessage<PostQrCodeResponse>) {
    super();
    proto3.util.initPartial(data, this);
  }

  static readonly runtime: typeof proto3 = proto3;
  static readonly typeName = "qrurl.v1.PostQrCodeResponse";
  static readonly fields: FieldList = proto3.util.newFieldList(() => [
    { no: 1, name: "url", kind: "scalar", T: 9 /* ScalarType.STRING */ },
  ]);

  static fromBinary(bytes: Uint8Array, options?: Partial<BinaryReadOptions>): PostQrCodeResponse {
    return new PostQrCodeResponse().fromBinary(bytes, options);
  }

  static fromJson(jsonValue: JsonValue, options?: Partial<JsonReadOptions>): PostQrCodeResponse {
    return new PostQrCodeResponse().fromJson(jsonValue, options);
  }

  static fromJsonString(jsonString: string, options?: Partial<JsonReadOptions>): PostQrCodeResponse {
    return new PostQrCodeResponse().fromJsonString(jsonString, options);
  }

  static equals(a: PostQrCodeResponse | PlainMessage<PostQrCodeResponse> | undefined, b: PostQrCodeResponse | PlainMessage<PostQrCodeResponse> | undefined): boolean {
    return proto3.util.equals(PostQrCodeResponse, a, b);
  }
}

