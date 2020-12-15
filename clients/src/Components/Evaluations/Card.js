import React from "react";
import api from '../../Constants/APIEndpoints/APIEndpoints';
import { withRouter } from "react-router";
import { Link } from 'react-router-dom'


class Card extends React.Component {
  setError = (error) => {
    this.setState({ error })
  }

  constructor(props) {
    super(props);

    this.state = {
      evalID: 0,
      authToken: null,
      newDescr: ""
    };
  }

  goToCarddetails = (cardId) => {
      localStorage.setItem("selectedCard", cardId);
      this.props.history.push(cardId);
  }

  render() {
      return ( 
          //<div className="card" onClick = {()=>this.goToCarddetails(this.props.classInfo.id.toString())}>
          <div className="card">
          <div className="card-body">
            <h3 className="card-title">Evaluation #{this.props.value+1}</h3>
            <h5>Credit Type: {this.props.classInfo.creditType}</h5>
            <h5>Instructors: {this.props.classInfo.instructors.map((instructor) => instructor.name)}</h5>
            <h5>Quarter: {this.props.classInfo.quarter}</h5>
            <h5>Credit Type: {this.props.classInfo.creditType}</h5>
            <h5>Workload: {this.props.classInfo.workload}</h5>
            <h5>Grading fairness: {this.props.classInfo.gradingTechniques}</h5>
            <h5>Comments: {this.props.classInfo.description}</h5>
          </div>
          <div className="card-buttons">
            <button type="button" class="btn btn-primary" onClick={this.unhide}>Edit</button>
            <button type="button" class="btn btn-primary" onClick={this.handleRemove}>Remove</button>
          </div>
          <div className="card-update-textbox">
            <input type="text" placeholder="Description..." id={"input-" + this.props.classInfo.id} ref={input => this.updt = input} hidden></input>
            <button type="button" class="btn btn-primary" id={"submit-" + this.props.classInfo.id} onClick={this.handleEdit} hidden>Submit</button>
          </div>
        </div>
      )
  }

  unhide = () => {
    this.setState({
      evalID: this.props.classInfo.id
    }, () => {
      document.getElementById("input-" + this.state.evalID).hidden = false;
      document.getElementById("submit-" + this.state.evalID).hidden = false;
    })
  }

  handleEdit = () => {
    this.setState({
      evalID: this.props.classInfo.id,
      authToken: localStorage.getItem("Authorization") || null,
      newDescr: this.updt.value
    }, () => {
      this.patchEval()
    })
  }

  handleRemove = () => {
    this.setState({
      evalID: this.props.classInfo.id,
      authToken: localStorage.getItem("Authorization") || null
    }, () => {
      this.removeEval()
    })
  }

  patchEval = async() => {
    console.log("in patch eval")
    if (!this.state.authToken) {
      return;
    }

    // PATCH /v1/evaluations/:evaluationID
    const response = await fetch(api.base + api.handlers.evaluations + "/" + this.state.evalID, 
      { 
        method: "PATCH", 
        headers: new Headers({
          "Authorization": this.state.authToken,
          "Content-type": "application/json"
        }),
        body: JSON.stringify({description: this.state.newDescr})
      });

    // hide submit and text input for edit 
    document.getElementById("input-" + this.state.evalID).hidden = true;
    document.getElementById("submit-" + this.state.evalID).hidden = true;
    
    // check response status for errors
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      alert(error)
    }

    console.log(response.status)
  }

  removeEval = async() => {
    if (!this.state.authToken) {
      return;
    }

    // DELETE /v1/evaluations/:evaluationID
    const response = await fetch(api.base + api.handlers.evaluations + "/" + this.state.evalID, 
      { 
        method: "DELETE", 
        headers: new Headers({
        "Authorization": this.state.authToken
      })
    });

    // check response for errors
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      alert(error);
      return;
    }
  }
}

  export default withRouter(Card);