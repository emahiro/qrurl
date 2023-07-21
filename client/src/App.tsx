import "./App.css";
import { postQrCode } from "./libs/apis/qrurl";

function App() {
  const handleImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files?.length === 0 || files === null) {
      return;
    }
    const file = files[0];

    const reader = new FileReader();
    reader.readAsArrayBuffer(file);
    reader.onload = (e) => {
      const imageBinary = e.target?.result;
      console.log(imageBinary);
      const result = postQrCode(imageBinary as Uint8Array);
      console.log(result);
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
