import reactLogo from "./assets/react.svg";
import "./App.css";
import Navbar from "./Components/Navbar";
import LandingPage from "./Components/LandingPage";
import React from "react";

function App(): JSX.Element {
  const [value, setValue] = React.useState(""); // State for controlled input

  return (
    <div className="min-h-screen">
      <div className=" z-50 fixed w-full">
        <Navbar />
      </div>
      <div className="pt-16 min-h-screen flex flex-col items-center justify-center bg-gradient-to-tr from-amber-800 via-blue-800 to-black text-white">
        <h1 className="text-7xl font-bold">Best Free URL Shortener!</h1>
        <p className="text-xl font-bold my-5 text-yellow-400 z-10 animate-pulse">
          GET YOUR LINK SHORTENED
        </p>
        <div className="flex flex-row space-x-3 justify-center w-full mb-6">
          <label htmlFor="url" className="text-2xl font-bold">
            Enter your URL:
          </label>
          <input
            id="url"
            type="text"
            placeholder="Enter your URL"
            className="border border-gray-300 p-2 rounded-lg w-3/6 text-gray-800 focus:outline-none focus:ring-2 focus:ring-blue-500"
            value={value}
            onChange={(e) => setValue(e.target.value)}
          />
        </div>
        <div className="my-6">
          <LandingPage />
        </div>
      </div>
    </div>
  );
}

export default App;


