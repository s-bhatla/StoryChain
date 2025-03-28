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


//CHANGE Chatroom:
// 1) Game start, read broadcasts for storylines, next prompts, submit prompts
// 1.5) Add time functionality (no backend security for now, just make it)
// 2) Redirect to results page after game has ended