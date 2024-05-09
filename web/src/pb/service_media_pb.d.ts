import * as jspb from 'google-protobuf'

import * as google_api_annotations_pb from './google/api/annotations_pb'; // proto import: "google/api/annotations.proto"
import * as protoc$gen$openapiv2_options_annotations_pb from './protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"


export class UploadImageRequest extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): UploadImageRequest;

  getFilename(): string;
  setFilename(value: string): UploadImageRequest;

  getAlt(): string;
  setAlt(value: string): UploadImageRequest;

  getProductId(): string;
  setProductId(value: string): UploadImageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageRequest): UploadImageRequest.AsObject;
  static serializeBinaryToWriter(message: UploadImageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageRequest;
  static deserializeBinaryFromReader(message: UploadImageRequest, reader: jspb.BinaryReader): UploadImageRequest;
}

export namespace UploadImageRequest {
  export type AsObject = {
    data: Uint8Array | string,
    filename: string,
    alt: string,
    productId: string,
  }
}

export class UploadImageResponse extends jspb.Message {
  getUrl(): string;
  setUrl(value: string): UploadImageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageResponse): UploadImageResponse.AsObject;
  static serializeBinaryToWriter(message: UploadImageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageResponse;
  static deserializeBinaryFromReader(message: UploadImageResponse, reader: jspb.BinaryReader): UploadImageResponse;
}

export namespace UploadImageResponse {
  export type AsObject = {
    url: string,
  }
}

export class DeleteImageRequest extends jspb.Message {
  getFilename(): string;
  setFilename(value: string): DeleteImageRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteImageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteImageRequest): DeleteImageRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteImageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteImageRequest;
  static deserializeBinaryFromReader(message: DeleteImageRequest, reader: jspb.BinaryReader): DeleteImageRequest;
}

export namespace DeleteImageRequest {
  export type AsObject = {
    filename: string,
  }
}

export class DeleteImageResponse extends jspb.Message {
  getSuccess(): string;
  setSuccess(value: string): DeleteImageResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteImageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteImageResponse): DeleteImageResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteImageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteImageResponse;
  static deserializeBinaryFromReader(message: DeleteImageResponse, reader: jspb.BinaryReader): DeleteImageResponse;
}

export namespace DeleteImageResponse {
  export type AsObject = {
    success: string,
  }
}

export class UploadRequest extends jspb.Message {
  getChunk(): Uint8Array | string;
  getChunk_asU8(): Uint8Array;
  getChunk_asB64(): string;
  setChunk(value: Uint8Array | string): UploadRequest;

  getSequenceNumber(): number;
  setSequenceNumber(value: number): UploadRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadRequest): UploadRequest.AsObject;
  static serializeBinaryToWriter(message: UploadRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadRequest;
  static deserializeBinaryFromReader(message: UploadRequest, reader: jspb.BinaryReader): UploadRequest;
}

export namespace UploadRequest {
  export type AsObject = {
    chunk: Uint8Array | string,
    sequenceNumber: number,
  }
}

export class UploadResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): UploadResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadResponse): UploadResponse.AsObject;
  static serializeBinaryToWriter(message: UploadResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadResponse;
  static deserializeBinaryFromReader(message: UploadResponse, reader: jspb.BinaryReader): UploadResponse;
}

export namespace UploadResponse {
  export type AsObject = {
    message: string,
  }
}

