import React, { Component } from "react";
import "../Styles/App.css";
import {Link} from "react-router-dom";

class Header extends Component {
  render() {
    return (
      <header>
        <nav>
          <div className="title">
            <h1 >RateTheCourse</h1>
          </div>
          <div className="navTitle">
            <Link to="/">
              <img src={require("../img/home.png")} aria-hidden="true" alt="Logo" className="imgSize"/> 
            </Link>
            <Link to="/Profile">
              <img src={require("../img/user.png")} aria-hidden="true" alt="Logo" className="imgSize"/>
            </Link> 
          </div> 
        </nav>
      </header>
    );
  }
}

export default Header;
