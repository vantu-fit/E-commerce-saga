import { ServiceMediaClient } from "@/pb/Service_mediaServiceClientPb";

export default function Video() {
  var client = new ServiceMediaClient("http://localhost:50055")
  
  return (
    <div>Home</div>
  );
}
