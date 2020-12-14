import React, { Component } from "react";
import "../Styles/App.css";

class Card extends Component {
  render() {
    return (
        <div class="card text-center">
          <div class="card-header">
            {this.props.course} 
          </div>
          <div class="card-body">
            <h5 class="card-title">{this.props.name}</h5>
            <p class="card-text">{this.props.d}</p>
            <a href="#" class="btn btn-primary">Add Evaluation</a>
          </div>
        </div>
    );
  }
}
class CardList extends Component {
  render() {
    return <Card course={this.props.classInfo[1]} name={this.props.classInfo[2]} d={this.props.classInfo[3]}/>;
  }
}
export default CardList;
