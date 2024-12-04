import reactLogo from "./assets/react.svg";
import "./App.css";
import Navbar from "./Components/Navbar";
import LandingPage from "./Components/LandingPage";
import React, { useState, useEffect } from "react";
import axios from "axios";
import { use } from "motion/react-client";

const getShortUrl = async (url: string) => {
  const data = {
    url: url,
  };
  console.log(url);

  try {
    const resp = await axios.post("http://localhost:3000/api/v1", data);
    console.log(resp);

    try {
      navigator.clipboard.writeText(resp.data.customShort);
      alert("Link Copied to Clipboard");
    } catch (error) {
      console.log("Error copying to clipboard", error);
    }
  } catch (error: any) {
    // Check if the error response has a rate limit exceeded message
    if (
      error.response &&
      error.response.data &&
      error.response.data.error === "Rate limit exceeded"
    ) {
      const retryAfter = error.response.data.retryAfter || "a few seconds";
      alert(`Too many requests, please try again after ${retryAfter}`);
    } else if (
      error.response &&
      error.response.data &&
      error.response.data.error === "URL cannot be empty"
    ) {
      alert("URL cannot be empty");
    } else {
      alert("An error occurred. Please try again.");
    }
    console.log("Error", error);
  }
};

const getQrCode = async (
  url: string,
  setImageSrc: React.Dispatch<React.SetStateAction<string | undefined>>,
  
) => {
  const data = {
    url: url,
  };
  try {
    const resp = await axios.post("http://localhost:3000/api/v1/qr", data, 
      // ? REMOVED THIS BELOW SINCE THE SERVER IS SENDING THE BASE64 IMAGE==================================== 
    //   {
    //   responseType: "arraybuffer", // Expect binary data (ArrayBuffer) for qrCode
    // }
    //? =================================== 
  );
    const qrCode = resp.data.qrCode;
    // ? COMMENTED BELOW SINCE NO NEED TO CONVERT .SINCE WE ARE GETTING THE IMAGE IN BASE64 FORMAT AND NOT IN []Byte ARRAY(ARRAYBUFFER) FORMAT
    //qeCode is in the []byte format so below we are converting it into blob
    // const blob = new Blob([qrCode], { type: "image/png" }); //! Convert ArrayBuffer to Blob
    // const url = URL.createObjectURL(blob); //! Create an object URL from the Blob
    // ? =================================================================================================


    const imageUrl = `data:image/png;base64,${qrCode}`; // Create Base64 image URL
    setImageSrc(imageUrl); // Set the URL as the src of the image

  } catch (error: any) {
    // Check if the error response has a rate limit exceeded message
    if (
      error.response &&
      error.response.data &&
      error.response.data.error === "Rate limit exceeded"
    ) {
      const retryAfter = error.response.data.retryAfter || "a few seconds";
      alert(`Too many requests, please try again after ${retryAfter}`);
    } else if (
      error.response &&
      error.response.data &&
      error.response.data.error === "URL cannot be empty"
    ) {
      alert("URL cannot be empty");
    } else {
      alert("An error occurred. Please try again.");
    }
    console.log("Error", error);
  }
};

function App(): JSX.Element {
  const [value, setValue] = useState<string>(""); // State for controlled input
  const [imageSrc, setImageSrc] = useState<string | undefined>(undefined);
  

  //   useEffect(() => {
  //     if (showQrCode) {
  //       // Clear the QR code image after 5 seconds
  //       const timeout = setTimeout(() => {
  //         setImageSrc(null);
  //         setShowQrCode(false);
  //       }, 5000);
  //       return () => clearTimeout(timeout);
  //     }
  // }, [showQrCode]);

  return (
    // give the min-h-screen to the parent div so that the height of the parent div is atleast the height of the screen
    <div className="min-h-screen">
      <div className=" z-50 fixed w-full">
        <Navbar />
      </div>
      {/* this padding below is very imp since it compasate the height of the above fixed navbar.and doesnt disturb the flow of the UI */}

      {/* Fixed Navbar Impacting Layout: The Navbar is set to position: fixed, which removes it from the normal document flow. As a result, it doesn't contribute to the height of the parent container, potentially causing overlapping or layout misalignment.
    Fix: Add padding to the top of the next section to account for the height of the fixed navbar. For example: see below padding is added of pt-16 */}
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
            onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
              setValue(e.target.value)
            }
          />
        </div>
        <div className="flex space-x-10 flex-row">
          <button
            className="border-2 px-5 border-purple-400 font-bold  rounded-lg p-2 hover:bg-gradient-to-tr from-amber-800 via-blue-800 to-black"
            // ! THERE IS DIFFERENCE IN THE BELOW TWO "onClick" FUNCTION
            // onClick={() => {
            //   getShortUrl(value), console.log("hello hardik");
            // }}
            onClick={() => {
              getShortUrl(value);
              console.log("hello hardik");
            }}
          >
            Get Link
          </button>
          <button
            className="border-2 border-purple-400 font-bold  rounded-lg p-2 hover:bg-gradient-to-tr from-amber-800 via-blue-800 to-black"
            onClick={() => {
              getQrCode(value, setImageSrc );
              console.log("hello hardik");
            }}
          >
            Get QR Code
          </button>
        </div>
        {  imageSrc && (
          <div className="mt-6">
             

            <img
              src={imageSrc}
              alt="Generated QR Code"
              style={{ width: "200px", height: "200px" }}
            />
          </div>
        )}
        <div className="my-6">
          <LandingPage />
        </div>
      </div>
    </div>
  );
}

export default App;

/*The reason I used `min-h-screen` in both places is to ensure that the background and layout cover the entire height of the viewport, even with a fixed navbar. Let me break it down for clarity:

### 1. **First `min-h-screen`:**
```jsx
<div className="min-h-screen">
```
This ensures that the outer container (the parent div of the entire app) spans at least the full height of the viewport. Without this, the background or the content might not fully cover the screen height if there's no content filling it. It is a safeguard to ensure that the entire page will cover the viewport height, even if the content is minimal.

### 2. **Second `min-h-screen` inside the content container:**
```jsx
<div className="pt-16 min-h-screen flex flex-col items-center justify-center bg-gradient-to-tr from-amber-800 via-blue-800 to-black text-white">
```
In this case, the second `min-h-screen` ensures that the content section (where the heading, input, and other components are) also takes up at least the full height of the viewport. Here's why:

- **Fixed Navbar (`pt-16`)**: Since the navbar is `fixed` at the top, it does not affect the height of the content section below it. However, without the `min-h-screen` here, the content section might only take the height of its content and not fill the remaining space below the navbar.
  
- **Vertical Centering (`flex` + `justify-center`)**: This ensures that the content within this section is vertically centered. But, without `min-h-screen`, if the content is smaller than the screen, the section won't fill the entire height, which could make it look unbalanced.

### Why `min-h-screen` in both:
1. **The Outer Container (`min-h-screen`)** ensures the whole page takes up at least the full viewport height.
2. **The Content Container (`min-h-screen`)** ensures that the section beneath the navbar takes up the remaining screen space and fills the whole viewport vertically, avoiding any potential layout issues caused by the navbar taking up space.

If you didn't use `min-h-screen` on the content container, and the content wasn't tall enough, you'd see a large gap between the content and the bottom of the screen.

In summary, both `min-h-screen` are used to ensure full-screen coverage and correct vertical alignment, especially with a fixed navbar that doesn't affect the normal flow of the page. */
