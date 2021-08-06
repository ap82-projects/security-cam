import logo from './logo.svg';
import './App.css';
import React, { useState, useEffect } from "react";
import axios from "axios";

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
  }, [])

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.js</code> and save to reload.
        </p>
        <p>Your current test is {testResult}</p>
      </header>
    </div>
  );
}

export default App;
