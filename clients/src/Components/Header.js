import React, { Component } from "react";
import "../Styles/App.css";

class Header extends Component {
  render() {
    return (
      <header>
        <nav>
          <div className="headerButton">
            <h1>RateTheCourse</h1>
            <div className="navTitle">
              <img src={require("../img/home.png")} aria-hidden="true" alt="Logo"/>
              <img src={require("../img/user.png")} aria-hidden="true" alt="Logo"/>
              <img src={require("../img/edit.png")} aria-hidden="true" alt="Logo"/>   
            </div>
          </div>
        </nav>
      </header>
    );
  }
}

export default Header;
