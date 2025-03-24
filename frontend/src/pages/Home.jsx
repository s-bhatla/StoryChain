import { useState } from 'react'
import { Link, useNavigate } from "react-router-dom";

import '../App.css'

function Home() {
  const [room, setRoom] = useState(0)
  const [name, setName] = useState("")
  const [errmsg, setErrmsg] = useState("")

  async function checkUsernameAvailability(roomID, username) {

    const response = await fetch(`http://127.0.0.1:3000/check-username/${roomID}/${username}`);
    console.log(response)
    const data = await response.json();
    
    if (response.ok && data.available) {
      console.log("Username is available");
      return true;
    } else {
      console.error(data.message);
      return false;
    }
  }

  const navigate = useNavigate();

  const goToRoom = () => {
    checkUsernameAvailability(room, name).then((isAvailable) => {
      if (isAvailable) {
        console.log("You can use this username.");
        navigate(`/room/${room}`, { state: {username: name} }); // Update the path dynamically
      } else {
        console.log("Please choose another username.");
        setErrmsg("Please choose another username.")
      }
    });
    
  };

  return (
      <div className="min-h-screen w-full bg-gray-50 text-gray-900 flex flex-col items-center justify-center">
        <div className="flex flex-col items-center justify-center ">
          <h3>Enter Name</h3>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Enter your name"
        />
        </div>
        
      <div className="flex flex-col items-center justify-center">
        <h3>Enter Room Code</h3>
        <input
          type="number"
          value={room}
          onChange={(e) => setRoom(e.target.value)}
          placeholder="Enter room number"
        />
       <button onClick={goToRoom}>Send</button>
      </div>
      <div>{errmsg}</div>
      
    </div>
  )
}

export default Home
