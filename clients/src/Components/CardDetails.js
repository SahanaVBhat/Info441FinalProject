import React from "react";


class CardDetails extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            courseID: localStorage.getItem("selectedCard"),
            evals: [],
            error: ""
        };
    }

    getEvals = async () => {
        const response = await fetch(api.base + api.handlers.courses + "/" + this.state.coourseID + "/evaluations", { method: "GET" });
        if (response.status >= 300) {
            const error = await response.text();
            this.setError(error);
            return;
        }

        const evals = await response.json();

        // check if there is more than one evaluation for course 
        if (evals.length > 0) {
            // get current evals saved in state
            let currEvals = this.state.evals;

            // for each evaluation:
            // -- create array containg all eval information to display
            // -- add eval to currEvals (for state)
            evals.forEach(function () {
                let evalInfo = [eval.id, eval.studentID, eval.courseID, eval.instructors, eval.year, eval.quarter, eval.creditType, eval.credits, eval.workload, eval.gradingTechniques, eval.description, eval.likedUsers, eval.dislikedUsers, eval.createdAt, eval.editedAt]

                currEvals.push(evalInfo);
            })

            this.setState({
                evals: currEvals
            })
            this.setError("");
        }
    }
    render() {
        this.getEvals();
        //let selectedCardId = localStorage.getItem("selectedCard");
        // you can get this cardId anywhere in the component as per your requirement 
        return (
        <div>
            {this.state.coourseID }
            {this.state.evals}
        </div>
        )
    }
}

export default CardDetails;