import './App.css';
import React from 'react';
import MainPage from './components/MainPage';
import firebase from 'firebase/app';
import 'firebase/firestore';
import 'firebase/auth';
import { useAuthState } from 'react-firebase-hooks/auth';

firebase.initializeApp({
  apiKey: process.env.REACT_APP_API_KEY,
  authDomain: process.env.REACT_APP_AUTH_DOMAIN,
  projectId: process.env.REACT_APP_PROJECT_ID,
  storageBucket: process.env.REACT_APP_STORAGE_BUCKET,
  messagingSenderId: process.env.REACT_APP_MESSAGING_SENDER_ID,
  appId: process.env.REACT_APP_APP_ID
});

const auth = firebase.auth();

function App() {
  const [userAuth] = useAuthState(auth);
  console.log("in app.js")
  return (
    <div className="App">
      <h1>Security Cam</h1>
      {userAuth ? <MainPage userAuth={userAuth} auth={auth} /> : <SignIn />}
    </div>
  );
}

function SignIn() {
  const signInWithGoogle = () => {
    const provider = new firebase.auth.GoogleAuthProvider();
    auth.signInWithPopup(provider);
  }

  

  return (
    <button onClick={signInWithGoogle}>Sign in with Google</button>
  )
}

export default App;
