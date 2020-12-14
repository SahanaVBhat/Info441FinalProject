import React from "react";
import { withRouter } from "react-router";


class Card extends React.Component {

goToCarddetails = (cardId) => {
    localStorage.setItem("selectedCard", cardId);
    this.props.history.push(cardId);
}

render() {
    return ( 
        <div className="card"onClick = {()=>this.goToCarddetails(this.props.classInfo[0].toString())}>
        <div className="card-body">
          <h3 className="card-title">{this.props.classInfo[1]} {this.props.classInfo[2]}</h3>
          <h4>{this.props.classInfo}</h4>
        </div>
      </div>
    
    )
  }
}

export default withRouter(Card);