import React, { Component } from "react";
import "../Styles/App.css";

class Card extends Component {
  render() {
    let topic = this.props.category;
    return (
      <div className="card">
        <div className="card-body">
          <h3 className="card-title">{this.props.course} {this.props.name}</h3>
          <h4>{this.props.d}</h4>
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
