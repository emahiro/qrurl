import "./App.css";
import { useState } from "react";
import { postQrCode } from "./libs/apis/qrurl";

const urlRegex = /https?/;

function App() {
  const [url, setUrl] = useState("");

  const handleImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files?.length === 0 || files === null) {
      return;
    }
    const file = files[0];

    const reader = new FileReader();
    reader.readAsBinaryString(file);
    reader.onload = async (e) => {
      const binaryStr = e.target?.result;
      const encoded = btoa(binaryStr as string);
      const resp = await postQrCode(encoded);

      const url = resp.url;
      const isURL = urlRegex.test(url);
      if (!isURL) {
        alert("読み込まれた QR コードは URL ではありません。");
        return;
      }
      setUrl(resp.url);
    };
  };

  return (
    <>
      <h1>QRコードリーダー</h1>
      <input type="file" accept="image/*" onChange={(e) => handleImage(e)} />
      <p>
        読み込まれたURL ▶
        <a href={url} target="_blank">
          {url}
        </a>
      </p>
    </>
  );
}

export default App;
