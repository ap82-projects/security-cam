import './Incidents.css';
import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import Button from 'react-bootstrap/Button';
import Card from 'react-bootstrap/Card';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';

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
        <Card key={incident.Time}>
          <Row>
            <Card.Text>{(new Date(Number(incident.Time))).toString()}</Card.Text>
          </Row>
          <Row>
            <Col>
              <img className='incident-image' src={incident.Image}></img>
            </Col>
            <Col>
              <Button variant='danger' id={incident.Time} onClick={deleteIncident}>Delete Incident</Button>
            </Col>
          </Row>
        </Card>
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
