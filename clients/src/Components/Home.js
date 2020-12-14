import React, { Component } from "react";
import "../Styles/App.css";
import api from './../Constants/APIEndpoints/APIEndpoints';
import Errors from './Errors/Errors';
import Footer from "./Footer";
import CardList from "./CardList";

class Home extends Component {
  // submitForm = async (e) => {
  //   e.preventDefault();
  //   const { email, password } = this.state;
  //   const sendData = { email, password };
  //   const response = await fetch(api.base + api.handlers.sessions, {
  //     method: "POST",
  //     body: JSON.stringify(sendData),
  //     headers: new Headers({
  //       "Content-Type": "application/json"
  //     })
  //   });
  //   if (response.status >= 300) {
  //     const error = await response.text();
  //     this.setError(error);
  //     return;
  //   }
  // }

  constructor(props) {
    super(props);

    this.state = {
      courseCode: "",
      //courseID: 0,
      results: [],
      evals: [],
      error: ""
    };
  }


  /**
   * @description setError sets the error message
   */
  setError = (error) => {
      this.setState({ error })
  }

  // GET /v1/courses?code=
  getInfo = async () => {
    const courseCode = this.state.courseCode;
    //const sendData = { courseCode };
    const response = await fetch(api.base + api.handlers.courses + '?code=' + courseCode, {
      method: "GET",
    });
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      return;
    }
    const courses = await response.json();
    if (courses.length >= 1) {  
      this.setState({
        results: [courses[0].id, courses[0].code, courses[0].title, courses[0].description], 
        //courseID: courses[0].id
      })
      this.setError("");
    } else {
      this.setState({
        results: ["no results"]
      })
    }
  }

  getEvals = async (courseID) => {
    const response = await fetch(api.base + api.handlers.courses + "/" + courseID + "/evaluations", {method: "GET"});
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      return;
    }

    const evals = await response.json();
    
    // check if there is more than one evaluation for course 
    if (evals.length > 0) {
      // get current evals saved in state
      let currEvals = this.state.evals;

      // for each evaluation:
      // -- create array containg all eval information to display
      // -- add eval to currEvals (for state)
      evals.forEach(function(e) {
        let evalInfo = [e.id, e.studentID, e.courseID, e.instructors[0]['name'], e.year, e.quarter, e.creditType, e.credits, e.workload, e.gradingTechniques, e.description, e.likedUsers.length, e.dislikedUsers.length, e.createdAt, e.editedAt]

        currEvals.push(evalInfo);
      })

      this.setState({
        evals: currEvals
      })
      this.setError("");
    }
  }

  handleInputChange = () => {
    this.setState({
      courseCode: this.search.value,
      evals: []
    }, () => {
      this.getInfo()
      
      this.getEvals(this.state.results[0])
    })
  }


  render() {
    const { error } = this.state;
    return (
      <div>
        <div>
          <div class="card mb-3 text-center">
            <div class="card-header">
              <h1>Find Course Evaluations</h1>
            </div>
            <div class="card-body">
              <form >
                <div class="searchBar">
                  <input type="text" class="btn btn-lg " id="textInput"
                    placeholder="e.g. INFO, INFO441..."
                    ref={input => this.search = input}
                    onChange={this.handleInputChange}
                  />
                  <button type="button" class="btn btn-lg">Search</button>
                </div>
              </form>
              <Errors error={error} setError={this.setError} />
            </div>
          </div>
          <div style={{ verticalAlign: "center" }}>
            <main>
              <div>
                <section>
                  <CardList classInfo={this.state.results}/>
                  <p>{this.state.evals}</p>
                </section>
              </div>
            </main>
            <Footer />
          </div>
        </div>
      </div>
    );
  }
}
export default Home;
