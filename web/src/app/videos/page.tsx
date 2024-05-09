"use client"
import React, { useState } from 'react';
import { pb } from "@/pb/service_media"
import { Metadata, credentials } from '@grpc/grpc-js';

function FileUploader() {
  const [file, setFile] = useState(null);

  const uploadLargeFile = async (file : any) => {
    const client = new pb.ServiceMediaClient('http://localhost:50055', credentials.createInsecure());
    const request = new pb.UploadRequest();
    request.chunk = (file);

    const metadata = new Metadata();

    const stream = client.uploadLargeFile(metadata, (err : any , response : any ) => {});
    stream.on('data', (response : any) => {
      console.log('Server response:', response.getMessage());
    });
    stream.on('end', () => {
      console.log('Upload complete');
    });
    stream.on('error', (error : any ) => {
      console.error('Error uploading file:', error);
    });

    stream.write(request);
    stream.end();
  };

  const handleFileUpload = (event : any ) => {
    const selectedFile = event.target.files[0];
    setFile(selectedFile);
    if (selectedFile) {
      uploadLargeFile(selectedFile);
    }
  };

  return (
    <div>
      <input type="file" onChange={handleFileUpload} />
    </div>
  );
}

export default FileUploader;
