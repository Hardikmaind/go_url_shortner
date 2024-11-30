import React from "react";

export default function Navbar(): JSX.Element {
  return (
    <>
      <nav className="bg-white border-gray-200 bg-gradient-to-br from-gray-900 via-gray-800 to-black">
        <div className="max-w-screen-xl flex items-center justify-center mx-auto p-2">
          <a
            href="https://flowbite.com/"
            className="flex flex-row items-center space-x-3 rtl:space-x-reverse"
          >
           
            <span className="self-center font-mono text-2xl font-semibold whitespace-nowrap dark:text-white animate-pulse ">
              WELCOME TO FREE URL SHORTNER!
            
            </span>
          </a>
        </div>
      </nav>
    </>
  );
}
