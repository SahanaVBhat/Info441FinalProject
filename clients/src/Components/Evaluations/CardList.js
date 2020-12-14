import React, { Component } from "react";
import "../../Styles/App.css";
import api from '../../Constants/APIEndpoints/APIEndpoints';
import Card from './Card';
import Profile from '../Profile';
import Button from '@material-ui/core/Button';
import {withRouter} from 'react-router-dom';
import { Link } from 'react-router-dom'


class CardList extends Component {
  setError = (error) => {
    this.setState({ error })
  }

  constructor(props) {
    super(props);

    this.state = {
      evals: [],
      show: false
    };
  }

  nextPath(path) {
    this.props.history.push(path);
  }

  render() {

    if (this.props.classInfo.length != 0) {
      // if (this.state.evals[0] != null) {
      //     console.log(this.state.evals[0][3])
      //     var instructor = document.createElement("P");
      //     var instrText = document.createTextNode("Instructor(s): " + this.state.evals[0][3]);
      //     instructor.appendChild(instrText);

      //     var descr = document.createElement("P");
      //     var descrText = document.createTextNode("Description: " + this.state.evals[0][10]);
      //     descr.appendChild(descrText);

      //     document.getElementById("hi").appendChild(document.createElement("HR"));
      //     document.getElementById("hi").appendChild(instructor);
      //     document.getElementById("hi").appendChild(descr);
      // }
      return (
        <div class="card text-center">
          <div class="card-header">
            {this.props.classInfo[1]}
          </div>
          <div class="card-body">
            <h5 class="card-title">{this.props.classInfo[2]}</h5>
            <p class="card-text">{this.props.classInfo[3]}</p>
            <a href="#" class="btn btn-primary">
            <Link to={{
                  pathname: '/Profile',
                  state: {
                    courseID: this.props.classInfo[0]
                  }}}>Add Evaluation</Link>

            </a>
            <a href="#" class="btn btn-primary" onClick={this.handleEvalClick}>Read Evaluations</a>
          </div>
          <div class="eval-body" id="hi">
            {this.state.show && this.state.evals.map((data, index) => {
              return (<Card classInfo={data} value={index}/>);
            })}
          </div>
        </div>
      );
    } else {
      return (
        <div></div>
      );
    }
  }

  handleEvalClick = () => {
    this.setState({
      evals: []
    }, () => {
      this.getEvals(this.props.classInfo[0])
      //this.getEvals()
    })
  }

  getEvals = async (courseID) => {
    const response = await fetch(api.base + api.handlers.courses + "/" + courseID + "/evaluations", { method: "GET" });
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      return;
    }

    const evals = await response.json();

    // check if there is more than one evaluation for course 
    if (evals.length > 0) {
      // get current evals saved in state
      //let currEvals = this.state.evals;

      // for each evaluation:
      // -- create array containg all eval information to display
      // -- add eval to currEvals (for state)

      let evalsMap = evals.map((data, key) => {
        var oneEval = {};
        oneEval.id = data.id;
        oneEval.studentID = data.studentID;
        oneEval.courseID = data.courseID;
        oneEval.instructors = data.instructors;
        oneEval.year = data.year;
        oneEval.quarter = data.quarter;
        oneEval.creditType = data.creditType;
        oneEval.credits = data.credits;
        oneEval.workload = data.workload;
        oneEval.gradingTechniques = data.gradingTechniques;
        oneEval.description = data.description;

        return oneEval;
      });

      console.log(evalsMap)

      // evals.forEach(function(e) {
      //     let evalInfo = [e.id, e.studentID, e.courseID, e.instructors[0]['name'], e.year, e.quarter, e.creditType, e.credits, e.workload, e.gradingTechniques, e.description, e.likedUsers.length, e.dislikedUsers.length, e.createdAt, e.editedAt]

      //     currEvals.push(evalInfo);
      // })

      // add evaluations to page 
      if (evalsMap.length >= 1) {
        this.setState({
          evals: [ ...this.state.evals,...evalsMap ],
          show: true
        })
        this.setError("");
        // instructor 
        // var instructor = document.createElement("P");
        // var instrText = document.createTextNode("Instructor(s): " + this.state.evals[0][3]);
        // instructor.appendChild(instrText);


        // // created at 
        // var created = document.createElement("P");
        // var createdText = document.createTextNode("Created At: " + this.state.evals[0][11]);
        // created.appendChild(createdText);

        // // created at 
        // var created = document.createElement("P");
        // var createdText = document.createTextNode("Created At: " + this.state.evals[0][13]);
        // created.appendChild(createdText);

        // document.getElementById("hi").appendChild(document.createElement("HR"));
        // document.getElementById("hi").appendChild(instructor);
        // document.getElementById("hi").appendChild(descr);
        // document.getElementById("hi").appendChild(created);

      }
    } else {
      this.setState({
        evals: [],
        show: false
      })
    }
  }
}
// class CardList extends Component {
//   render() {
//     return <Card id={this.props.classInfo[0]} course={this.props.classInfo[1]} name={this.props.classInfo[2]} d={this.props.classInfo[3]}/>;
//   }
// }
export default CardList;