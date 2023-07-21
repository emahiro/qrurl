import "./App.css";
import { useState } from "react";
import { postQrCode } from "./libs/apis/qrurl";
import titleImg from "./assets/sp_qr_code_man.png";

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
      <h1 className="text-4xl font-bold font-mono">QR Code Reader</h1>
      <div className="logoImg">
        <img src={titleImg} alt="QRコードを読み込む男性" />
      </div>
      <div className="qrCodeBox">
        <form>
          <label className="block">
            <span className="sr-only">Choose File</span>
            <input
              type="file"
              className="block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-gray-700 hover:file:bg-blue-100"
              accept="image/*"
              onChange={(e) => handleImage(e)}
            />
          </label>
        </form>
        <div className="qrCodeResult">
          <p className="text-sm">読み込まれたURL ↓</p>
          <a className="font-mono text-lg" href={url} target="_blank">
            {url}
          </a>
        </div>
      </div>
    </>
  );
}

export default App;
