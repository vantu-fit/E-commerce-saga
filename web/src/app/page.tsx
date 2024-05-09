"use client"
import { ServiceMediaClient } from "@/pb/Service_mediaServiceClientPb";
import { UploadImageRequest } from "@/pb/service_media_pb";
import React from "react";

export default function Home() {
  const [files, setFiles] = React.useState<File>();
  const handleFileChange = () => {
    if (!files) {
      return;
    }
    uploadFile(files);
  };

  const uploadFile = async (file: File) => {
    console.log("file", file);
    var client = new ServiceMediaClient("http://localhost:8086");
    var request = new UploadImageRequest();
    request.setFilename(file.name);

    // Đọc dữ liệu từ file và chuyển đổi thành Uint8Array
    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = () => {
      const data = reader.result as ArrayBuffer;
      const dataArray = new Uint8Array(data);

      // Đặt dữ liệu vào request
      request.setData(dataArray);

      // Đặt các giá trị khác
      request.setProductId("0e96402d-3677-475b-90af-6e459b260fb4");
      request.setAlt("alt");

      // Gửi request
      client.uploadImage(request, {}, (err, response) => {
        if (err) {
          console.log(err);
        } else {
          console.log(response.toObject());
        }
      });
    };
  };

  return (
    <div>
      <input type="file" onChange={e => setFiles(e.target.files![0])}/>
      <button onClick={() => handleFileChange()}> Submit </button>
    </div>
  );
}
