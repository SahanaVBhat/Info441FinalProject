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
              return (
                <Card classInfo={data} value={index}/>
              );
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

  // on click of button 'Read Evaluations' call getEvals() passing course ID from classInfo
  handleEvalClick = () => {
    this.setState({
      evals: []
    }, () => {
      this.getEvals(this.props.classInfo[0])
    })
  }

  // call 'GET /v1/courses/:courseID/evaluations` to get evaluations for specific course 
  getEvals = async (courseID) => {
    const response = await fetch(api.base + api.handlers.courses + "/" + courseID + "/evaluations", { method: "GET" });
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      return;
    }

    const evals = await response.json();

    // check if there is at least one evaluation for course 
    if (evals.length > 0) {

      // for each evaluation create array containg all eval information to display
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

      // add evaluations to page 
      if (evalsMap.length >= 1) {
        this.setState({
          evals: [ ...this.state.evals,...evalsMap ],
          show: true
        })
        this.setError("");
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