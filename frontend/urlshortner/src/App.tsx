import reactLogo from "./assets/react.svg";
import "./App.css";
import Navbar from "./Components/Navbar";

function App(): JSX.Element {
  return (
	  <div className=" h-screen overflow-hidden">
	  <Navbar />
      <div className="min-h-screen flex flex-col  items-center justify-center bg-gradient-to-tr from-amber-800 via-blue-800 to-black text-white">
        <main className=" h-full">
          <h1 className="text-7xl font-bold">Free Url Shortner!</h1>

		  <p className="flex flex-row justify-center text-xl font-bold text-black">GET YOUR LINK SHORTNED</p>
        </main>
      </div>
    </div>
  );
}

export default App;
