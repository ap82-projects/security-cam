import logo from './logo.svg';
import './App.css';
import React, { useState, useEffect } from "react";
import axios from "axios";
import env from "@beam-australia/react-env";

function App() {
  const [ testResult, setTestResult ] = useState("pending");

  useEffect(() => {
    console.log("before axios call")
    axios.get("api/test")
    .then((response) => {
      console.log("response");
      console.log(response);
      setTestResult(response.data.message);
    })
    console.log("before process.env")
    console.log(process.env)
    console.log("after process.env")
    console.log("before env")
    console.log(env("REACT_APP_MESSAGE"));
    console.log("after env")
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <p>{process.env.REACT_APP_MESSAGE} {testResult} process</p>
      </header>
    </div>
  );
}

export default App;
