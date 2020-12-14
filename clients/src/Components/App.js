import React, { Component } from "react";
import "../Styles/App.css";

import Header from "./Header";
import Home from "./Home";
import Profile from "./Profile";
import { BrowserRouter as Router, Switch, Route} from "react-router-dom";

class App extends Component {
  render() {
    return (
      <Router>
        <div className="App">
          <Header />
          <Switch>
            <Route path="/" exact component={Home}/>
            <Route path="/Profile" component={Profile}/>
          </Switch>
        </div>
      </Router>
    );
  }
}
export default App;
