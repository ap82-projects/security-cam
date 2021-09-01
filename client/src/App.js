import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css'
import React, { useState, useEffect } from 'react';
import LoginPage from './components/LoginPage';
import axios from 'axios';
import io from 'socket.io-client';

function App() {
  const [testResult, setTestResult] = useState();
  const [socket, setSocket] = useState();
  // const socket = io("localhost:8080/socket.io", {transports: ["websocket"]});
  // listen for messages
  // socket.on('message', function(message) {

  //   console.log('new message');
  //   console.log(message);
  // });

  // socket.on('disconnect', () => {
  //   console.log('disconnected')
  // })

  // socket.on('connect', function () {

  //   console.log('socket connected');

  //   //send something
  //   socket.emit('send', {text: "connected"}, function(result) {

  //     console.log('sended successfully');
  //     console.log(result);
  //   });
  // });
  // console.log("socket before")
  // console.log(socket)
  // console.log("socket after")



  useEffect(() => {
    axios.get("/api/test")
      .then((response) => {
        setTestResult(response.data.message)
      })


    // const socket = io();
    // const socket = io("localhost:3811/socket.io/", {transports: ["websocket"]});
    // const socket = io('http://localhost:8080/api/socket.io/')
    // const socket = io({
    //   path: "/api/socket.io",
    //   port: "8080",
    //   transports: ["websocket"]
    // });


    const socket = io("ws://localhost:8080/socket.io", { transports: ["websocket"] });
    // listen for messages
    socket.on('message', function (message) {
      console.log('new message');
      console.log(message);
    });

    socket.on('disconnect', () => {
      console.log('disconnected')
    })

    socket.on('connect', function () {
      console.log('socket connected');
      //send something
      socket.emit('send', { text: "connected" }, function (result) {

        console.log('sended successfully');
        console.log(result);
      });
    });
    
    console.log("socket before")
    console.log(socket)
    console.log("socket after")
    setSocket(socket)

  }, [])

  return (
    <div className="App">
      {testResult ? <LoginPage /> : <h1>Loading</h1>}
    </div>
  );
}

export default App;
