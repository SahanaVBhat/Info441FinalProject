import React, { Component } from "react";
import "../Styles/App.css";

class Header extends Component {
  render() {
    return (
      <header>
        <button type="button" className="headerButton">
          <img src={require("../img/husky.png")} aria-hidden="true" alt="Logo"/> 
        </button>
        <nav>
          <p >Post</p>
          <p >Profile</p>
        </nav>
      </header>
    );
  }
}

export default Header;
