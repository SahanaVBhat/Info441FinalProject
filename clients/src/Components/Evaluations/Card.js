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
          <div className="card-buttons" >
            <button type="button" class="btn btn-primary" onClick={this.handleEdit}>Edit</button>
            <button type="button" class="btn btn-primary" onClick={this.handleRemove}>Remove</button>
          </div>
        </div>
      )
  }

  handleEdit = () => {
    this.setState({
      evalID: this.props.classInfo.id,
      authToken: localStorage.getItem("Authorization") || null
    }, () => {
      // create text box for user to input new description 
      // var updateDescr = document.createElement("INPUT");
      // updateDescr.setAttribute("type", "text");
      // updateDescr.setAttribute("placeholder", "Type updated description here...")
      // document.getElementById("buttons-" + this.state.evalID).appendChild(updateDescr);
      // set description in setState to text box value 

      // trigger patchEval()
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

    // const response = await fetch(api.base + api.handlers.evaluations + "/" + this.state.evalID, { method: "PATCH" });
    // if (response.status >= 300) {
    //   const error = await response.text();
    //   this.setError(error);
    //   return;
    // }
  }

  removeEval = async() => {
    if (!this.state.authToken) {
      return;
    }

    const response = await fetch(api.base + api.handlers.evaluations + "/" + this.state.evalID, 
      { 
      method: "DELETE", 
      headers: new Headers({
      "Authorization": this.state.authToken
      })
    });

    console.log(response.status)
  }
}

  export default withRouter(Card);