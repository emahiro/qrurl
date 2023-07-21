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
    reader.readAsBinaryString(file);
    reader.onload = async (e) => {
      const binaryStr = e.target?.result;
      const encoded = btoa(binaryStr as string);
      const resp = await postQrCode(encoded);
      console.log(resp);
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
