import ReactDOM from "react-dom/client";
import { BrowserRouter, Routes, Route } from "react-router-dom";

import Home from './pages/Home.jsx';
import Chatroom from './pages/Chatroom.jsx';

import './App.css'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/room/:roomID" element={<Chatroom />} />
        <Route path="/results/:roomID" component={Results} />
      </Routes>
    </BrowserRouter>
  );
}


export default App
