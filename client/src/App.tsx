import "./App.css";
import { Ping } from "./libs/apis/ping";

function App() {
  void Ping().then((res) => {
    console.log(res);
  });

  return (
    <>
      <h1>QRコードリーダー</h1>
      <p>test</p>
    </>
  );
}

export default App;
