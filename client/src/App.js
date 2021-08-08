import './App.css';
import React, { useState, useEffect } from 'react';
// import { React } from './react';
// import { useState, useEffect } from React;
import LoginPage from './components/LoginPage';
import axios from 'axios';

function App() {
  const [testResult, setTestResult] = useState();

  useEffect(() => {
    console.log("in App")
    axios.get("/api/test")
      .then((response) => {
        console.log("response.data.message")
        console.log(response.data.message)
        setTestResult(response.data.message)
      })
    // axios.get("/api/privatetest")
    //   .then((response) => {
    //     console.log("private response.data.messge")
    //     console.log(response.data.message)
    //   })
    console.log("afterAxios")
  }, [])

  return (
    <div className="App">
      {testResult ? <LoginPage /> : <h1>Loading</h1>}
    </div>
  );
}

// function SignIn() {
//   const signInWithGoogle = () => {
//     const provider = new firebase.auth.GoogleAuthProvider();
//     auth.signInWithPopup(provider);
//   }

//   return (
//     <button onClick={signInWithGoogle}>Sign in with Google</button>
//   )
// }
// export default firebaseConfig;
export default App;
