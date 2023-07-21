import "./App.css";
import { Ping } from "./libs/apis/ping";

function App() {
  void Ping().then((res) => {
    console.log(res);
  });

  const handleImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files?.length === 0 || files === null) {
      return;
    }
    const file = files[0];
    console.log(file);

    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = (e) => {
      const imageBinary = e.target?.result;
      console.log(imageBinary);
    };
  };

  return (
    <>
      <h1>QRコードリーダー</h1>
      <input type="file" accept="image/*" onChange={(e) => handleImage(e)} />
    </>
  );
}

export default App;
