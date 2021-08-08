// import React from 'react';
import { React } from '../react';
import MainPage from './MainPage';
import firebase from 'firebase/app';
import 'firebase/firestore';
import 'firebase/auth';
import { useAuthState } from 'react-firebase-hooks/auth';
import axios from 'axios';

let auth
axios.get("/api/firebase")
  .then((response) => {
    firebase.initializeApp(response.data);
    auth = firebase.auth()
  })

function LoginPage() {
  const [userAuth] = useAuthState(auth);

  return (
    <div className="LoginPage">
      <h1>Security Cam</h1>
      {userAuth ? <MainPage userAuth={userAuth} auth={auth} axios={axios} /> : <SignIn />}
    </div>
  );

  function SignIn() {
    const signInWithGoogle = () => {
      const provider = new firebase.auth.GoogleAuthProvider();
      auth.signInWithPopup(provider);
    }

    return (
      <button onClick={signInWithGoogle}>Sign in with Google</button>
    )
  }
}

export default LoginPage;