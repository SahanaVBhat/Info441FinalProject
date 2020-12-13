import React, { Component } from "react";
import "../Styles/App.css";
import Footer from "./Footer";
// import CardList from "./CardList";
import Header from "./Header";
import api from './../Constants/APIEndpoints/APIEndpoints';


class Dashboard extends Component {
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
      results: []
    };
  }

  getInfo = async () => {
    const courseCode = this.state.courseCode;
    //const sendData = { courseCode };
    const response = await fetch(api.base + api.handlers.courses + '/?code=' + courseCode, {
      method: "GET",
    });
    if (response.status >= 300) {
      const error = await response.text();
      this.setError(error);
      return;
    }
    const courses = await response.json();
    this.setState({
      results: courses
    })
  }

  handleInputChange = () => {
    this.setState({
      courseCode: this.search.value
    }, () => {
      this.getInfo()
      // if (this.state.courseCode && this.state.courseCode.length > 1) {
      //   if (this.state.courseCode.length % 2 === 0) {
      //     this.getInfo()
      //   }
      // }
    })
  }


  render() {
    return (
      <div>
        <div>
          <Header />
        </div>
        <div>
          <div class="card mb-3 text-center">
            <div class="card-header">
              <h1>Find Course Evaluations</h1>
            </div>
            <div class="card-body">
              <form>
                <input type="text" class="btn-lg" id="textInput"
                  placeholder="e.g. INFO, INFO441..."
                  ref={input => this.search = input}
                  onChange={this.handleInputChange}
                />
                <p>{this.state.query}</p>
                <button type="button" class="btn personalButton btn-lg">Search</button>

              </form>
            </div>
          </div>

          <div style={{ verticalAlign: "center" }}>
            <main>
              <div>
                <section>
                  {/* <CardList/> */}
                  <p>{this.state.results}</p>
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
export default Dashboard;
