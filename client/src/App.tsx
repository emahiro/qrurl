import "./App.css";
import { useState, useEffect } from "react";
import { postQrCode } from "./libs/apis/qrurl";
import titleImg from "./assets/sp_qr_code_man.png";

const urlRegex = /https?/;

function App() {
  const [url, setUrl] = useState("");
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    const savedTheme = localStorage.getItem("theme");
    const prefersDark = window.matchMedia(
      "(prefers-color-scheme: dark)",
    ).matches;
    setIsDarkMode(savedTheme === "dark" || (!savedTheme && prefersDark));
  }, []);

  useEffect(() => {
    document.documentElement.classList.toggle("dark", isDarkMode);
    localStorage.setItem("theme", isDarkMode ? "dark" : "light");
  }, [isDarkMode]);

  const handleImage = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (files?.length === 0 || files === null) {
      return;
    }
    const file = files[0];

    setIsLoading(true);
    setError("");
    setUrl("");

    try {
      const reader = new FileReader();
      reader.readAsBinaryString(file);
      reader.onload = async (e) => {
        const binaryStr = e.target?.result;
        const encoded = btoa(binaryStr as string);
        const resp = await postQrCode(encoded);

        const url = resp.url;
        const isURL = urlRegex.test(url);
        if (!isURL) {
          setError("読み込まれた QR コードは URL ではありません。");
          return;
        }
        setUrl(resp.url);
      };
    } catch (err) {
      setError("QR コードの読み込みに失敗しました。");
    } finally {
      setIsLoading(false);
    }
  };

  const toggleTheme = () => {
    setIsDarkMode(!isDarkMode);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 transition-colors duration-300">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {/* Header */}
        <header className="text-center mb-12">
          <div className="flex justify-between items-center mb-8">
            <div className="flex-1"></div>
            <h1 className="text-4xl md:text-5xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              QR Code Reader
            </h1>
            <div className="flex-1 flex justify-end">
              <button
                onClick={toggleTheme}
                className="p-3 rounded-full bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-all duration-200 border border-gray-200 dark:border-gray-700"
                aria-label="Toggle theme"
              >
                {isDarkMode ? (
                  <svg
                    className="w-5 h-5 text-yellow-500"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z"
                      clipRule="evenodd"
                    />
                  </svg>
                ) : (
                  <svg
                    className="w-5 h-5 text-gray-700"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" />
                  </svg>
                )}
              </button>
            </div>
          </div>
          <p className="text-lg text-gray-600 dark:text-gray-300 mb-8">
            QR コードから URL を簡単に読み取ることができます
          </p>
        </header>

        {/* Main Content */}
        <main className="grid lg:grid-cols-2 gap-8 items-start">
          {/* Left Side - Image and Upload */}
          <div className="space-y-8">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8 border border-gray-200 dark:border-gray-700">
              <div className="text-center mb-6">
                <div className="mx-auto w-48 h-48 mb-6 rounded-2xl overflow-hidden shadow-lg">
                  <img
                    src={titleImg}
                    alt="QRコードを読み込む男性"
                    className="w-full h-full object-cover"
                  />
                </div>
              </div>

              <div className="space-y-6">
                <div className="text-center">
                  <h2 className="text-2xl font-semibold text-gray-800 dark:text-white mb-2">
                    QR コードをアップロード
                  </h2>
                  <p className="text-gray-600 dark:text-gray-300">
                    画像ファイルを選択してください
                  </p>
                </div>

                <form className="space-y-4">
                  <label className="block">
                    <span className="sr-only">Choose File</span>
                    <div className="relative">
                      <input
                        type="file"
                        className="block w-full text-sm text-gray-500 dark:text-gray-400 file:mr-4 file:py-3 file:px-6 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100 dark:file:bg-blue-900 dark:file:text-blue-300 dark:hover:file:bg-blue-800 file:transition-colors file:duration-200 cursor-pointer"
                        accept="image/*"
                        onChange={handleImage}
                        disabled={isLoading}
                      />
                      {isLoading && (
                        <div className="absolute inset-0 bg-white dark:bg-gray-800 bg-opacity-75 flex items-center justify-center rounded-lg">
                          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
                        </div>
                      )}
                    </div>
                  </label>
                </form>
              </div>
            </div>
          </div>

          {/* Right Side - Results */}
          <div className="space-y-6">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8 border border-gray-200 dark:border-gray-700">
              <h2 className="text-2xl font-semibold text-gray-800 dark:text-white mb-6">
                読み取り結果
              </h2>

              {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl p-4 mb-4">
                  <div className="flex items-center">
                    <svg
                      className="w-5 h-5 text-red-500 mr-2"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clipRule="evenodd"
                      />
                    </svg>
                    <p className="text-red-800 dark:text-red-200 text-sm">
                      {error}
                    </p>
                  </div>
                </div>
              )}

              {url && (
                <div className="bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-xl p-6">
                  <div className="flex items-center mb-3">
                    <svg
                      className="w-5 h-5 text-green-500 mr-2"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      />
                    </svg>
                    <p className="text-green-800 dark:text-green-200 font-medium">
                      URL を読み取りました
                    </p>
                  </div>
                  <div className="bg-white dark:bg-gray-900 rounded-lg p-4 border border-green-200 dark:border-green-800">
                    <a
                      className="font-mono text-sm text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-200 break-all transition-colors duration-200"
                      href={url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {url}
                    </a>
                  </div>
                  <button
                    onClick={() => { void navigator.clipboard.writeText(url); }}
                    className="mt-4 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors duration-200"
                  >
                    URL をコピー
                  </button>
                </div>
              )}

              {!url && !error && !isLoading && (
                <div className="text-center py-12">
                  <svg
                    className="w-16 h-16 text-gray-400 mx-auto mb-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={1}
                      d="M12 4v1m6 11h2m-6 0h-2v4m-2 0h-2m3-4h2m-6 0h2v4m-2 0h-2m9-4h2m-6 0h-2v4m-2 0h-2"
                    />
                  </svg>
                  <p className="text-gray-500 dark:text-gray-400">
                    QR コードを選択してください
                  </p>
                </div>
              )}
            </div>
          </div>
        </main>

        {/* Footer */}
        <footer className="mt-16 text-center text-gray-600 dark:text-gray-400">
          <p className="text-sm">
            © 2024 QR Code Reader. All rights reserved.
          </p>
        </footer>
      </div>
    </div>
  );
}

export default App;
