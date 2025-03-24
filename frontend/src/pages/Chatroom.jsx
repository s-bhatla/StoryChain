import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { useLocation } from 'react-router-dom';



const Chatroom = () => {
  const { roomID }  = useParams();
  const [socket, setSocket] = useState(null);
  const [messages, setMessages] = useState([]);
  const [message, setMessage] = useState("");
  const location = useLocation()
  const { username } = location.state || {};

  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:3000/ws/${roomID}`);
    setSocket(ws);

    ws.onopen = () => {
      console.log("USername sui ", username)
      ws.send(JSON.stringify(username));
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      const receivedMessage = {
        text: JSON.parse(event.data),
        sendrec: "rec", //flag: message is sent or recieved?
        timestamp: new Date().toISOString
      }

      setMessages((prev) => [...prev, receivedMessage]);
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
    };

    return () => ws.close();
  }, [roomID]);

  useEffect(() => {
    console.log("Printing messages array, size=", messages.length)
    messages.map((msg, index) => (
      
      console.log(msg)
    ))
  }, [message]);

  const sendMessage = () => {
    if (socket && message.trim()) {
      socket.send(JSON.stringify(message));
      const sentMessage = {
        text: message,
        sendrec: "send",
        timestamp: new Date().toISOString
      }
      setMessages((prev) => [...prev, sentMessage]);
      setMessage(""); // Clear the input
    }
  };

  return (
    <div className="min-h-screen w-full bg-gray-20 text-gray-900 flex-col items-center justify-center ps-[2rem]">
      <h1 className="w-fit">Chatroom: {roomID}</h1>
      <div className="messages w-1/2">
      <>
        {messages.map((msg, index) => (
          <div className={msg.sendrec} key={index}><span className="w-fit bg-blue-500 text-white ps-1 pe-1 rounded-lg shadow-lg ">{msg.text}</span></div>
        ))}
        </>
      </div>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Type a message"
      />
      <button onClick={sendMessage}>Send</button>
    </div>
  );
};

export default Chatroom;
