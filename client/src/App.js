import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css'
import React, { useState, useEffect } from 'react';
import LoginPage from './components/LoginPage';
import axios from 'axios';
import io from 'socket.io-client';

function App() {
  const [testResult, setTestResult] = useState();
  
  
  
  useEffect(() => {
    axios.get("/api/test")
    .then((response) => {
      setTestResult(response.data.message)
    })
    const socket = io("localhost:8080/api/socket.io", {transports: ["websocket"]});
    // const socket = io('http://localhost:8080/api/socket.io/')
    // const socket = io({
    //   path: "/api/socket.io",
    //   port: "8080",
    //   transports: ["websocket"]
    // });
    console.log("socket before")
    console.log(socket)
    console.log("socket after")
    
  }, [])

  return (
    <div className="App">
      {testResult ? <LoginPage /> : <h1>Loading</h1>}
    </div>
  );
}

export default App;
