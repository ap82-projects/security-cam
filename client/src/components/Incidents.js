import './Incidents.css';
import React from 'react';

function Incidents(props) {
  const {
    user,
    setUser,
    userDocumentId,
    getUserData,
    axios
  } = props;

  return (
    <div className='Incidents'>
      <h3>Incidents</h3>
      {parceIncidents()}
    </div>
  )

  function parceIncidents() {
    if (user.Incidents) {
      return user.Incidents.map(incident => (
        <div className="card" key={incident.Time}>
          <div className="row">
            <p className="card-text">{(new Date(Number(incident.Time))).toString()}</p>
          </div>
          <div className="row">
            <div className="col">
              <img className='incident-image' src={incident.Image}></img>
            </div>
            <div className="col">
              <button type="button" className="btn btn-danger" id={incident.Time} onClick={deleteIncident}>Delete Incident</button>
            </col>
          </div>
        </div>
      ));
    } else {
      return <div></div>
    }
  }

  async function deleteIncident(e) {
    const response = await axios.delete(`/api/user/incident?id=${userDocumentId}&time=${e.target.id}`);
    const updatedUserData = await getUserData(userDocumentId);
    setUser(updatedUserData)
  }
}

export default Incidents;
