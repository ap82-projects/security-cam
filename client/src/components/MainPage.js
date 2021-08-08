import './MainPage.css'
import React, { useState, useEffect } from 'react';
import Incidents from './Incidents';
import SecurityCam from './SecurityCam';
import VideoChat from './VideoChat';

function MainPage(props) {
  const { userAuth, auth, axios } = props;
  const [userGoogleId, setUserGoogleId] = useState(userAuth.providerData[0].uid);
  const [userGoogleName, setUserGoogleName] = useState(userAuth.displayName);
  const [userGoogleEmail, setUserGoogleEmail] = useState(userAuth.email);
  const [userGooglePhone, setUserGooglePhone] = useState(userAuth.phoneNumber);
  const [userDocumentId, setUserDocumentId] = useState("")
  const [user, setUser] = useState({});
  const [asSecurityCam, setAsSecurityCam] = useState(false);
  const [watchSecurityCam, setWatchSecurityCam] = useState(false);

  useEffect(async () => {
    const existingUserDoc = await getUserDocId(userGoogleId);
    if (existingUserDoc && existingUserDoc.id) {
      setUserDocumentId(existingUserDoc.id)
      const userData = await getUserData(existingUserDoc.id);
      setUser(userData)
      /////////////////////////////////////////
      // Possible error handling functionality
      // if (userData) {
      //   setUser(userData);
      // } else {
      //   // ERROR
      //   auth.signOut()
      // }
      /////////////////////////////////////////
    } else {
      // User doesn't exist in database
      const newUserDoc = await addUser();
      setUserDocumentId(newUserDoc.id)
      const newUserData = await getUserData(newUserDoc.id)
      setUser(newUserData);
      /////////////////////////////////////////
      // Possible error handling functionality
      // if (newUserDoc && newUserDoc.id) {
      //   setUserDocumentId(newUserDoc.id)
      //   setUser(await getUserData(newUserDoc.id));
      // } else {
      //   // ERROR, user not added
      //   auth.signOut();
      // }
      /////////////////////////////////////////
    }
  }, []);

  useEffect(async () => {
    const updatedUserData = await getUserData(userDocumentId);
    setUser(updatedUserData);
  }, [asSecurityCam])
  
  return (
    <div className="MainPage">
      <div>
        <button variant="danger" onClick={() => auth.signOut()}>Sign Out</button>
        {/* DISABLED UNTIL IMPLEMENTATION CLEANED UP
        <button onClick={watchSecurityCam ? cutSecurityFeed : viewSecurityFeed}>
          {watchSecurityCam ? 'Cut Security Feed' : 'View Security Feed'}
        </button> */}
        <button onClick={() => setAsSecurityCam(!asSecurityCam)}>
          {asSecurityCam ? 'Stop Monitoring' : 'Security Cam'}
        </button>
      </div>
      <div>
        {asSecurityCam
          ? <SecurityCam addIncident={addIncident} />
          : watchSecurityCam
            ? <VideoChat guestName={user.Name} guestRoom={userDocumentId} />
            : <Incidents
          user={user}
          setUser={setUser}
          userDocumentId={userDocumentId}
          getUserData={getUserData}
          axios={axios}
          />
        }
      </div>
    </div>
  )
  
  async function getUserDocId(googleId) {
    const response = await axios.get(`/api/user/google?id=${googleId}`);
    return response.data;
  }

  async function getUserData(docID) {
    const response = await axios.get(`/api/user?id=${docID}`);
    return response.data;
  }

  async function addUser() {
    const response = await axios.post(`/api/user`, {
      'email': userGoogleEmail,
      'googleid': userGoogleId,
      'incidents': [],
      'name': userGoogleName,
      'phone': userGooglePhone,
      'watching': false
    });
    return response.data;
  }

  async function viewSecurityFeed() {
    setWatchSecurityCam(true);
    const response = await axios.put(`/api/user/watching?id=${userDocumentId}`, {
      'watching': true
    });
  }

  async function cutSecurityFeed() {
    setWatchSecurityCam(false);
    const response = await axios.put(`/api/user/watching?id=${userDocumentId}`, {
      'watching': false
    });
  }

  async function addIncident(img) {
    const response = await axios.put(`/api/user/incident?id=${userDocumentId}`, {
      'image': img,
      'time': String(Date.now())
    });
  }

  async function testDeleteIncident() {
    const time = 'this time';
    const response = await axios.delete(`/api/user/incident?id=${userDocumentId}&time=${encodeURIComponent(time)}`);
  }
}

export default MainPage;
