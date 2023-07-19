import "./App.css";
import { createConnectTransport } from "@bufbuild/connect-web";
import { createPromiseClient } from "@bufbuild/connect";
import { PingService } from "../gen/proto/ping/v1/ping_connectweb";

function App() {
  // init gRPC Client
  const transport = createConnectTransport({
    baseUrl: "http://localhost:8080",
  });

  const client = createPromiseClient(PingService, transport);

  void client.ping({ message: "Hello" }, {}).then((resp) => {
    console.log(resp);
  });

  return (
    <>
      <h1>QRURL</h1>
      <p>test</p>
    </>
  );
}

export default App;
