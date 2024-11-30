import React from "react";
import { motion } from "framer-motion";

export default function Navbar(): JSX.Element {
  return (
    <nav className="  w-full bg-gradient-to-br from-gray-900 via-gray-800 to-black z-50">
      <motion.div
        className="self-center font-mono text-2xl font-semibold whitespace-nowrap text-white text-center py-3"
        animate={{ x: ["100%", "-100%"] }} // Move from right to left
        transition={{
          duration: 20,
          repeat: Infinity,
          ease: "linear",
        }}
      >
        WELCOME TO FREE URL SHORTENER!
      </motion.div>
    </nav>
  );
}
