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
      results: [],
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
      // var currCourses = this.state.results;
      // courses.forEach(function(c) {
      //   var cToAdd = [c.id, c.code, c.title, c.description];
      //   currCourses.push(cToAdd);
      // })
      this.setState({
        results: [courses[0].id, courses[0].code, courses[0].title, courses[0].description], 
        //currCourses: currCourses
      })
      this.setError("");
    } else {
      this.setState({
        results: ["no results"]
      })
    }
  }

  handleInputChange = () => {
    this.setState({
      courseCode: this.search.value
    }, () => {
      this.getInfo()
      
      //this.getEvals(this.state.results[0])
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
                  <CardList classInfo={this.state.results} />
                  
                  {/* {this.state.results.map(data => {
                      return (
                          <CardList classInfo={this.state.results}/>
                      );
                  })} */}
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
