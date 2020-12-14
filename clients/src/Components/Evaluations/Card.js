import React from "react";
import { withRouter } from "react-router";


class Card extends React.Component {

goToCarddetails = (cardId) => {
    localStorage.setItem("selectedCard", cardId);
    this.props.history.push(cardId);
}

render() {
    return ( 
        <div className="card"onClick = {()=>this.goToCarddetails(this.props.classInfo.id.toString())}>
        <div className="card-body">
    <h3 className="card-title">Evaluation #{this.props.value+1}</h3>
          <h5>Credit Type: {this.props.classInfo.creditType}</h5>
          <h5>Instructors: {this.props.classInfo.instructors.map((instructor) => instructor.name)}</h5>
          <h5>Quarter: {this.props.classInfo.quarter}</h5>
          <h5>Credit Type: {this.props.classInfo.creditType}</h5>
          <h5>Workload: {this.props.classInfo.workload}</h5>
          <h5>Grading fairness: {this.props.classInfo.gradingTechniques}</h5>
          <h5>Comments: {this.props.classInfo.description}</h5>
        </div>
      </div>
    
    )
  }
}

export default withRouter(Card);