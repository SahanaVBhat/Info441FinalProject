import React, { Component } from "react";
import "../Styles/App.css";
import Footer from "./Footer";
// import CardList from "./CardList";
import Header from "./Header";


class Dashboard extends Component {
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
                <input type="text" class="btn-lg" id="textInput"></input>
                <button type="button" class="btn personalButton btn-lg">Search</button>
              </div>
            </div>

          <div style={{ verticalAlign: "center" }}>
            <main>
              <div>
                <section>
                  {/* <CardList/> */}
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
