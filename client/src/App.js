import './App.css';
import React, { useState, useEffect } from 'react';
import LoginPage from './components/LoginPage';
import axios from 'axios';

function App() {
  const [testResult, setTestResult] = useState();

  useEffect(() => {
    axios.get("/api/test")
      .then((response) => {
        setTestResult(response.data.message)
      })
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
